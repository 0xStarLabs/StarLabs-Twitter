# quick api for the https://capmonster.cloud/ captcha service #
import requests as default_requests
from curl_cffi import requests
from loguru import logger
from time import sleep
from typing import Any


class Capmonstercloud:
    def __init__(self, account_index: int, api_key: str, client: requests.Session, proxy: str):
        self.account_index = account_index
        self.api_key = api_key
        self.client = client
        self.proxy = proxy

    # returns gRecaptchaResponse, resp_key and success status (False if something failed)
    def solve_hcaptcha(self, site_key: str, website_url: str, captcha_rqdata: str, user_agent: str) -> tuple[str, bool]:
        try:
            if self.client.proxies:
                captcha_data = {
                    "clientKey": self.api_key,
                    "task":
                        {
                            "type": "HCaptchaTask",
                            "websiteURL": website_url,
                            "websiteKey": site_key,
                            "data": captcha_rqdata,
                            "userAgent": user_agent,
                            "fallbackToActualUA": True,
                            "proxyType": "http",
                            "proxyAddress": self.proxy.split("@")[1].split(":")[0],
                            "proxyPort": self.proxy.split("@")[1].split(":")[1],
                            "proxyLogin": self.proxy.split("@")[0].split(":")[0],
                            "proxyPassword": self.proxy.split("@")[0].split(":")[1]
                        }
                }

            else:
                captcha_data = {
                    "clientKey": self.api_key,
                    "task":
                        {
                            "type": "HCaptchaTaskProxyless",
                            "websiteURL": website_url,
                            "websiteKey": site_key,
                            "data": captcha_rqdata,
                            "userAgent": user_agent,
                            "fallbackToActualUA": True,
                        }
                }

            resp = self.client.post("https://api.capmonster.cloud/createTask",
                                    json=captcha_data)

            if resp.json()["errorId"] == 0:
                logger.info(f"{self.account_index} | Starting to solve the captcha...")
                return self.get_captcha_result(resp.json()["taskId"])

            else:
                logger.error(f"{self.account_index} | Failed to send captcha request: {resp.json()['errorDescription']}")
                return "", False

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to send captcha request to capmonster: {err}")
            return "", False

    # returns gRecaptchaResponse, resp_key and success status (False if something failed)
    def solve_turnstile(self, site_key: str, website_url: str) -> tuple[str, bool]:
        try:
            captcha_data = {
                "clientKey": self.api_key,
                "task":
                    {
                        "type": "TurnstileTaskProxyless",
                        "websiteURL": website_url,
                        "websiteKey": site_key,
                    }
            }

            resp = self.client.post("https://api.capmonster.cloud/createTask",
                                    json=captcha_data)

            if resp.json()["errorId"] == 0:
                logger.info(f"{self.account_index} | Starting to solve the CF captcha...")
                return self.get_captcha_result(resp.json()["taskId"], "cf")

            else:
                logger.error(f"{self.account_index} | Failed to send CF captcha request: {resp.json()['errorDescription']}")
                return "", False

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to send CF captcha request to capmonster: {err}")
            return "", False

    def solve_funcaptcha(self, site_key: str, website_url: str) -> tuple[str, bool]:
        try:
            captcha_data = {
                "clientKey": self.api_key,
                "task":
                    {
                        "type": "FunCaptchaTaskProxyless",
                        "websiteURL": website_url,
                        "websitePublicKey": site_key
                    }
            }

            resp = self.client.post("https://api.capmonster.cloud/createTask",
                                    json=captcha_data)

            if resp.json()["errorId"] == 0:
                logger.info(f"{self.account_index} | Starting to solve the funcaptcha...")
                return self.get_captcha_result(resp.json()["taskId"], "funcaptcha")

            else:
                logger.error(f"{self.account_index} | Failed to send funcaptcha request: {resp.json()['errorDescription']}")
                return "", False

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to send funcaptcha request to capmonster: {err}")
            return "", False

    def get_captcha_result(self, task_id: str, captcha_type: str = "hcaptcha") -> tuple[Any, bool] | tuple[str, str, bool]:
        for i in range(30):
            try:
                resp = self.client.post("https://api.capmonster.cloud/getTaskResult/",
                                        json={
                                            "clientKey": self.api_key,
                                            "taskId": int(task_id)
                                        })

                if resp.json()["errorId"] == 0:
                    if resp.json()["status"] == "ready":

                        logger.success(f"{self.account_index} | Captcha solved!")

                        if captcha_type == "cf" or "funcaptcha":
                            g_recaptcha_response = resp.json()['solution']["token"]

                        elif captcha_type == "image_to_text":
                            g_recaptcha_response = resp.json()['solution']["text"]

                        else:
                            g_recaptcha_response = resp.json()['solution']["gRecaptchaResponse"]

                        return g_recaptcha_response, True
                elif resp.json()['errorId'] == 1:
                    if "ERROR_CAPTCHA_UNSOLVABLE" in resp.text:
                        raise Exception("ERROR CAPTCHA UNSOLVABLE")
                else:
                    pass

            except Exception as err:
                logger.error(f"{self.account_index} | Failed to get captcha solution: {err}")
                return "", False
            # sleep between result requests
            sleep(6)

        logger.error(f"{self.account_index} | Failed to get captcha solution")
        return "", False


