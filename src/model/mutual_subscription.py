import asyncio
import random
from loguru import logger
from typing import List, Dict, Any
import time

from src.model.instance import Instance
from src.utils.config import Config
from src.utils.constants import Account, DataForTasks
from src.utils.logs import create_progress_tracker


async def create_mutual_subscription_plan(
    accounts: List[Instance], followers_per_account: int
) -> Dict[str, List[str]]:
    """
    Create a subscription plan for mutual follows.

    Args:
        accounts: List of account instances
        followers_per_account: Number of followers each account should have

    Returns:
        A dictionary mapping username to a list of usernames that should follow it
    """
    subscription_plan = {}

    # Get all valid usernames
    valid_accounts = [acc for acc in accounts if acc.username is not None]

    if not valid_accounts:
        logger.error("No valid accounts found for mutual subscription")
        return {}

    for account in valid_accounts:
        # For each account, determine how many followers it should get
        follower_count = min(followers_per_account, len(valid_accounts) - 1)

        # Available followers (excluding self)
        available_followers = [
            acc for acc in valid_accounts if acc.username != account.username
        ]

        # Shuffle available followers
        random.shuffle(available_followers)

        # Select required number of followers
        selected_followers = available_followers[:follower_count]

        # Add to subscription plan - who should follow this account
        for follower in selected_followers:
            if follower.username not in subscription_plan:
                subscription_plan[follower.username] = []

            subscription_plan[follower.username].append(account.username)

    return subscription_plan


async def execute_mutual_subscription(
    accounts: List[Instance],
    config: Config,
    followers_per_account: int,
    progress_tracker=None,
) -> Dict[str, Any]:
    """
    Execute the mutual subscription plan.

    Args:
        accounts: List of account instances
        config: Application configuration
        followers_per_account: Number of followers each account should have
        progress_tracker: Optional progress tracker to update

    Returns:
        Statistics of the operation
    """
    start_time = time.time()
    stats = {
        "total": len(accounts),
        "processed": 0,
        "success": 0,
        "failed": 0,
    }

    logger.info(f"Starting mutual subscription with {len(accounts)} accounts")

    # Filter accounts that were successfully initialized - parallelize this step
    async def initialize_account(account):
        try:
            result = await account.initialize()
            if result:
                # Update progress tracker if provided
                if progress_tracker:
                    await progress_tracker.increment(1)
                return account
            if progress_tracker:
                await progress_tracker.increment(1)
            return None
        except Exception as e:
            # Handle authentication errors, blocks, etc. gracefully
            logger.error(f"Failed to initialize account {account.account_index}: {e}")
            # Update progress tracker if provided
            if progress_tracker:
                await progress_tracker.increment(1)
            # Return None for failed accounts
            return None

    # Use a semaphore to limit concurrency
    thread_count = config.SETTINGS.THREADS
    semaphore = asyncio.Semaphore(thread_count)

    async def bounded_initialize(account):
        async with semaphore:
            return await initialize_account(account)

    # Initialize accounts in parallel
    logger.info(f"Initializing accounts using {thread_count} threads")
    initialization_tasks = [bounded_initialize(account) for account in accounts]

    # Use gather with return_exceptions=True to prevent exceptions from propagating
    initialized_results = await asyncio.gather(
        *initialization_tasks, return_exceptions=True
    )

    # Filter out None results and exceptions
    initialized_accounts = []
    for result in initialized_results:
        if isinstance(result, Exception):
            logger.error(f"Account initialization task failed with error: {result}")
            continue
        if result is not None:
            initialized_accounts.append(result)

    logger.info(
        f"Successfully initialized {len(initialized_accounts)}/{len(accounts)} accounts"
    )

    if not initialized_accounts:
        logger.error("No accounts could be initialized. Aborting.")
        stats["processed"] = len(accounts)
        stats["failed"] = len(accounts)
        return stats

    # Create a new progress tracker for the subscription process if the previous one exists
    if progress_tracker:
        progress_tracker = await create_progress_tracker(
            total=len(initialized_accounts),
            description="Mutual subscriptions completed",
        )

    # Create subscription plan
    subscription_plan = await create_mutual_subscription_plan(
        initialized_accounts, followers_per_account
    )

    if not subscription_plan:
        logger.error("Failed to create subscription plan. Aborting.")
        stats["processed"] = len(accounts)
        stats["failed"] = len(accounts)
        return stats

    # Execute the subscription plan - also parallelize this step using the semaphore
    async def bounded_process_subscription(account, usernames_to_follow):
        try:
            async with semaphore:
                await process_account_subscription(
                    account, usernames_to_follow, config, stats, progress_tracker
                )
        except Exception as e:
            logger.error(
                f"Error during subscription process for account {account.username}: {e}"
            )
            stats["processed"] += 1
            stats["failed"] += 1
            if progress_tracker:
                await progress_tracker.increment(1)

    tasks = []
    for account in initialized_accounts:
        if account.username in subscription_plan:
            usernames_to_follow = subscription_plan[account.username]
            logger.info(
                f"Account {account.username} will follow {len(usernames_to_follow)} accounts"
            )

            # Execute mutual subscription for this account with bounded concurrency
            task = asyncio.create_task(
                bounded_process_subscription(account, usernames_to_follow)
            )
            tasks.append(task)
        else:
            logger.info(f"Account {account.username} has no accounts to follow")
            stats["processed"] += 1
            stats["success"] += 1
            if progress_tracker:
                await progress_tracker.increment(1)

    # Wait for all tasks to complete, handling exceptions
    if tasks:
        await asyncio.gather(*tasks, return_exceptions=True)

    # Calculate elapsed time
    elapsed_time = time.time() - start_time
    logger.info(f"Mutual subscription completed in {elapsed_time:.2f} seconds")
    logger.info(
        f"Processed: {stats['processed']}/{stats['total']}, Success: {stats['success']}, Failed: {stats['failed']}"
    )

    return stats


