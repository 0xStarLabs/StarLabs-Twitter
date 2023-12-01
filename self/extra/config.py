from configparser import ConfigParser


# read config.ini
def read_config() -> dict:
    settings = {}
    config = ConfigParser()
    config.read('config.ini')

    settings['max_tasks_retries'] = int(config['info']['max_tasks_retries'])
    settings["1stcaptcha_api_key"] = str(config['info']['1stcaptcha_api_key'])
    settings["auto_unfreeze"] = str(config['info']['auto_unfreeze']).lower().strip()

    pause_between_tasks = config['info']['pause_between_tasks']
    settings['pause_start'] = int(pause_between_tasks.split("-")[0])
    settings['pause_end'] = int(pause_between_tasks.split("-")[1])

    return settings
