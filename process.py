import asyncio
import random
from loguru import logger


import src.utils
import src.model
from src.utils.check_github_version import check_version
from src.utils.logs import (
    ProgressTracker,
    create_progress_tracker,
    save_task_logs_to_excel,
)
from src.utils.config_browser import run
from src.utils.reader import read_accounts_from_excel
from src.utils.telegram_logger import send_telegram_message


async def start():
    async def launch_wrapper(index, account):
        async with semaphore:
            await account_flow(
                index,
                account,
                config,
                progress_tracker,
            )

    print("\nAvailable options:\n")
    print("[1] ⭐️ Start farming")
    print("[2] 🔧 Edit config")
    print("[3] 👋 Exit")
    print()

    try:
        choice = input("Enter option (1-3): ").strip()
    except Exception as e:
        logger.error(f"Input error: {e}")
        return

    if choice == "3" or not choice:
        return
    elif choice == "2":
        run()
        return
    elif choice == "1":
        pass
    else:
        logger.error(f"Invalid choice: {choice}")
        return

    # Get tasks from menu selection
    selected_tasks = src.utils.show_menu(src.utils.MAIN_MENU_OPTIONS)

    # Check if user selected exit
    if "Exit" in selected_tasks:
        return

    config = src.utils.get_config()
    prepare_data = await src.model.prepare_data(selected_tasks)
    config.FLOW.TASKS_DATA = prepare_data

    # Save selected tasks to config for use in Start class
    config.FLOW.TASKS = selected_tasks

    # Читаем аккаунты из Excel вместо базы данных
    accounts_to_process = read_accounts_from_excel("data/accounts.xlsx")

    if not accounts_to_process:
        logger.error("No accounts found in Excel file. Please add accounts first.")
        input("Press Enter to continue...")
        return

    # Определяем диапазон аккаунтов
    start_index = config.SETTINGS.ACCOUNTS_RANGE[0]
    end_index = config.SETTINGS.ACCOUNTS_RANGE[1]

    # Если оба 0, проверяем EXACT_ACCOUNTS_TO_USE
    if start_index == 0 and end_index == 0:
        if config.SETTINGS.EXACT_ACCOUNTS_TO_USE:
            # Валидируем, что указанные номера не превышают количество аккаунтов
            valid_indices = [
                i
                for i in config.SETTINGS.EXACT_ACCOUNTS_TO_USE
                if i <= len(accounts_to_process)
            ]
            # Преобразуем номера аккаунтов в индексы (номер - 1)
            selected_indices = [i - 1 for i in valid_indices]
            final_accounts = [
                accounts_to_process[i]
                for i in selected_indices
                if i < len(accounts_to_process)
            ]
            logger.info(f"Using specific accounts: {valid_indices}")

            # Для совместимости с остальным кодом
            start_index = min(valid_indices) if valid_indices else 1
            end_index = (
                max(valid_indices) if valid_indices else len(accounts_to_process)
            )
        else:
            # Если список пустой, берем все аккаунты
            final_accounts = accounts_to_process
            start_index = 1
            end_index = len(accounts_to_process)
    else:
        # Проверяем, что индексы в пределах списка
        if start_index < 1:
            start_index = 1
        if end_index > len(accounts_to_process):
            end_index = len(accounts_to_process)

        # Python slice не включает последний элемент, поэтому +1
        final_accounts = accounts_to_process[start_index - 1 : end_index]

    threads = config.SETTINGS.THREADS

    # Создаем список индексов
    indices = list(range(len(final_accounts)))

    # Перемешиваем индексы только если включен SHUFFLE_ACCOUNTS
    if config.SETTINGS.SHUFFLE_ACCOUNTS:
        random.shuffle(indices)
        shuffle_status = "random"
    else:
        shuffle_status = "sequential"

    # Создаем строку с порядком аккаунтов
    if config.SETTINGS.EXACT_ACCOUNTS_TO_USE:
        # Создаем список номеров аккаунтов в нужном порядке
        valid_indices = [
            i
            for i in config.SETTINGS.EXACT_ACCOUNTS_TO_USE
            if i <= len(accounts_to_process)
        ]
        ordered_accounts = [valid_indices[i] for i in indices if i < len(valid_indices)]
        account_order = " ".join(map(str, ordered_accounts))
        logger.info(f"Starting with specific accounts in {shuffle_status} order...")
    else:
        account_order = " ".join(str(start_index + idx) for idx in indices)
        logger.info(
            f"Starting with accounts {start_index} to {end_index} in {shuffle_status} order..."
        )
    logger.info(f"Accounts order: {account_order}")

    semaphore = asyncio.Semaphore(value=threads)
    tasks_list = []

    # Add before creating tasks
    progress_tracker = await create_progress_tracker(
        total=len(final_accounts), description="Accounts completed"
    )

    # Используем индексы для создания задач
    for idx in indices:
        if idx >= len(final_accounts):
            continue

        actual_index = (
            config.SETTINGS.EXACT_ACCOUNTS_TO_USE[idx]
            if config.SETTINGS.EXACT_ACCOUNTS_TO_USE
            and idx < len(config.SETTINGS.EXACT_ACCOUNTS_TO_USE)
            else start_index + idx
        )
        tasks_list.append(
            asyncio.create_task(
                launch_wrapper(
                    actual_index,
                    final_accounts[idx],
                )
            )
        )

    await asyncio.gather(*tasks_list)

    # After all accounts are processed, save task logs to Excel
    logger.info("Saving task logs to Excel files...")
    summary = await save_task_logs_to_excel()

    # Send summary report to Telegram
    if config.SETTINGS.SEND_TELEGRAM_LOGS and summary:
        try:
            # Create a detailed message about the execution results
            message = (
                f"🌟 StarLabs Twitter Bot Summary Report 🌟\n\n"
                f"📊 Statistics:\n"
                f"Total Accounts: {summary['total_accounts']}\n"
                f"Total Task Executions: {summary['total_tasks']}\n"
                f"Successful Executions: {summary['success_count']}\n"
                f"Failed Executions: {summary['total_tasks'] - summary['success_count']}\n"
                f"Success Rate: {(summary['success_count'] / summary['total_tasks'] * 100) if summary['total_tasks'] > 0 else 0:.1f}%\n\n"
            )

            # Add task-specific statistics
            message += "✅ Task Success Breakdown:\n"
            for i, (task_name, success_count) in enumerate(
                summary["task_success_counts"].items(), 1
            ):
                total_for_task = summary["task_counts"].get(task_name, 0)
                success_rate = (
                    (success_count / total_for_task * 100) if total_for_task > 0 else 0
                )
                message += f"{i}. {task_name.capitalize()}: {success_count}/{total_for_task} ({success_rate:.1f}%)\n"

            message += f"\n⚙️ Settings:\n"
            message += f"Skip Failed: {'Yes' if config.FLOW.SKIP_FAILED_TASKS else 'No'}\n"
            message += f"Tasks: {', '.join(config.FLOW.TASKS)}\n"

            # Send the message
            await send_telegram_message(config, message)
        except Exception as e:
            logger.error(f"Error sending summary Telegram message: {e}")

    logger.success("Task completed for all accounts.")
    input("Press Enter to continue...")