class TwoCaptcha:
    def __init__(self, account_index: int, api_key: str, client: requests.Session, proxy: str):
        self.account_index = account_index
        self.api_key = api_key
        self.client = client
        self.proxy = proxy

    # returns gRecaptchaResponse, resp_key and success status (False if something failed)
    def solve_image_to_text(self, image_base64: str, captcha_options: {}) -> tuple[str, bool]:
        try:
            captcha_data = {
                "key": self.api_key,
                "method": "base64",
                "body": str(image_base64),
                "json": 1
            }

            captcha_data.update(captcha_options)

            resp = default_requests.post("http://2captcha.com/in.php",
                                         data=captcha_data)

            if resp.json()["status"] == 1:
                logger.info(f"{self.account_index} | Starting to solve the captcha...")
                return self.get_captcha_result(resp.json()["request"])

            else:
                logger.error(f"{self.account_index} | Failed to send captcha request: {resp.json()['errorDescription']}")
                return "", False

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to send captcha request to 2captcha: {err}")
            return "", False

    def solve_funcaptcha(self, site_key: str, website_url: str, user_agent: str) -> tuple[str, bool]:
        try:
            captcha_data = {
                "key": self.api_key,
                "method": "funcaptcha",
                "publickey": site_key,
                "surl": "https://client-api.arkoselabs.com",
                "pageurl": website_url,
                "userAgent": user_agent,
                "json": 1
            }

            resp = default_requests.post("http://2captcha.com/in.php",
                                         data=captcha_data)

            if resp.json()["status"] == 1:
                logger.info(f"{self.account_index} | Starting to solve funcaptcha...")
                return self.get_captcha_result(resp.json()["request"])

            else:
                logger.error(f"{self.account_index} | Failed to send funcaptcha request: {resp.json()['errorDescription']}")
                return "", False

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to send funcaptcha request to 2captcha: {err}")
            return "", False

    def get_captcha_result(self, task_id: str) -> tuple[Any, bool] | tuple[str, bool]:
        for i in range(50):
            try:
                resp = default_requests.post("http://2captcha.com/res.php",
                                             params={
                                                 "key": self.api_key,
                                                 "action": "get",
                                                 "id": int(task_id),
                                                 "json": 1
                                             })

                if resp.json()["status"] == 1:
                    logger.success(f"{self.account_index} | Captcha solved!")

                    response = resp.json()['request']
                    return response, True

            except Exception as err:
                logger.error(f"{self.account_index} | Failed to get captcha solution: {err}")
                return "", False
            # sleep between result requests
            sleep(5)

        logger.error(f"{self.account_index} | Failed to get captcha solution")
        return "", False


