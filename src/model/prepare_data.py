from src.utils.constants import DataForTasks
from src.utils.config import Config
from src.utils.reader import read_pictures, read_txt_file

def extract_tweet_id(tweet_link: str) -> str:
    tweet_link = tweet_link.strip()

    if "tweet_id=" in tweet_link:
        parts = tweet_link.split("tweet_id=")
        tweet_id = parts[1].split("&")[0]
    elif "?" in tweet_link and "status/" in tweet_link:
        parts = tweet_link.split("status/")
        tweet_id = parts[1].split("?")[0]
    elif "status/" in tweet_link:
        parts = tweet_link.split("status/")
        tweet_id = parts[1]
    else:
        raise ValueError(f"Failed to get tweet ID from your link: {tweet_link}")

    return tweet_id


async def prepare_data(tasks: list[str]) -> DataForTasks:
    data = DataForTasks(
        USERNAMES_TO_FOLLOW=[],
        USERNAMES_TO_UNFOLLOW=[],
        LINKS_TO_LIKE=[],
        LINKS_TO_RETWEET=[],
        TEXT_FOR_TWEETS=[],
        IMAGES=[],
        LINKS_FOR_COMMENTS=[],
        LINKS_FOR_QUOTES=[],
        COMMENTS=[],
    )

    for task in tasks:
        task = task.lower()
        if task == "like":
            user_input = input("Enter links to tweets for likes: ")
            tweet_ids = [extract_tweet_id(link) for link in user_input.replace(",", " ").split()]
            data.LINKS_TO_LIKE = tweet_ids

        if task == "retweet":
            user_input = input("Enter links to tweets for retweets: ")
            tweet_ids = [extract_tweet_id(link) for link in user_input.replace(",", " ").split()]
            data.LINKS_TO_RETWEET = tweet_ids

        if "comment" in task:
            user_input = input("Enter the link to the tweet for comment: ")
            tweet_id = extract_tweet_id(user_input)
            data.LINKS_FOR_COMMENTS = tweet_id
            data.COMMENTS = read_txt_file("comments", "data/comment_text.txt")
        
        if "quote" in task:
            user_input = input("Enter links to tweets for quotes: ")
            tweet_ids = [link.strip() for link in user_input.replace(",", " ").split()]
            data.LINKS_FOR_QUOTES = tweet_ids

        if task == "follow":
            user_input = input("Enter usernames to follow: ")
            data.USERNAMES_TO_FOLLOW = user_input.replace(",", " ").replace("@", "").split()

        if task == "unfollow":
            user_input = input("Enter usernames to unfollow: ")
            data.USERNAMES_TO_UNFOLLOW = user_input.replace(",", " ").replace("@", "").split()

        if "tweet" in task or "quote" in task:
            data.TEXT_FOR_TWEETS = read_txt_file("tweets", "data/tweet_text.txt")
        
        if "image" in task:
            data.IMAGES = await read_pictures("data/images")

    return data