async def process_account_subscription(
    account: Instance,
    usernames_to_follow: List[str],
    config: Config,
    stats: Dict[str, Any],
    progress_tracker=None,
) -> None:
    """
    Process mutual subscription for a single account.

    Args:
        account: The account instance
        usernames_to_follow: List of usernames this account should follow
        config: Application configuration
        stats: Statistics object to update
        progress_tracker: Optional progress tracker to update
    """
    try:
        logger.info(
            f"Account {account.username} following {len(usernames_to_follow)} accounts"
        )

        # Instead of calling mutual_subscription that follows all accounts at once,
        # we'll follow them one by one with random delays in between
        all_success = True
        for username in usernames_to_follow:
            try:
                # Follow a single user
                success = await account.twitter.follow(username)

                if success:
                    logger.success(
                        f"Account {account.username} successfully followed {username}"
                    )
                else:
                    logger.error(
                        f"Account {account.username} failed to follow {username}"
                    )
                    all_success = False

                # Check account status after follow - if account is locked or has other issues, break the loop
                if account.twitter.account_status != "ok":
                    logger.warning(
                        f"Account {account.username} status changed to {account.twitter.account_status}, stopping follow operations"
                    )
                    all_success = False
                    break

            except Exception as e:
                logger.error(
                    f"Error while account {account.username} was following {username}: {e}"
                )
                all_success = False
                # Continue with next username rather than breaking completely

            # Add a random delay between follows within the same account
            # Use a smaller range than between accounts
            random_pause = random.uniform(
                config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACTIONS[0],
                config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACTIONS[1],
            )
            logger.info(
                f"Account {account.username} waiting {random_pause:.2f} seconds before next follow"
            )
            await asyncio.sleep(random_pause)

        # Update statistics
        stats["processed"] += 1
        if all_success:
            stats["success"] += 1
            logger.success(
                f"Account {account.username} successfully completed mutual subscription"
            )
        else:
            stats["failed"] += 1
            logger.error(
                f"Account {account.username} had some failures during mutual subscription"
            )

        # Random pause between accounts
        random_pause = random.uniform(
            config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[0],
            config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[1],
        )
        logger.info(f"Waiting {random_pause:.2f} seconds before next account")
        await asyncio.sleep(random_pause)

        # Update progress tracker if provided
        if progress_tracker:
            await progress_tracker.increment(1)

    except Exception as e:
        logger.error(f"Error processing account {account.username}: {e}")
        stats["processed"] += 1
        stats["failed"] += 1
        # Also update progress tracker on error
        if progress_tracker:
            await progress_tracker.increment(1)
