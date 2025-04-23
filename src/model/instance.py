import asyncio
import random
from loguru import logger
from src.model.twitter.client import Twitter
from src.utils.config import Config
from src.utils.constants import Account, DataForTasks


class Instance:
    def __init__(
        self,
        account: Account,
        config: Config,
        account_index: int,
        prepare_data: DataForTasks,
    ):
        self.account = account
        self.config = config
        self.account_index = account_index
        self.prepare_data = prepare_data

        self.account_status: str = "ok"
        self.username: str | None = None

    async def initialize(self):
        self.twitter = Twitter(
            self.account_index, self.account.auth_token, self.account.proxy, self.config
        )
        ok = await self.twitter.initialize()
        if self.twitter.username:
            self.username = self.twitter.username

        if not ok:
            return False

        return True

    async def like(self):
        try:
            if not self.prepare_data.LINKS_TO_LIKE:
                logger.error(f"{self.account_index} | No links to like found")
                return False

            success = True
            for tweet_id in self.prepare_data.LINKS_TO_LIKE:
                success = await self.twitter.like(tweet_id)
                if not success:
                    success = False
                random_pause = random.uniform(
                    self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACTIONS[0],
                    self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACTIONS[1],
                )
                logger.info(
                    f"{self.account_index} | Waiting {random_pause} seconds before next like"
                )
                await asyncio.sleep(random_pause)
            return success

        except Exception as e:
            logger.error(f"{self.account_index} | Error liking tweet: {e}")
            return False

    async def retweet(self):
        try:
            if not self.prepare_data.LINKS_TO_RETWEET:
                logger.error(f"{self.account_index} | No links to retweet found")
                return False

            success = True
            for tweet_id in self.prepare_data.LINKS_TO_RETWEET:
                success = await self.twitter.retweet(tweet_id)
                if not success:
                    success = False
                random_pause = random.uniform(
                    self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[0],
                    self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[1],
                )
                logger.info(
                    f"{self.account_index} | Waiting {random_pause} seconds before next retweet"
                )
                await asyncio.sleep(random_pause)
            return success

        except Exception as e:
            logger.error(f"{self.account_index} | Error retweeting tweet: {e}")
            return False

    async def unfollow(self):
        try:
            if not self.prepare_data.USERNAMES_TO_UNFOLLOW:
                logger.error(f"{self.account_index} | No usernames to unfollow found")
                return False

            success = True
            for username in self.prepare_data.USERNAMES_TO_UNFOLLOW:
                success = await self.twitter.unfollow(username)
                if not success:
                    success = False
                random_pause = random.uniform(
                    self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[0],
                    self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[1],
                )
                logger.info(
                    f"{self.account_index} | Waiting {random_pause} seconds before next unfollow"
                )
                await asyncio.sleep(random_pause)
            return success

        except Exception as e:
            logger.error(f"{self.account_index} | Error unfollowing user: {e}")
            return False

    async def follow(self):
        try:
            if not self.prepare_data.USERNAMES_TO_FOLLOW:
                logger.error(f"{self.account_index} | No usernames to follow found")
                return False

            success = True
            for username in self.prepare_data.USERNAMES_TO_FOLLOW:
                success = await self.twitter.follow(username)
                if not success:
                    success = False
                random_pause = random.uniform(
                    self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[0],
                    self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[1],
                )
                logger.info(
                    f"{self.account_index} | Waiting {random_pause} seconds before next follow"
                )
                await asyncio.sleep(random_pause)
            return success

        except Exception as e:
            logger.error(f"{self.account_index} | Error following user: {e}")
            return False

    async def tweet(self, task: str):
        """
        Types: default, quote, image,
        """
        try:
            image = None
            quote_link = None
            success = True

            if not self.prepare_data.TEXT_FOR_TWEETS:
                logger.error(f"{self.account_index} | No text for tweets found")
                return False

            if self.config.TWEETS.RANDOM_TEXT_FOR_TWEETS:
                tweet_text = random.choice(self.prepare_data.TEXT_FOR_TWEETS)
            else:
                if len(self.prepare_data.TEXT_FOR_TWEETS) >= self.account_index - 1:
                    tweet_text = self.prepare_data.TEXT_FOR_TWEETS[
                        self.account_index - 1
                    ]
                else:
                    logger.error(
                        f"{self.account_index} | In your tweets.txt file, you have less tweets than accounts."
                    )
                    return False

            if "image" in task:
                if not self.prepare_data.IMAGES:
                    logger.error(f"{self.account_index} | No images found")
                    return False

                if self.config.TWEETS.RANDOM_PICTURE_FOR_TWEETS:
                    image = random.choice(self.prepare_data.IMAGES)
                else:
                    if len(self.prepare_data.IMAGES) >= self.account_index - 1:
                        image = self.prepare_data.IMAGES[self.account_index - 1]
                    else:
                        logger.error(
                            f"{self.account_index} | In your images.txt file, you have less images than accounts."
                        )
                        return False

            if "quote" in task:
                if not self.prepare_data.LINKS_FOR_QUOTES:
                    logger.error(f"{self.account_index} | No links for quotes found")
                    return False

                for link in self.prepare_data.LINKS_FOR_QUOTES:
                    success = await self.twitter.tweet(tweet_text, link, image)
                    if not success:
                        success = False
                    random_pause = random.uniform(
                        self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[0],
                        self.config.SETTINGS.RANDOM_PAUSE_BETWEEN_ACCOUNTS[1],
                    )
                    logger.info(
                        f"{self.account_index} | Waiting {random_pause} seconds before next tweet"
                    )
                    await asyncio.sleep(random_pause)
            else:
                return await self.twitter.tweet(tweet_text, quote_link, image)

            return success

        except Exception as e:
            logger.error(f"{self.account_index} | Error following user: {e}")
            return False

    async def comment(self, task: str):
        """
        Types: default, picture
        """
        try:
            image = None
            success = True

            if not self.prepare_data.COMMENTS:
                logger.error(f"{self.account_index} | No text for comments found")
                return False

            if not self.prepare_data.LINKS_FOR_COMMENTS:
                logger.error(f"{self.account_index} | No links for comments found")
                return False

            if self.config.COMMENTS.RANDOM_TEXT_FOR_COMMENTS:
                comment_text = random.choice(self.prepare_data.COMMENTS)
            else:

                if len(self.prepare_data.COMMENTS) >= self.account_index - 1:
                    comment_text = self.prepare_data.COMMENTS[self.account_index - 1]
                else:
                    logger.error(
                        f"{self.account_index} | In your comments.txt file, you have less comments than accounts."
                    )
                    return False

            if "image" in task:
                if not self.prepare_data.IMAGES:
                    logger.error(f"{self.account_index} | No images found")
                    return False

                if self.config.COMMENTS.RANDOM_PICTURE_FOR_COMMENTS:
                    image = random.choice(self.prepare_data.IMAGES)
                else:
                    if len(self.prepare_data.IMAGES) >= self.account_index - 1:
                        image = self.prepare_data.IMAGES[self.account_index - 1]
                    else:
                        logger.error(
                            f"{self.account_index} | In your images.txt file, you have less images than accounts."
                        )
                        return False

            link_for_comment = self.prepare_data.LINKS_FOR_COMMENTS

            return await self.twitter.comment(comment_text, link_for_comment, image)

        except Exception as e:
            logger.error(f"{self.account_index} | Error commenting: {e}")
            return False

    async def mutual_subscription(self, usernames_to_follow: list[str]):
        """
        Follow a list of usernames for mutual subscription functionality.

        Args:
            usernames_to_follow (list[str]): List of usernames to follow

        Returns:
            bool: True if successful, False otherwise
        """
        try:
            if not usernames_to_follow:
                logger.error(f"{self.account_index} | No usernames to follow found")
                return False

            success = await self.twitter.mutual_subscription(usernames_to_follow)
            return success

        except Exception as e:
            logger.error(f"{self.account_index} | Error in mutual subscription: {e}")
            return False
