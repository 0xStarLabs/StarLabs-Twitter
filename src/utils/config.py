from dataclasses import dataclass, field
from typing import List, Tuple
import yaml
import asyncio

from src.utils.constants import DataForTasks


@dataclass
class SettingsConfig:
    THREADS: int
    ATTEMPTS: int
    ACCOUNTS_RANGE: Tuple[int, int]
    EXACT_ACCOUNTS_TO_USE: List[int]
    PAUSE_BETWEEN_ATTEMPTS: Tuple[int, int]
    RANDOM_PAUSE_BETWEEN_ACCOUNTS: Tuple[int, int]
    RANDOM_PAUSE_BETWEEN_ACTIONS: Tuple[int, int]
    RANDOM_INITIALIZATION_PAUSE: Tuple[int, int]
    TELEGRAM_USERS_IDS: List[int]
    TELEGRAM_BOT_TOKEN: str
    SEND_TELEGRAM_LOGS: bool
    SEND_ONLY_SUMMARY: bool
    SHUFFLE_ACCOUNTS: bool


@dataclass
class FlowConfig:
    TASKS_DATA: DataForTasks | None
    SKIP_FAILED_TASKS: bool
    TASKS: List[str] = field(default_factory=list)


@dataclass
class TweetsConfig:
    RANDOM_TEXT_FOR_TWEETS: bool
    RANDOM_PICTURE_FOR_TWEETS: bool


@dataclass
class CommentsConfig:
    RANDOM_TEXT_FOR_COMMENTS: bool
    RANDOM_PICTURE_FOR_COMMENTS: bool


@dataclass
class OthersConfig:
    SSL_VERIFICATION: bool


@dataclass
class Config:
    SETTINGS: SettingsConfig
    FLOW: FlowConfig
    TWEETS: TweetsConfig
    COMMENTS: CommentsConfig
    OTHERS: OthersConfig
    lock: asyncio.Lock = field(default_factory=asyncio.Lock)

    @classmethod
    def load(cls, path: str = "config.yaml") -> "Config":
        """Load configuration from yaml file"""
        with open(path, "r", encoding="utf-8") as file:
            data = yaml.safe_load(file)

        return cls(
            SETTINGS=SettingsConfig(
                THREADS=data["SETTINGS"]["THREADS"],
                ATTEMPTS=data["SETTINGS"]["ATTEMPTS"],
                ACCOUNTS_RANGE=tuple(data["SETTINGS"]["ACCOUNTS_RANGE"]),
                EXACT_ACCOUNTS_TO_USE=data["SETTINGS"]["EXACT_ACCOUNTS_TO_USE"],
                PAUSE_BETWEEN_ATTEMPTS=tuple(
                    data["SETTINGS"]["PAUSE_BETWEEN_ATTEMPTS"]
                ),
                RANDOM_PAUSE_BETWEEN_ACCOUNTS=tuple(
                    data["SETTINGS"]["RANDOM_PAUSE_BETWEEN_ACCOUNTS"]
                ),
                RANDOM_PAUSE_BETWEEN_ACTIONS=tuple(
                    data["SETTINGS"]["RANDOM_PAUSE_BETWEEN_ACTIONS"]
                ),
                RANDOM_INITIALIZATION_PAUSE=tuple(
                    data["SETTINGS"]["RANDOM_INITIALIZATION_PAUSE"]
                ),
                TELEGRAM_USERS_IDS=data["SETTINGS"]["TELEGRAM_USERS_IDS"],
                TELEGRAM_BOT_TOKEN=data["SETTINGS"]["TELEGRAM_BOT_TOKEN"],
                SEND_TELEGRAM_LOGS=data["SETTINGS"]["SEND_TELEGRAM_LOGS"],
                SEND_ONLY_SUMMARY=data["SETTINGS"]["SEND_ONLY_SUMMARY"],
                SHUFFLE_ACCOUNTS=data["SETTINGS"].get("SHUFFLE_ACCOUNTS", True),
            ),
            FLOW=FlowConfig(
                TASKS_DATA=None,
                SKIP_FAILED_TASKS=data["FLOW"]["SKIP_FAILED_TASKS"],
                TASKS=[],
            ),
            TWEETS=TweetsConfig(
                RANDOM_TEXT_FOR_TWEETS=data["TWEETS"]["RANDOM_TEXT_FOR_TWEETS"],
                RANDOM_PICTURE_FOR_TWEETS=data["TWEETS"]["RANDOM_PICTURE_FOR_TWEETS"],
            ),
            COMMENTS=CommentsConfig(
                RANDOM_TEXT_FOR_COMMENTS=data["COMMENTS"]["RANDOM_TEXT_FOR_COMMENTS"],
                RANDOM_PICTURE_FOR_COMMENTS=data["COMMENTS"][
                    "RANDOM_PICTURE_FOR_COMMENTS"
                ],
            ),
            OTHERS=OthersConfig(
                SSL_VERIFICATION=data["OTHERS"]["SSL_VERIFICATION"],
            ),
        )


# Singleton pattern
def get_config() -> Config:
    """Get configuration singleton"""
    if not hasattr(get_config, "_config"):
        get_config._config = Config.load()
    return get_config._config
