from curl_cffi import requests


def create_client(proxy: str, user_agent: str) -> requests.Session:
    session = requests.Session(impersonate="chrome110", timeout=60)

    if proxy:
        session.proxies.update({
            "http": "http://" + proxy,
            "https": "http://" + proxy,
        })

    session.headers.update({
        'Accept-Encoding': 'gzip, deflate, br',
        'Accept': '*/*',
        'Connection': 'keep-alive',
        'User-Agent': user_agent,
        'Origin': 'https://mobile.twitter.com',
        'Referer': 'https://mobile.twitter.com/',
        'x-twitter-active-user': 'yes',
        'x-twitter-auth-type': 'OAuth2Session',
        'x-twitter-client-language': 'en',
        'content-type': 'application/x-www-form-urlencoded'
    })

    return session
