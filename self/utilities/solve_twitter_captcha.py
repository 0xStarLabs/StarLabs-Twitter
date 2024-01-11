from playwright.sync_api import Error as PlaywrightError
from playwright.sync_api import sync_playwright
from playwright_stealth import stealth_sync
from loguru import logger
import time

from . import TwoCaptcha


class SolveTwitterCaptcha:
    def __init__(self, auth_token: str, ct0: str, two_captcha: TwoCaptcha, account_index: int, user_agent: str):
        self.auth_token = auth_token
        self.account_index = account_index
        self.ct0 = ct0
        self.two_captcha = two_captcha
        self.user_agent = user_agent
    @staticmethod
    def wait_for_url(page, url, timeout=90):
        try:
            start_time = time.time()
            while time.time() - start_time < timeout:
                if url in page.url:
                    return True
                time.sleep(1)
            logger.error('Timeout waiting for URL')
        except:
            pass
        return False

    def wait_for_multiple_conditions(self, page, selector, url, timeout=60000) -> tuple[any, any]:
        try:
            start_time = time.time()
            while time.time() - start_time < timeout / 1000:
                element = page.query_selector(selector)
                if element:
                    return None, element

                if self.wait_for_url(page, url, 1):
                    return True, None

            return None, None

        except (PlaywrightError, TimeoutError) as e:
            return None, None

    def solve_captcha(self, proxy: str | None) -> bool:
        for _ in range(5):
            try:
                with sync_playwright() as p:
                    context_options = {
                        'user_data_dir': '',
                        'viewport': None,
                    }

                    if proxy:
                        context_options['proxy'] = {
                            "server": f"http://{proxy.split('@')[1].split(':')[0]}:{proxy.split('@')[1].split(':')[1]}",
                            "username": proxy.split('@')[0].split(':')[0],
                            "password": proxy.split('@')[0].split(':')[1],
                        }

                    context = p.firefox.launch_persistent_context(**context_options)

                    cock = [
                        {
                            "name": "auth_token",
                            "value": self.auth_token,
                            "domain": "twitter.com",
                            "path": "/",
                        },
                        {
                            "name": "ct0",
                            "value": self.ct0,
                            "domain": "twitter.com",
                            "path": "/",
                        },
                    ]

                    context.add_cookies(cock)

                    page = context.new_page()
                    stealth_sync(page)
                    page.goto('https://twitter.com/account/access')
                    page.wait_for_load_state(state='networkidle', timeout=60000)

                    home_page, element = self.wait_for_multiple_conditions(page=page,
                                                                           selector="#arkose_iframe, input[type='submit'].Button.EdgeButton.EdgeButton--primary",
                                                                           url="twitter.com/home")

                    if not home_page and not element:
                        logger.error(f'{self.account_index} | Failed to detect captcha element on the page')
                        continue

                    if home_page:
                        logger.success(f'{self.account_index} | Account successfully unfrozen')
                        return True

                    if element and element.get_attribute('value') == 'Continue to X':
                        element.click()
                        logger.success(f'{self.account_index} | Account successfully unfrozen')
                        return True

                    elif element and element.get_attribute('value') == 'Delete':
                        element.click()
                        continue

                    elif element and element.get_attribute('value') == 'Start':
                        element.click()

                        page.goto('https://twitter.com/account/access')
                        page.wait_for_selector('#arkose_iframe')

                    for _ in range(6):
                        captcha_result, ok = self.two_captcha.solve_funcaptcha('0152B4EB-D2DC-460A-89A1-629838B529C9', "https://twitter.com/account/access", self.user_agent)

                        if not ok:
                            continue

                        logger.info(f'{self.account_index} | Captcha solution received, trying to send')
                        break

                    iframe_element = page.query_selector('#arkose_iframe')
                    if not iframe_element:
                        if 'twitter.com/home' in page.url:
                            logger.success(f'{self.account_index} | Account successfully unfrozen')
                            return True

                        logger.error(f'{self.account_index} | Failed to detect captcha element on the page')
                        continue

                    iframe = iframe_element.content_frame()
                    iframe.evaluate(
                        f'parent.postMessage(JSON.stringify({{eventId:"challenge-complete",payload:{{sessionToken:"{captcha_result}"}}}}),"*")')

                    page.wait_for_load_state(state='networkidle',
                                             timeout=60000)

                    if not self.wait_for_url(page=page, url='twitter.com/home', timeout=20):
                        continue

                    logger.success(f'{self.account_index} | Account successfully unfrozen')
                    context.close()
                    return True

            except Exception as error:
                logger.error(f'{self.account_index} | Unknown error while trying to unfreeze the account: {error}')
                # print(traceback.print_exc())
                continue

        return False
