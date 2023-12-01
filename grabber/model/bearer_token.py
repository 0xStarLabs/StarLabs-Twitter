from curl_cffi import requests
from loguru import logger


def set_bearer_token(account_index: int, client: requests.Session) -> bool:
    for i in range(3):
        try:
            resp = client.get("https://abs.twimg.com/responsive-web/client-web/main.f3ada2b5.js")
            bearer_token = "Bearer " + resp.text.split('const r="ACTION_FLUSH",i="ACTION_REFRESH')[1].split(',l="d_prefs"')[0].split(',s="')[1].split('"')[0]

            client.headers.update({
                "Authorization": bearer_token
            })

            return True

        except Exception as err:
            logger.error(f"{account_index} | Failed to get bearer token: {err}")

    return False