async def account_flow(
    account_index: int,
    account,
    config: src.utils.config.Config,
    progress_tracker: ProgressTracker,
):
    try:
        pause = random.randint(
            config.SETTINGS.RANDOM_INITIALIZATION_PAUSE[0],
            config.SETTINGS.RANDOM_INITIALIZATION_PAUSE[1],
        )
        logger.info(f"[{account_index}] Sleeping for {pause} seconds before start...")
        await asyncio.sleep(pause)

        # Create instance with tasks from config
        instance = src.model.Start(
            account_index, account.proxy, account.auth_token, config
        )

        result = await wrapper(instance.initialize, config)
        if not result:
            # If initialization failed, update the progress tracker and return
            logger.error(f"[{account_index}] Account initialization failed")
            await progress_tracker.increment(1)
            return

        # Call the flow method with tasks from config
        await instance.flow()

        pause = random.randint(
            config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[0],
            config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[1],
        )
        logger.info(f"Sleeping for {pause} seconds before next account...")
        await asyncio.sleep(pause)

        # Add progress update
        await progress_tracker.increment(1)

    except Exception as err:
        logger.error(f"{account_index} | Account flow failed: {err}")
        # Update progress even if there's an error
        await progress_tracker.increment(1)


async def wrapper(function, config: src.utils.config.Config, *args, **kwargs):
    attempts = config.SETTINGS.ATTEMPTS
    for attempt in range(attempts):
        result = await function(*args, **kwargs)
        if isinstance(result, tuple) and result and isinstance(result[0], bool):
            if result[0]:
                return result
        elif isinstance(result, bool):
            if result:
                return True

        if attempt < attempts - 1:  # Don't sleep after the last attempt
            pause = random.randint(
                config.SETTINGS.PAUSE_BETWEEN_ATTEMPTS[0],
                config.SETTINGS.PAUSE_BETWEEN_ATTEMPTS[1],
            )
            logger.info(
                f"Sleeping for {pause} seconds before next attempt {attempt+1}/{config.SETTINGS.ATTEMPTS}..."
            )
            await asyncio.sleep(pause)

    return result


def task_exists_in_config(task_name: str, tasks_list: list) -> bool:
    """Рекурсивно проверяет наличие задачи в списке задач, включая вложенные списки"""
    for task in tasks_list:
        if isinstance(task, list):
            if task_exists_in_config(task_name, task):
                return True
        elif task == task_name:
            return True
    return False
