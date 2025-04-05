from loguru import logger
import random
import asyncio

from src.model.instance import Instance
from src.model.twitter.client import Twitter
from src.utils.config import Config
from src.utils.telegram_logger import send_telegram_message
from src.utils.constants import Account
from src.utils.logs import update_account_in_excel, add_task_log_entry


class Start:
    def __init__(
        self,
        account_index: int,
        proxy: str,
        auth_token: str,
        config: Config,
        task_names: list = None,
    ):
        self.account_index = account_index
        self.proxy = proxy
        self.auth_token = auth_token
        self.config = config
        # Use tasks from config if not explicitly passed
        self.task_names = task_names or self.config.FLOW.TASKS
        self.completed_tasks = []

        self.twitter_instance: None | Instance = None

    async def initialize(self):
        try:
            self.twitter_instance = Instance(
                Account(auth_token=self.auth_token, proxy=self.proxy),
                self.config,
                self.account_index,
                self.config.FLOW.TASKS_DATA,
            )
            ok = await self.twitter_instance.initialize()
            if self.twitter_instance.username:
                pass

            await update_account_in_excel(
                "data/accounts.xlsx",
                self.auth_token,
                self.twitter_instance.username,
                self.twitter_instance.account_status,
            )

            if not ok:
                logger.error(f"{self.account_index} | Failed to initialize")
                return False

            return True
        except Exception as e:
            logger.error(f"{self.account_index} | Error: {e}")
            return False

    async def flow(self):
        try:
            if not self.task_names:
                logger.warning(
                    f"{self.account_index} | No tasks specified for this account. Exiting..."
                )
                return True

            task_plan_msg = [f"{i+1}. {task}" for i, task in enumerate(self.task_names)]
            logger.info(
                f"{self.account_index} | Task execution plan: {' | '.join(task_plan_msg)}"
            )

            self.completed_tasks = []
            failed_tasks = []

            # Выполняем задачи
            for task_name in self.task_names:
                if task_name == "skip" or task_name == "exit":
                    logger.info(f"{self.account_index} | Skipping task: {task_name}")
                    continue

                logger.info(f"{self.account_index} | Executing task: {task_name}")

                success = await self.execute_task(task_name)

                # Log the task execution result
                log_status = "success" if success else "failed"
                await add_task_log_entry(
                    task_name.lower(), self.account_index, self.auth_token, log_status
                )

                if success:
                    self.completed_tasks.append(task_name)
                    logger.info(f"{self.account_index} | Task '{task_name}' completed")
                    await self.sleep(task_name)
                else:
                    failed_tasks.append(task_name)
                    if not self.config.FLOW.SKIP_FAILED_TASKS:
                        logger.error(
                            f"{self.account_index} | Failed to complete task {task_name}. Stopping account execution."
                        )
                        break
                    else:
                        logger.warning(
                            f"{self.account_index} | Failed to complete task {task_name}. Skipping to next task."
                        )
                        await self.sleep(task_name)

            # Отправляем сообщение в Telegram только в конце всей работы
            if self.config.SETTINGS.SEND_TELEGRAM_LOGS and not self.config.SETTINGS.SEND_ONLY_SUMMARY:
                message = (
                    f"😎 Twitter StarLabs Bot Report\n\n"
                    f"💳 Account: {self.account_index} | <code>{self.auth_token[:6]}...{self.auth_token[-4:]}</code>\n\n"
                )

                if self.completed_tasks:
                    message += f"✅ Completed Tasks:\n"
                    for i, task in enumerate(self.completed_tasks, 1):
                        message += f"{i}. {task}\n"
                    message += "\n"

                if failed_tasks:
                    message += f"❌ Failed Tasks:\n"
                    for i, task in enumerate(failed_tasks, 1):
                        message += f"{i}. {task}\n"
                    message += "\n"

                total_tasks = len(self.task_names)
                completed_count = len(self.completed_tasks)
                message += (
                    f"📊 Statistics:\n"
                    f"Total Tasks: {total_tasks}\n"
                    f"Completed: {completed_count}\n"
                    f"Failed: {len(failed_tasks)}\n"
                    f"Success Rate: {(completed_count/total_tasks)*100:.1f}%\n\n"
                    f"⚙️ Settings:\n"
                    f"Skip Failed: {'Yes' if self.config.FLOW.SKIP_FAILED_TASKS else 'No'}\n"
                )

                await send_telegram_message(self.config, message)

            return len(failed_tasks) == 0

        except Exception as e:
            logger.error(f"{self.account_index} | Error: {e}")

            if self.config.SETTINGS.SEND_TELEGRAM_LOGS and not self.config.SETTINGS.SEND_ONLY_SUMMARY:
                error_message = (
                    f"⚠️ Error Report\n\n"
                    f"Account #{self.account_index}\n"
                    f"Token: <code>{self.auth_token[:6]}...{self.auth_token[-4:]}</code>\n"
                    f"Error: {str(e)}"
                )
                await send_telegram_message(self.config, error_message)

            return False

    async def execute_task(self, task):
        """Execute a single task"""
        task = task.lower()

        if task == "like":
            return await self.twitter_instance.like()

        if task == "retweet":
            return await self.twitter_instance.retweet()

        if "comment" in task:
            return await self.twitter_instance.comment(task)

        if "tweet" in task or "quote" in task:
            return await self.twitter_instance.tweet(task)

        if task == "follow":
            return await self.twitter_instance.follow()

        if task == "unfollow":
            return await self.twitter_instance.unfollow()

        if task == "check valid":
            # For check valid, we just return True as the account has already been initialized
            logger.info(f"{self.account_index} | Account is valid")
            return True

        logger.error(f"{self.account_index} | Task {task} not found")
        return False

    async def sleep(self, task_name: str):
        """Делает рандомную паузу между действиями"""
        pause = random.randint(
            self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACTIONS[0],
            self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACTIONS[1],
        )
        logger.info(
            f"{self.account_index} | Sleeping {pause} seconds after {task_name}"
        )
        await asyncio.sleep(pause)
