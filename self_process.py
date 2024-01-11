import random
import time
from concurrent.futures import ThreadPoolExecutor
from random import randint
import traceback
import requests
from loguru import logger
from time import sleep
import threading
import queue
import os

import self
import grabber
from self import extra


def options():
    def launch_wrapper(index, twitter, proxy):
        account_flow(lock, index, twitter, proxy, config, change_data, tasks_data, failed_queue, query_ids, usernames_queue)

    extra.show_logo()
    extra.show_dev_info()
    extra.show_menu(extra.MENU_ITEMS)

    all_tasks = extra.get_user_choice(extra.MENU_ITEMS, "Choose one or a few tasks")
    tasks_data = extra.ask_for_task_data(all_tasks)
    if tasks_data is None:
        return

    os.system("cls")
    change_data, ok = extra.get_txt_data(tasks_data)
    if not ok:
        return

    usernames_queue = queue.Queue()
    if tasks_data['mutual subscription']['collect usernames']:
        with open("data/my_usernames.txt", "w") as f:
            f.write("")
    elif tasks_data['mutual subscription']['start']:
        users = change_data['mutual subscription']
        random.shuffle(users)
        for user in users:
            usernames_queue.put(user)

    threads = int(input("How many threads do you want: ").strip())
    twitters = extra.read_txt_file("cookies", "data/accounts.txt")
    proxies = extra.read_txt_file("proxies", "data/proxies.txt")
    indexes = [i + 1 for i in range(len(twitters))]
    query_ids = extra.read_query_ids()
    mobile_proxy_queue = queue.Queue()
    config = extra.read_config()
    failed_queue = queue.Queue()
    lock = threading.Lock()

    use_proxy = True
    if len(proxies) == 0:
        if not extra.no_proxies():
            return
        else:
            use_proxy = False

    if config['mobile_proxy'].lower() == "yes":
        ip_change_links = extra.read_txt_file("ip change links", "data/ip_change_links.txt")

        for i in range(len(twitters)):
            mobile_proxy_queue.put(i)
        cycle = []
        for i in range(len(proxies)):
            data_list = (proxies[i], ip_change_links[i], mobile_proxy_queue, config, lock, twitters, change_data, tasks_data, failed_queue, query_ids, usernames_queue)
            cycle.append(data_list)

        with ThreadPoolExecutor() as executor:
            executor.map(mobile_proxy_wrapper, cycle)

    else:
        if not use_proxy:
            proxies = ["" for _ in range(len(twitters))]
        elif len(proxies) < len(twitters):
            proxies = [proxies[i % len(proxies)] for i in range(len(twitters))]

        logger.info("Starting...")
        with ThreadPoolExecutor(max_workers=threads) as executor:
            executor.map(launch_wrapper, indexes, twitters, proxies)

    failed_count = 0
    while not failed_queue.empty():
        failed_count += failed_queue.get()

    print()
    logger.info(f"STATISTICS: {len(twitters) - failed_count} SUCCESS | {failed_count} FAILED")
    input("Press Enter to exit...")


