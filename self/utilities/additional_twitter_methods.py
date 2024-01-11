from curl_cffi import requests
from random import randbytes
from loguru import logger
from hashlib import md5
import json


def set_auth_cookies(account_index: int, client: requests.Session, twitter_auth: str, bearer_token: str) -> str:
    try:
        bearer_token = "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"
        # json cookies
        if '[' in twitter_auth and ']' in twitter_auth:
            json_part = twitter_auth.split('[')[-1].split(']')[0]

            try:
                cookies_json = json.loads(json_part)
            except json.JSONDecodeError:
                logger.error(f"{account_index} | Failed to decode account json cookies.")
                return ""

            auth_token = ""
            for current_cookie in cookies_json:
                if 'name' in current_cookie and 'value' in current_cookie:
                    if 'domain' not in current_cookie or current_cookie['domain']:
                        client.cookies.set(current_cookie['name'], current_cookie['value'])

                        if current_cookie['name'] == 'ct0':
                            csrf_token = current_cookie['value']

                        if current_cookie['name'] == 'auth_token':
                            auth_token = current_cookie['value']

            client.headers.update({'authorization': bearer_token, 'x-csrf-token': csrf_token})
            return auth_token

        # auth_token
        elif len(twitter_auth) < 60:
            csrf_token = md5(randbytes(32)).hexdigest()
            cookie_str = f"des_opt_in=Y; auth_token={twitter_auth}; ct0={csrf_token};"
            client.headers.update({'cookie': cookie_str})
            client.headers.update({'authorization': bearer_token, 'x-csrf-token': csrf_token})
            return twitter_auth

        # requests browser cookies
        else:
            client.headers.update({'cookie': twitter_auth})
            csrf_token = twitter_auth.split('ct0=')[-1].split(';')[0]
            client.headers.update({'authorization': bearer_token, 'x-csrf-token': csrf_token})
            return twitter_auth.split('auth_token=')[-1].split(';')[0]

    except Exception as err:
        logger.error(f'{account_index} | Failed to set headers: {err}')
        return ""


def get_account_username(account_index: int, client: requests.Session) -> tuple[str, str]:
    try:
        username_url = 'https://mobile.twitter.com/i/api/1.1/account/settings.json?include_mention_filter=true&include_nsfw_user_flag=true&include_nsfw_admin_flag=true&include_ranked_timeline=true&include_alt_text_compose=true&ext=ssoConnections&include_country_code=true&include_ext_dm_nsfw_media_filter=true&include_ext_sharing_audiospaces_listening_data_with_followers=true'

        for x in range(4):
            resp = client.get(username_url, verify=False, timeout=120)

            ct0 = resp.cookies.get("ct0")

            if resp.json().get('screen_name') is not None:
                return resp.json()['screen_name'], ct0

            elif "this account is temporarily locked" in resp.text:
                logger.error(f"{account_index} | This account is temporarily locked.")
                return "locked", ct0

            elif "Could not authenticate you" in resp.text:
                raise Exception("Could not authenticate you.")

            else:
                raise Exception(resp.text)

    except Exception as err:
        logger.error(f"{account_index} | Failed to get account username: {err}")
        return "", ""


def mass_dm():
    pass


def set_ct0_response_token(resp, client: requests.Session):
    ct0 = resp.json().get("ct0")
    client.cookies.update({"ct0": ct0})
    client.headers.update({"x-csrf-token": ct0})


