import requests as default_requests
from random import randint
from loguru import logger
from time import sleep


def get_user_choice(tasks_list: list, text_to_show: str) -> list:
    user_choice = input(text_to_show + ": ").split()
    return [tasks_list[int(task.strip())-1] for task in user_choice]


def ask_for_task_data(all_tasks: list) -> dict | None:
    tasks_user_data = {
        "follow": {},
        "retweet": {},
        "like": {},
        "tweet": False,
        "tweet with picture": False,
        "quote tweet": {},
        "comment": {},
        "comment with picture": {},
        "change description": False,
        "change username": False,
        "change name": False,
        "change background": False,
        "change password": False,
        "change birthdate": False,
        "change location": False,
        "change profile picture": False,
        "send direct message": {},
        # "mass dm": False,
        "check dm": {},
        "check valid": False,
        "unfreeze": False,
        "mutual subscription": {
            "collect usernames": False,
            "start": False
        }
    }

    if "Follow" in all_tasks:
        tasks_user_data['follow']['usernames'] = input("Paste the usernames you want to follow: ").split(" ")

    if "Retweet" in all_tasks:
        tasks_user_data['retweet']['tweet_id'] = get_user_tweet_id_input("Paste your link for retweet: ")

    if "Like" in all_tasks:
        tasks_user_data['like']['tweet_id'] = get_user_tweet_id_input("Paste your link for like: ")

    tasks_user_data['tweet'] = True if "Tweet" in all_tasks else False

    tasks_user_data['tweet with picture'] = True if "Tweet with picture" in all_tasks else False

    if "Quote Tweet" in all_tasks:
        tasks_user_data['quote tweet']['tweet_id'] = input("Paste your link for quote tweet: ").strip()

    if "Comment" in all_tasks:
        tasks_user_data['comment']['tweet_id'] = get_user_tweet_id_input("Paste your link for the comment: ")

    if "Comment with picture" in all_tasks:
        tasks_user_data['comment with picture']['tweet_id'] = get_user_tweet_id_input("Paste your link for the comment with picture: ")

    tasks_user_data['change description'] = True if "Change Description" in all_tasks else False

    tasks_user_data['change username'] = True if "Change Username" in all_tasks else False

    tasks_user_data['change name'] = True if "Change Name" in all_tasks else False

    tasks_user_data['change background'] = True if "Change Background" in all_tasks else False

    tasks_user_data['change password'] = True if "Change Password" in all_tasks else False

    tasks_user_data['change birthdate'] = True if "Change Birthdate" in all_tasks else False

    tasks_user_data['change location'] = True if "Change Location" in all_tasks else False

    tasks_user_data['change profile picture'] = True if "Change profile picture" in all_tasks else False

    if "Send direct message" in all_tasks:
        tasks_user_data['send direct message']['recipient_id'] = input("Paste recipient ID: ").strip()
        tasks_user_data['send direct message']['content'] = input("Paste content text: ").strip()

    tasks_user_data['mass dm'] = True if "Mass DM" in all_tasks else False

    if "Check account messages" in all_tasks:
        tasks_user_data['check dm']['days'] = input("Enter the number of days since the last tweet: ").strip()

    tasks_user_data['check valid'] = True if "Check if account is valid" in all_tasks else False

    tasks_user_data['unfreeze'] = True if "Unfreeze Accounts" in all_tasks else False

    if "Mutual Subscription" in all_tasks:
        user_choice = int(input("Choose want you want to do (1 or 2):\n"
                                "[1] Collect usernames (required before start subscription)\n"
                                "[2] Start mutual subscription\n>> ").strip())

        if user_choice == 1:
            tasks_user_data['mutual subscription']['collect usernames'] = True
        else:
            tasks_user_data['mutual subscription']['start'] = True

    return tasks_user_data


def no_proxies() -> bool:
    user_choice = int(input("No proxies were detected. Do you want to continue without proxies? (1 or 2)\n"
                            "[1] Yes\n"
                            "[2] No\n>> ").strip())

    return True if user_choice == 1 else False


def get_user_tweet_id_input(text_to_ask: str) -> str:
    try:
        user_input = input(text_to_ask).strip()

        if "tweet_id=" in user_input:
            tweet_id = user_input.strip().split('tweet_id=')[-1].split("&")[0]
        elif "?" in user_input:
            tweet_id = user_input.strip().split('status/')[-1].split("?")[0]
        else:
            tweet_id = user_input.strip().split('status/')[-1]

        return tweet_id

    except Exception as err:
        logger.error(f"Failed to get tweet id: {err}")
        return ""