def account_flow(lock: threading.Lock, account_index: int, twitter: str, proxy: str, config: dict, txt_data: dict, tasks_data: dict, failed_queue: queue.Queue, query_ids: dict, usernames_queue: queue.Queue):
    report = False

    try:
        if len(str(twitter)) == 40 and ":" not in twitter:
            # account_type = "auth_token"
            auth_token = twitter

        elif ":" in twitter and len(twitter.split(":")) == 2:
            # account_type = "login:pass"

            grabber_instance = grabber.CookieGrabber(account_index, twitter.split(":")[0], twitter.split(":")[1], proxy, config)
            ok, auth_token, cookies = wrapper(grabber_instance.grab_cookies, config['max_tasks_retries'])
            if not ok:
                return False
            else:
                update_grabbed_cookies(twitter.split(":")[0], twitter.split(":")[1], auth_token, cookies, lock)

        elif ":" in twitter and len(twitter.split(":")[2]) == 40:
            # login:pass:auth_token:json_cookies
            auth_token = twitter.split(":")[2]

        else:
            logger.error(f"{account_index} | Wrong account format: {twitter}")
            return False

        twitter_instance = self.Twitter(account_index, auth_token, proxy, config, query_ids)

        if tasks_data['mutual subscription']['start'] or tasks_data['mutual subscription']['collect usernames']:
            if tasks_data['mutual subscription']['collect usernames']:
                with lock:
                    with open("data/my_usernames.txt", "a") as f:
                        f.write(f"{twitter_instance.username}\n")
                        logger.success(f"{account_index} | Collected username -> {twitter_instance.username}.")

            else:
                for _ in range(5):
                    user_for_sub = usernames_queue.get()

                    if twitter_instance.username != user_for_sub:
                        ok = wrapper(twitter_instance.follow, config['max_tasks_retries'], user_for_sub)
                        if not ok:
                            report = True
                        random_pause(config['pause_start'], config['pause_end'])
                        break
                    else:
                        usernames_queue.put(user_for_sub)

        if tasks_data['follow']:
            for username in tasks_data['follow']['usernames']:
                ok = wrapper(twitter_instance.follow, config['max_tasks_retries'], username)
                if not ok:
                    report = True
                random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['retweet']:
            ok = wrapper(twitter_instance.retweet, config['max_tasks_retries'], tasks_data['retweet']['tweet_id'])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['like']:
            ok = wrapper(twitter_instance.like, config['max_tasks_retries'], tasks_data['like']['tweet_id'])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['tweet']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['tweet'])
            else:
                content = txt_data['tweet'][account_index - 1]

            ok = wrapper(twitter_instance.tweet, config['max_tasks_retries'], content)
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['tweet with picture']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['tweet'])
                pic = random.choice(txt_data['pictures'])
            else:
                content = txt_data['tweet'][account_index - 1]
                pic = txt_data['pictures'][account_index - 1]

            ok = wrapper(twitter_instance.tweet_with_picture, config['max_tasks_retries'], content, pic)
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['quote tweet']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['tweet'])
            else:
                content = txt_data['tweet'][account_index - 1]

            ok = wrapper(twitter_instance.quote_tweet, config['max_tasks_retries'], content, tasks_data['quote tweet']['tweet_id'])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['comment']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['comment'])
            else:
                content = txt_data['comment'][account_index - 1]

            ok = wrapper(twitter_instance.comment, config['max_tasks_retries'], content, tasks_data['comment']['tweet_id'])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['comment with picture']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['comment'])
                pic = random.choice(txt_data['pictures'])
            else:
                content = txt_data['comment'][account_index - 1]
                pic = txt_data['pictures'][account_index - 1]

            ok = wrapper(twitter_instance.comment_with_picture, config['max_tasks_retries'], content, pic, tasks_data['comment with picture']['tweet_id'])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['change description']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['description'])
            else:
                content = txt_data['description'][account_index - 1]

            ok = wrapper(twitter_instance.change_description, config['max_tasks_retries'], content)
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['change username']:
            ok = wrapper(twitter_instance.change_username, config['max_tasks_retries'], txt_data['username'][account_index - 1])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['change name']:
            ok = wrapper(twitter_instance.change_name, config['max_tasks_retries'], txt_data['name'][account_index - 1])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['change background']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['pictures'])
            else:
                content = txt_data['pictures'][account_index - 1]

            ok = wrapper(twitter_instance.change_background, config['max_tasks_retries'], content)
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['change password']:
            ok = wrapper(twitter_instance.change_password, config['max_tasks_retries'], txt_data['current password'][account_index - 1], txt_data['new password'][account_index - 1])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['change birthdate']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['birthdate'])
            else:
                content = txt_data['birthdate'][account_index - 1]

            ok = wrapper(twitter_instance.change_birthdate, config['max_tasks_retries'], content)
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['change location']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['location'])
            else:
                content = txt_data['location'][account_index - 1]

            ok = wrapper(twitter_instance.change_location, config['max_tasks_retries'], content)
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['change profile picture']:
            if config['data_random'] == "yes":
                content = random.choice(txt_data['pictures'])
            else:
                content = txt_data['pictures'][account_index - 1]

            ok = wrapper(twitter_instance.change_profile_picture, config['max_tasks_retries'], content)
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['send direct message']:
            ok = wrapper(twitter_instance.send_direct_message, config['max_tasks_retries'], tasks_data['send direct message']['content'], tasks_data['send direct message']['recipient_id'])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['check dm']:
            ok = wrapper(twitter_instance.check_dm, config['max_tasks_retries'], tasks_data['check dm']['days'])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['check valid']:
            ok, status = wrapper(twitter_instance.check_suspended, config['max_tasks_retries'])
            if not ok:
                report = True
            else:
                if status != "ban":
                    with lock:
                        with open("data/valid_accounts.txt", "a") as f:
                            f.write(f"{twitter_instance.auth_token}\n")
                else:
                    report_locked_twitter(auth_token, lock)

            random_pause(config['pause_start'], config['pause_end'])

        if tasks_data['unfreeze']:
            ok = wrapper(twitter_instance.start_unfreezer, config['max_tasks_retries'])
            if not ok:
                report = True
            random_pause(config['pause_start'], config['pause_end'])

    except Exception as err:
        logger.error(f"{account_index} | Account flow failed: {err}")
        report = True

    if report:
        report_failed_twitter(twitter, lock, failed_queue)