class OneCaptcha:
    def __init__(self, account_index: int, api_key: str, client: requests.Session):
        self.account_index = account_index
        self.api_key = api_key
        self.client = client

    def solve_funcaptcha(self, site_key: str, site_url: str) -> str:
        try:
            captcha_data = {
                'apikey': self.api_key,
                'sitekey': site_key,
                'siteurl': site_url,
                'affiliateid': 33644
            }

            resp = self.client.get("https://api.1stcaptcha.com/funcaptchatokentask",
                                   params=captcha_data)

            if resp.json()["Code"] == 0:
                logger.info(f"{self.account_index} | Starting to solve the funcaptcha...")
                return self.get_captcha_result(resp.json()["TaskId"])

            else:
                logger.error(f"{self.account_index} | Failed to send funcaptcha request: {resp.json()['Message']}")
                return ""

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to send funcaptcha request to capmonster: {err}")
            return ""

    def get_captcha_result(self, task_id: int) -> str:
        for i in range(30):
            try:
                resp = self.client.get("https://api.1stcaptcha.com/getresult",
                                       params={
                                           "apikey": self.api_key,
                                           "taskId": int(task_id)
                                       })

                if resp.json()["Code"] == 0:
                    if resp.json()["Status"] == "SUCCESS":
                        logger.success(f"{self.account_index} | Captcha solved!")

                        return resp.json()['Data']["Token"]

                    elif resp.json()["Status"] == "ERROR":
                        raise Exception(resp.json()['Message'])

                else:
                    raise Exception(resp.text)

            except Exception as err:
                logger.error(f"{self.account_index} | Failed to get captcha solution: {err}")
                return ""
            # sleep between result requests
            sleep(6)

        logger.error(f"{self.account_index} | Failed to get captcha solution")
        return ""


class Capsolver:
    def __init__(self, account_index: int, api_key: str, client: requests.Session, proxy: str):
        self.account_index = account_index
        self.api_key = api_key
        self.client = client
        self.proxy = proxy

    def solve_funcaptcha(self, site_key: str, url: str) -> tuple[str, bool]:
        try:
            if self.proxy != "":
                captcha_data = {
                    "clientKey": self.api_key,
                    "task": {
                        "type": "FunCaptchaTaskProxyLess",
                        "websiteURL": url,
                        "websiteKey": site_key,
                    }
                }
            else:
                logger.error(f"{self.account_index} | Capsolver requires proxy for solving FunCaptcha.")
                return "", False

            resp = default_requests.post("https://api.capsolver.com/createTask",
                                         json=captcha_data)

            if resp.status_code == 200:
                logger.info(f"{self.account_index} | Starting to solve FunCaptcha...")
                return self.get_captcha_result(resp.json()["taskId"])

            else:
                logger.error(f"{self.account_index} | Failed to send FunCaptcha request: {resp.json()['errorDescription']}")
                return "", False

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to send FunCaptcha request to Capsolver: {err}")
            return "", False

    def get_captcha_result(self, task_id: str) -> tuple[Any, bool] | tuple[str, bool]:
        for i in range(30):
            try:
                resp = default_requests.post("https://api.capsolver.com/getTaskResult",
                                             json={
                                                 "clientKey": self.api_key,
                                                 "taskId": task_id
                                             })

                if resp.status_code == 200:
                    if resp.json()['errorId'] != 0:
                        logger.error(f"{self.account_index} | FunCaptcha failed!")
                        return "", False

                    elif resp.json()['status'] == "ready":
                        logger.success(f"{self.account_index} | FunCaptcha solved!")

                        response = resp.json()['solution']['token']
                        return response, True

            except Exception as err:
                logger.error(f"{self.account_index} | Failed to get FunCaptcha solution: {err}")
                return "", False
            # sleep between result requests
            sleep(1)

        logger.error(f"{self.account_index} | Failed to get FunCaptcha solution")
        return "", False

