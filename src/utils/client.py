from curl_cffi.requests import AsyncSession
import secrets

HEADERS = {
    "accept": "*/*",
    "accept-language": "en-GB,en-US;q=0.9,en;q=0.8,ru;q=0.7,zh-TW;q=0.6,zh;q=0.5",
    "content-type": "application/json",
    "priority": "u=1, i",
    "sec-ch-ua": '"Google Chrome";v="133", "Chromium";v="133", "Not_A Brand";v="24"',
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": '"Windows"',
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "same-site",
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
}


async def create_twitter_client(
    proxy: str, auth_token: str, verify_ssl: bool = True
) -> tuple[AsyncSession, str]:
    session = AsyncSession(
        impersonate="chrome131",
        verify=verify_ssl,
        timeout=60,
    )

    if proxy:
        session.proxies.update({
            "http": "http://" + proxy,
            "https": "http://" + proxy,
        })

    generated_csrf_token = secrets.token_hex(16)

    cookies = {"ct0": generated_csrf_token, "auth_token": auth_token}
    headers = {"x-csrf-token": generated_csrf_token}

    session.headers.update(headers)
    session.cookies.update(cookies)

    session.headers["x-csrf-token"] = generated_csrf_token

    session.headers = get_headers(session)

    return session, generated_csrf_token


def get_headers(session: AsyncSession, **kwargs) -> dict:
    """
    Get the headers required for authenticated requests
    """
    cookies = session.cookies

    headers = kwargs | {
        "authorization": "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs=1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
        # "cookie": "; ".join(f"{k}={v}" for k, v in cookies.items()),
        "origin": "https://x.com",
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
        "x-csrf-token": cookies.get("ct0", ""),
        # "x-guest-token": cookies.get("guest_token", ""),
        "x-twitter-auth-type": "OAuth2Session",
        "x-twitter-active-user": "yes",
        "x-twitter-client-language": "en",
    }
    return dict(sorted({k.lower(): v for k, v in headers.items()}.items()))