def wrapper(function, attempts: int, *args, **kwargs):
    for _ in range(attempts):
        result = function(*args, **kwargs)
        if isinstance(result, tuple) and result and isinstance(result[0], bool):
            if result[0]:
                return result
        elif isinstance(result, bool):
            if result:
                return True

    return result


def report_failed_twitter(twitter: str, lock: threading.Lock, failed_queue: queue.Queue):
    try:
        with lock:
            with open("data/failed_accounts.txt", "a") as file:
                file.write(twitter + "\n")
                failed_queue.put(1)
                return

    except Exception as err:
        logger.error(f"Error while reporting failed account: {err}")


def random_pause(start: int, end: int):
    sleep(randint(start, end))


def report_locked_twitter(twitter: str, lock: threading.Lock):
    try:
        with lock:
            with open("data/locked_tokens.txt", "a") as file:
                file.write(twitter + "\n")
                return

    except Exception as err:
        logger.error(f"Error while reporting locked token: {err}")


def update_grabbed_cookies(login: str, password: str, auth_token: str, cookies: str, lock: threading.Lock):
    try:
        with lock:
            with open("data/accounts.txt", "r") as file:
                lines = file.readlines()

            updated_lines = []
            for line in lines:
                if line.startswith(login + ":"):
                    updated_lines.append(f"{login}:{password}:{auth_token}:{cookies}\n")
                else:
                    updated_lines.append(line)

            with open("data/accounts.txt", "w") as file:
                file.writelines(updated_lines)

    except Exception as err:
        logger.error(f"Error while saving grabbed cookies: {err}")


def mobile_proxy_wrapper(data):
    # proxies[i], ip_change_links[i], mobile_proxy_queue, config, lock, twitters, change_data, tasks_data, failed_queue, query_ids
    proxy, ip_change_link, q, config, lock, twitters, change_data, tasks_data, failed_queue, query_ids, usernames_queue = data[:11]

    while not q.empty():
        i = q.get()

        try:
            for _ in range(3):
                try:
                    requests.get(f"{ip_change_link}",
                                 headers={"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"},
                                 timeout=60)

                    time.sleep(config['change_ip_pause'])
                    logger.success(f"{i + 1} | Successfully changed IP")
                    break

                except Exception as err:
                    logger.error(f"{i + 1} | Mobile proxy error! Check your ip change link: {err}")
                    time.sleep(2)

            account_flow(lock, i + 1, twitters[i], proxy, config, change_data, tasks_data, failed_queue, query_ids, usernames_queue)

        except Exception as err:
            logger.error(f"{i + 1} | Mobile proxy flow error: {err}")
