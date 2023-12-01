from curl_cffi import requests
from loguru import logger
import json


def collect_auth_cookies(account_index: int, client: requests.Session) -> tuple[bool, str, str]:
    try:
        cookies_list = []
        auth_token = None

        for key, value in client.cookies.items():
            if value != "":
                cookies_list.append({'name': key, 'value': value})

                if key == 'auth_token':
                    auth_token = value

        cookies = str(json.dumps(cookies_list)).replace('\\"', '"').replace('""', '"')

        if auth_token is not None:
            logger.success(f"{account_index} | Grabbed auth token -> {auth_token}")
            return True, str(auth_token), cookies

        else:
            logger.error(f"{auth_token} | Failed to grab auth token.")
            return False, "", ""

    except Exception as err:
        logger.error(f"{account_index} | Failed to collect cookies: {err}")
        return False, "", ""
