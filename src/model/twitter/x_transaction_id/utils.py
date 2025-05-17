import re
import bs4
import math
import base64
from typing import Union
from .constants import (
    MIGRATION_REDIRECTION_REGEX,
    ON_DEMAND_FILE_REGEX,
    ON_DEMAND_FILE_URL,
)


class Math:
    @staticmethod
    def round(num: Union[float, int]):
        # using javascript...? just use the native Math.round(num)
        x = math.floor(num)
        if (num - x) >= 0.5:
            x = math.ceil(num)
        return math.copysign(x, num)


def generate_headers():
    headers = {
        "Authority": "x.com",
        "Accept-Language": "en-US,en;q=0.9",
        "Cache-Control": "no-cache",
        "Referer": "https://x.com",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
        "X-Twitter-Active-User": "yes",
        "X-Twitter-Client-Language": "en",
    }
    return headers


def validate_response(response: bs4.BeautifulSoup):
    if not isinstance(response, bs4.BeautifulSoup):
        raise TypeError(
            f"the response object must be bs4.BeautifulSoup, not {type(response).__name__}"
        )


def get_migration_url(response: bs4.BeautifulSoup):
    migration_url = response.select_one("meta[http-equiv='refresh']")
    migration_redirection_url = re.search(
        MIGRATION_REDIRECTION_REGEX, str(migration_url)
    ) or re.search(MIGRATION_REDIRECTION_REGEX, str(response.content))
    return migration_redirection_url


def get_migration_form(response: bs4.BeautifulSoup):
    migration_form = response.select_one("form[name='f']") or response.select_one(
        f"form[action='https://x.com/x/migrate']"
    )
    if not migration_form:
        return
    url = migration_form.attrs.get("action", "https://x.com/x/migrate")
    method = migration_form.attrs.get("method", "POST")
    request_payload = {
        input_field.get("name"): input_field.get("value")
        for input_field in migration_form.select("input")
    }
    return {"method": method, "url": url, "data": request_payload}


def get_ondemand_file_url(response: bs4.BeautifulSoup):
    """
    Get the URL for the ondemand.s JavaScript file.

    This function tries multiple methods to find the URL:
    1. Using regex pattern matching (original method)
    2. Directly looking for script tags with ondemand.s in the src attribute
    3. Finding ondemand hash pattern in any script content
    """
    try:
        # Method 1: Using the original regex pattern
        on_demand_file = ON_DEMAND_FILE_REGEX.search(str(response))
        if on_demand_file:
            filename = on_demand_file.group(1)
            return ON_DEMAND_FILE_URL.format(filename=filename)

        # Method 2: Look for script tags with ondemand.s in the src attribute
        script_tags = response.find_all("script", src=True)
        for script in script_tags:
            src = script.get("src", "")
            if "ondemand.s" in src and ".js" in src:
                return src

        # Method 3: Find ondemand hash pattern in any script content
        # Look for patterns like ondemand.xxxx.js or dist/ondemand.xxxx.js
        scripts = response.find_all("script")
        for script in scripts:
            if script.string:
                matches = re.findall(r"(ondemand\.[a-zA-Z0-9]+\.js)", script.string)
                for match in matches:
                    hash_value = match.split(".")[1]
                    return ON_DEMAND_FILE_URL.format(filename=hash_value)

                # Also look for full URLs
                matches = re.findall(
                    r'(https?://[^"\']+ondemand\.[a-zA-Z0-9]+\.js)', script.string
                )
                if matches:
                    return matches[0]

        # If all else fails, try a more general approach
        all_text = str(response)
        matches = re.findall(r'(https?://[^"\']+ondemand\.[a-zA-Z0-9]+\.js)', all_text)
        if matches:
            return matches[0]

        return None
    except Exception as e:
        print(f"Error in get_ondemand_file_url: {e}")
        return None


def handle_x_migration(session):
    # for python requests -> session = requests.Session()
    # session.headers = generate_headers()
    response = session.request(method="GET", url="https://x.com")
    home_page = bs4.BeautifulSoup(response.content, "html.parser")
    migration_redirection_url = get_migration_url(response=home_page)
    if migration_redirection_url:
        response = session.request(method="GET", url=migration_redirection_url.group(0))
        home_page = bs4.BeautifulSoup(response.content, "html.parser")
    migration_form = get_migration_form(response=home_page)
    if migration_form:
        response = session.request(**migration_form)
        home_page = bs4.BeautifulSoup(response.content, "html.parser")
    return home_page


async def handle_x_migration_async(session):
    """
    Handle Twitter to X migration process for async HTTP clients

    Args:
        session: An authenticated httpx.AsyncClient session with cookies already set

    Returns:
        bs4.BeautifulSoup: Processed home page response
    """
    # Use the existing session with its headers and cookies (already authenticated)
    response = await session.request(method="GET", url="https://x.com")

    # If response is empty or unauthorized, the session might not be properly authenticated
    if response.status_code == 401 or not response.content:
        raise Exception(
            "Unauthorized access to X.com. Make sure session is properly authenticated."
        )

    home_page = bs4.BeautifulSoup(response.content, "html.parser")

    migration_redirection_url = get_migration_url(response=home_page)
    if migration_redirection_url:
        response = await session.request(
            method="GET", url=migration_redirection_url.group(0)
        )
        home_page = bs4.BeautifulSoup(response.content, "html.parser")

    migration_form = get_migration_form(response=home_page)
    if migration_form:
        response = await session.request(**migration_form)
        home_page = bs4.BeautifulSoup(response.content, "html.parser")

    return home_page


def float_to_hex(x):
    result = []
    quotient = int(x)
    fraction = x - quotient

    while quotient > 0:
        quotient = int(x / 16)
        remainder = int(x - (float(quotient) * 16))

        if remainder > 9:
            result.insert(0, chr(remainder + 55))
        else:
            result.insert(0, str(remainder))

        x = float(quotient)

    if fraction == 0:
        return "".join(result)

    result.append(".")

    while fraction > 0:
        fraction *= 16
        integer = int(fraction)
        fraction -= float(integer)

        if integer > 9:
            result.append(chr(integer + 55))
        else:
            result.append(str(integer))

    return "".join(result)


def is_odd(num: Union[int, float]):
    if num % 2:
        return -1.0
    return 0.0


def base64_encode(string):
    string = string.encode() if isinstance(string, str) else string
    return base64.b64encode(string).decode()
