from curl_cffi import requests
from loguru import logger

from grabber import utils
from grabber import model


class CookieGrabber:
    def __init__(self, account_index: int, twitter_login: str, twitter_password: str, proxy: str, config: dict, twitter_email: str = None):
        self.account_index = account_index
        self.twitter_login = twitter_login
        self.twitter_password = twitter_password
        self.twitter_email = twitter_email
        self.proxy = proxy
        self.config = config

        self.client: requests.Session | None = None
        self.flow_token: str = ""

    def login(self):
        self.client = utils.create_client(self.proxy)

        resp = self.client.get("https://twitter.com/i/flow/login")

        if not model.set_first_response_cookies(self.account_index, self.client, resp):
            return False

        if not model.set_bearer_token(self.account_index, self.client):
            return False

        if not model.activate_guest_token(self.account_index, self.client):
            return False

        ok = self.make_login_tasks()
        if not ok:
            return False

        logger.success(f"{self.account_index} | Logged into account {self.twitter_login}")
        return True

    def make_login_tasks(self) -> bool:
        self.flow_token, ok = model.make_task_request(self.account_index,
                                                      self.client,
                                                      params={'flow_name': 'login'},
                                                      json_data=model.generate_first_task())
        if not ok:
            return False

        self.flow_token, ok = model.make_task_request(self.account_index,
                                                      self.client,
                                                      params={},
                                                      json_data=model.generate_second_task(self.flow_token))
        if not ok:
            return False

        self.flow_token, ok = model.make_task_request(self.account_index,
                                                      self.client,
                                                      params={},
                                                      json_data=model.generate_third_task(self.flow_token, self.twitter_login))
        if not ok:
            return False

        self.flow_token, ok = model.make_task_request(self.account_index,
                                                      self.client,
                                                      params={},
                                                      json_data=model.generate_fourth_task(self.flow_token, self.twitter_password))
        if not ok:
            return False

        self.flow_token, ok = model.make_task_request(self.account_index,
                                                      self.client,
                                                      params={},
                                                      json_data=model.generate_fifth_task(self.flow_token))

        if not ok:
            return False

        return True

    def grab_cookies(self) -> tuple[bool, str, str]:
        ok = self.login()
        if not ok:
            return False, "", ""

        return model.collect_auth_cookies(self.account_index, self.client)
