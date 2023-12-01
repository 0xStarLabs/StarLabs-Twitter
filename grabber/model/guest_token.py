from curl_cffi import requests
from loguru import logger


def activate_guest_token(account_index: int, client: requests.Session) -> bool:
    try:

        resp = client.post("https://api.twitter.com/1.1/guest/activate.json",
                           headers={
                               "content-type": "application/x-www-form-urlencoded",
                               "origin": "https://twitter.com",
                               "referer": "https://twitter.com/",
                               "x-twitter-active-user": "yes"
                           })

        guest_token = resp.json()["guest_token"]

        client.headers.update({"x-guest-token": str(guest_token)})
        client.cookies.update({"gt": str(guest_token)})

        return True

    except Exception as err:
        logger.error(f"{account_index} | Failed to get 'activate' guest token: {err}")
        return False
