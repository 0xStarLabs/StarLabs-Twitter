from dataclasses import dataclass


class Account:
    def __init__(self, auth_token: str, proxy: str = ""):
        self.auth_token = auth_token
        self.proxy = proxy

    def __repr__(self):
        return f"Account(auth_token={self.auth_token}, proxy={self.proxy})"


@dataclass
class DataForTasks:
    """
    Class for storing data for tasks
    """

    USERNAMES_TO_FOLLOW: list[str]
    USERNAMES_TO_UNFOLLOW: list[str]
    LINKS_TO_LIKE: list[str]
    LINKS_TO_RETWEET: list[str]
    TEXT_FOR_TWEETS: list[str]
    IMAGES: list[str]
    LINKS_FOR_COMMENTS: list[str]
    LINKS_FOR_QUOTES: list[str]
    COMMENTS: list[str]


MAIN_MENU_OPTIONS = [
    "Follow",
    "Like",
    "Retweet",
    "Comment",
    "Comment with image",
    "Tweet",
    "Tweet with image",
    "Quote",
    "Quote with image",
    "Unfollow",
    "Check Valid",
    "Exit",
]
