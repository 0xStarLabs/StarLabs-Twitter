from configparser import ConfigParser


# read config.ini
def read_config() -> dict:
    settings = {}
    config = ConfigParser()
    config.read('config.ini')

    settings['max_tasks_retries'] = int(config['info']['max_tasks_retries'])
    settings["2captcha_api_key"] = str(config['info']['2captcha_api_key'])
    settings["auto_unfreeze"] = str(config['info']['auto_unfreeze']).lower().strip()

    pause_between_tasks = config['info']['pause_between_tasks']
    settings['pause_start'] = int(pause_between_tasks.split("-")[0])
    settings['pause_end'] = int(pause_between_tasks.split("-")[1])

    settings["mobile_proxy"] = str(config['proxy']['mobile_proxy']).strip().lower()
    settings["change_ip_pause"] = int(config['proxy']['change_ip_pause'].strip())

    settings["data_random"] = str(config['data']['random']).strip().lower()

    return settings


def read_query_ids() -> dict:
    settings = {}
    config = ConfigParser()
    config.read('query_ids.ini')

    settings['bearer_token'] = str(config['info']['bearer_token'])
    settings["like"] = str(config['info']['like'])
    settings["retweet"] = str(config['info']['retweet']).lower().strip()
    settings['tweet'] = str(config['info']['tweet'])
    settings['comment'] = str(config['info']['comment'])

    return settings
