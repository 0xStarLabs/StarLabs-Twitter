from curl_cffi import requests
from loguru import logger


def set_first_response_cookies(account_index: int, client: requests.Session, response: requests.Response) -> bool:
    try:
        cookies = response.headers.get_list("set-cookie")
        for cookie in cookies:
            try:
                key, value = cookie.split(';')[0].strip().split("=")
                client.cookies.set(name=key, value=value, domain=".twitter.com", path="/")

            except:
                pass

        return True

    except Exception as err:
        logger.error(f"{account_index} | Failed to set response cookies: {err}")
        return False
