from curl_cffi import requests
from loguru import logger

from grabber.model import set_first_response_cookies


def make_task_request(account_index: int, client: requests.Session, params: dict, json_data: dict) -> tuple[str, bool]:
    try:

        resp = client.post("https://api.twitter.com/1.1/onboarding/task.json",
                           headers={
                               "Content-Type": "application/json",
                               "Origin": "https://twitter.com",
                               "Referer": "https://twitter.com/"
                           },
                           params=params,
                           json=json_data)

        if "we've sent a confirmation code to" in resp.text:
            logger.warning(f"{account_index} | Twitter sent a confirmation code to the mail.")
            return "", False

        set_first_response_cookies(account_index, client, resp)

        try:
            if resp.json()["status"] == "success":
                flow_token = resp.json()["flow_token"]
                return flow_token, True
            else:
                return "", False
        except:
            logger.error(f"{account_index} | Failed to make a task: {resp.json()['errors'][0]['message']}")
            return "", False

    except Exception as err:
        logger.error(f"{account_index} | Failed to make a task: {err}")
        return "", False


def generate_first_task() -> dict:
    return {
        'input_flow_data': {
            'flow_context': {
                'debug_overrides': {},
                'start_location': {
                    'location': 'manual_link',
                },
            },
        },
        'subtask_versions': {
            'action_list': 2,
            'alert_dialog': 1,
            'app_download_cta': 1,
            'check_logged_in_account': 1,
            'choice_selection': 3,
            'contacts_live_sync_permission_prompt': 0,
            'cta': 7,
            'email_verification': 2,
            'end_flow': 1,
            'enter_date': 1,
            'enter_email': 2,
            'enter_password': 5,
            'enter_phone': 2,
            'enter_recaptcha': 1,
            'enter_text': 5,
            'enter_username': 2,
            'generic_urt': 3,
            'in_app_notification': 1,
            'interest_picker': 3,
            'js_instrumentation': 1,
            'menu_dialog': 1,
            'notifications_permission_prompt': 2,
            'open_account': 2,
            'open_home_timeline': 1,
            'open_link': 1,
            'phone_verification': 4,
            'privacy_options': 1,
            'security_key': 3,
            'select_avatar': 4,
            'select_banner': 2,
            'settings_list': 7,
            'show_code': 1,
            'sign_up': 2,
            'sign_up_review': 4,
            'tweet_selection_urt': 1,
            'update_users': 1,
            'upload_media': 1,
            'user_recommendations_list': 4,
            'user_recommendations_urt': 1,
            'wait_spinner': 3,
            'web_modal': 1,
        },
    }


def generate_second_task(flow_token: str) -> dict:
    return {
        'flow_token': flow_token,
        'subtask_inputs': [
            {
                'subtask_id': 'LoginJsInstrumentationSubtask',
                'js_instrumentation': {
                    'response': '{"rf":{"a41eea33ff45860584e0a227cca4d3b195dc4c024f64f5e07f7f08c1f0045773":-1,"a3fbe732c22c0e16bd3424dd41be249e363d8bf64e46c5349aba728a1d99001b":95,"aff908271b9a8bc0140fcc1a221cb3bd6aafaab0ee4e8e246791a17eeda505da":234,"dcbb29c850b8d22588a377fafa2d2f742c38c6ecf20874c8b14e6b62e0473f77":180},"s":"YrwdUofQdV3MssGow7mi0ew0No0tgY57AdhypjW5h75VBvTMcWT5ClJmAXJO33usTnzi5-PxsykcIZ5IgJCnSJB7BUf2FLq946udjqKjT_CglZluS-if5E8nm9DPuHuLwF2E76XttS9lpxVST_eFfnl96EtadnGTxuLa471S0_1LcXxdrqXQpqkObHYNchGymzxTT1zc-r9KkBo800CXo5S32s0pE9eRZuBjOPGiTK0QHLPGZsrhOLTHUrs9yhJQPKfUAGz7zEIZe0sZVHggFZWpsuc9ozoD-Omxfpw3JZqiZ-JovhnPyWICzz4S58SXYbQBmfH_btRkvXM8Ad0uVgAAAYuP1-lw"}',
                    'link': 'next_link',
                },
            },
        ],
    }


def generate_third_task(flow_token: str, twitter_login: str) -> dict:
    return {
        'flow_token': flow_token,
        'subtask_inputs': [
            {
                'subtask_id': 'LoginEnterUserIdentifierSSO',
                'settings_list': {
                    'setting_responses': [
                        {
                            'key': 'user_identifier',
                            'response_data': {
                                'text_data': {
                                    'result': twitter_login,
                                },
                            },
                        },
                    ],
                    'link': 'next_link',
                },
            },
        ],
    }


def generate_fourth_task(flow_token: str, twitter_password: str) -> dict:
    return {
        'flow_token': flow_token,
        'subtask_inputs': [
            {
                'subtask_id': 'LoginEnterPassword',
                'enter_password': {
                    'password': twitter_password,
                    'link': 'next_link',
                },
            },
        ],
    }


def generate_fifth_task(flow_token: str) -> dict:
    return {
        'flow_token': flow_token,
        'subtask_inputs': [
            {
                'subtask_id': 'AccountDuplicationCheck',
                'check_logged_in_account': {
                    'link': 'AccountDuplicationCheck_false',
                },
            },
        ],
    }
