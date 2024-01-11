from curl_cffi import requests
from datetime import datetime
from loguru import logger

from . import utilities


class Twitter:
    def __init__(self, account_index: int, twitter_auth: str, proxy: str, config: dict, query_ids: dict):

        self.account_index = account_index
        self.twitter_auth = twitter_auth
        self.proxy = proxy
        self.config = config
        self.query_ids = query_ids

        self.client: requests.Session | None = None

        self.auth_token: str = ""
        self.ct0: str = ""
        self.username: str | None = None
        self.user_agent: str = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36'

        self.__prepare_data()

    def __prepare_data(self):
        for x in range(5):
            try:
                # create client
                self.client = utilities.create_client(self.proxy, self.user_agent)
                # set headers
                self.auth_token = utilities.set_auth_cookies(self.account_index, self.client, self.twitter_auth, self.query_ids['bearer_token'])
                if self.auth_token == "":
                    x = 4
                    raise Exception("...")

                self.username, self.ct0 = utilities.get_account_username(self.account_index, self.client)

                if self.username == "locked":
                    if self.config['auto_unfreeze'] == "yes":
                        self._unfreeze()
                    else:
                        raise Exception("account blocked.")

                elif self.username == "" or self.ct0 == "":
                    x = 4
                    raise Exception("...")

                self.client.headers.update({"x-csrf-token": self.ct0})
                self.client.cookies.update({"ct0": self.ct0})

                break
            except Exception as err:
                if x + 1 != 5:
                    logger.error(f"{self.account_index} | Failed to prepare data: {err}")
                else:
                    raise Exception(f"Failed to prepare data 5 times. Exit...")

    def like(self, tweet_id: str) -> bool:
        try:
            like_url = f'https://twitter.com/i/api/graphql/{self.query_ids["like"]}/FavoriteTweet'
            json_data = {
                'variables': {
                    'tweet_id': tweet_id,
                },
                'queryId': self.query_ids['like'],
            }

            like_headers = {'content-type': 'application/json'}
            resp = self.client.post(like_url,
                                    headers=like_headers,
                                    json=json_data)

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif "has already favorited tweet" in resp.text or '"favorite_tweet":"Done"' in resp.text:
                logger.success(f"{self.account_index} | Liked.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to like: {err}")
            return False

    def retweet(self, tweet_id: str) -> bool:
        try:
            retweet_url = 'https://twitter.com/i/api/graphql/' + self.query_ids['retweet'] + '/CreateRetweet'
            json_data = {
                'variables': {
                    'tweet_id': tweet_id,
                    'dark_request': False,
                },
                'queryId': self.query_ids['retweet'],
            }
            retweet_headers = {'content-type': 'application/json'}

            resp = self.client.post(retweet_url,
                                    headers=retweet_headers,
                                    json=json_data,
                                    verify=False)

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif "You have already retweeted this Tweet" in resp.text or "create_retweet" in resp.text:
                logger.success(f"{self.account_index} | Retweeted.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to retweet: {err}")
            return False

    def follow(self, user_to_follow: str) -> bool:
        try:
            resp = self.client.get(f"https://twitter.com/i/api/graphql/G3KGOASz96M-Qu0nwmGXNg/UserByScreenName",
                                   params={
                                       'variables': '{"screen_name":"' + user_to_follow + '","withSafetyModeUserFields":true}',
                                       'features': '{"hidden_profile_likes_enabled":true,"hidden_profile_subscriptions_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"subscriptions_verification_info_is_identity_verified_enabled":true,"subscriptions_verification_info_verified_since_enabled":true,"highlights_tweets_tab_ui_enabled":true,"creator_subscriptions_tweet_preview_api_enabled":true,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true}',
                                       'fieldToggles': '{"withAuxiliaryUserLabels":false}',
                                   }
                                   )

            rest_id = (resp.json().get('data', {})
                       .get('user', {})
                       .get('result', {})
                       .get('rest_id'))

            if rest_id is None:
                logger.error(f"{self.account_index} | Failed to get rest ID: {resp.text}")

            resp = self.client.post("https://twitter.com/i/api/1.1/friendships/create.json",
                                    data={
                                        'include_profile_interstitial_type': '1',
                                        'include_blocking': '1',
                                        'include_blocked_by': '1',
                                        'include_followed_by': '1',
                                        'include_want_retweets': '1',
                                        'include_mute_edge': '1',
                                        'include_can_dm': '1',
                                        'include_can_media_tag': '1',
                                        'include_ext_has_nft_avatar': '1',
                                        'include_ext_is_blue_verified': '1',
                                        'include_ext_verified_type': '1',
                                        'include_ext_profile_image_shape': '1',
                                        'skip_status': '1',
                                        'user_id': rest_id,
                                    }
                                    )

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif "You can't follow yourself." in resp.text:
                logger.error(f"{self.account_index} | You can't follow yourself.")
            elif "Your account is suspended" in resp.text:
                logger.error(f"{self.account_index} | Your account is suspended.")

            resp_data = resp.json()
            if resp_data.get('id') == int(rest_id):
                logger.success(f"{self.account_index} | Subscribed to {user_to_follow}")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to subscribe: {err}")
            return False

    def tweet(self, tweet_content: str) -> bool:
        try:
            tweet_url = f'https://twitter.com/i/api/graphql/{self.query_ids["tweet"]}/CreateTweet'
            tweet_json = {
                "variables": "{\"tweet_text\":\"" + tweet_content + "\",\"media\":{\"media_entities\":[],\"possibly_sensitive\":false},\"withDownvotePerspective\":false,\"withReactionsMetadata\":false,\"withReactionsPerspective\":false,\"withSuperFollowsTweetFields\":true,\"withSuperFollowsUserFields\":true,\"semantic_annotation_ids\":[],\"dark_request\":false,\"__fs_dont_mention_me_view_api_enabled\":false,\"__fs_interactive_text_enabled\":false,\"__fs_responsive_web_uc_gql_enabled\":false}"}
            tweet_headers = {'content-type': 'application/json'}

            resp = self.client.post(tweet_url,
                                    headers=tweet_headers,
                                    json=tweet_json,
                                    verify=False)

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif "rest_id" in resp.text:
                logger.success(f"{self.account_index} | {self.username} | Tweeted.")
            elif "Status is a duplicate" in resp.text:
                logger.warning(f"{self.account_index} | Tweet is a duplicate.")
                return True
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | {self.username} | Failed to tweet: {err}")
            return False

    def tweet_with_picture(self, tweet_content: str, tweet_image) -> bool:
        try:
            payload = {'media_data': tweet_image}

            resp = self.client.post("https://upload.twitter.com/1.1/media/upload.json", data=payload)
            media_id = resp.json()['media_id_string']

            resp = self.client.post(f'https://twitter.com/i/api/graphql/{self.query_ids["tweet"]}/CreateTweet',
                                    headers={"content-type": "application/json"},
                                    json={
                                        "variables": "{\"tweet_text\":\"" + tweet_content + "\",\"media\":{\"media_entities\":[{\"media_id\":\"" + media_id + "\",\"tagged_users\":[]}],\"possibly_sensitive\":false},\"withDownvotePerspective\":false,\"withReactionsMetadata\":false,\"withReactionsPerspective\":false,\"withSuperFollowsTweetFields\":true,\"withSuperFollowsUserFields\":true,\"semantic_annotation_ids\":[],\"dark_request\":false,\"__fs_dont_mention_me_view_api_enabled\":false,\"__fs_interactive_text_enabled\":false,\"__fs_responsive_web_uc_gql_enabled\":false}"})

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif "rest_id" in resp.text:
                logger.success(f"{self.account_index} | {self.username} | Tweeted with picture.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to tweet with picture: {err}")
            return False

    def comment(self, comment_content: str, tweet_id: str) -> bool:
        try:
            resp = self.client.post('https://twitter.com/i/api/graphql/' + self.query_ids["comment"] + '/CreateTweet',
                                    headers={"content-type": "application/json"},
                                    json={
                                        "variables": "{\"tweet_text\":\"" + comment_content + "\",\"reply\":{\"in_reply_to_tweet_id\":\"" + tweet_id + "\",\"exclude_reply_user_ids\":[]},\"media\":{\"media_entities\":[],\"possibly_sensitive\":false},\"withDownvotePerspective\":false,\"withReactionsMetadata\":false,\"withReactionsPerspective\":false,\"withSuperFollowsTweetFields\":true,\"withSuperFollowsUserFields\":true,\"semantic_annotation_ids\":[],\"dark_request\":false,\"withUserResults\":true,\"withBirdwatchPivots\":false}",
                                        "queryId": "" + self.query_ids['comment'] + ""})

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif "rest_id" in resp.text:
                logger.success(f"{self.account_index} | {self.username} | Commented.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to comment: {err}")
            return False

    def change_description(self, new_description: str) -> bool:
        try:
            resp = self.client.post("https://twitter.com/i/api/1.1/account/update_profile.json",
                                    headers={"content-type": 'application/x-www-form-urlencoded'},
                                    data={"description": new_description})

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif resp.status_code == 200:
                logger.success(f"{self.account_index} | {self.username} | Changed description.")
            elif "Rate limit exceeded" in resp.text:
                raise Exception("Rate limit exceeded.")
            elif "account is temporarily locked" in resp.text:
                raise Exception("Account is temporarily locked.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to change description: {err}")
            return False

    def change_username(self, new_username: str) -> bool:
        try:
            resp = self.client.post("https://twitter.com/i/api/1.1/account/settings.json",
                                    data={"screen_name": new_username})

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif resp.status_code == 200:
                logger.success(f"{self.account_index} | {self.username} | Changed username.")
            elif "Rate limit exceeded" in resp.text:
                raise Exception("Rate limit exceeded.")
            elif "account is temporarily locked" in resp.text:
                raise Exception("Account is temporarily locked.")
            elif "already been taken for Screen name" in resp.text:
                raise Exception(f"Uusername {new_username} is already taken.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to change username: {err}")
            return False

    def change_name(self, new_name: str) -> bool:
        try:
            resp = self.client.post("https://twitter.com/i/api/1.1/account/update_profile.json",
                                    headers={'content-type': 'application/x-www-form-urlencoded'},
                                    data={"name": new_name})

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif resp.status_code == 200:
                logger.success(f"{self.account_index} | {self.username} | Changed name.")
            elif "Rate limit exceeded" in resp.text:
                raise Exception("Rate limit exceeded.")
            elif "account is temporarily locked" in resp.text:
                raise Exception("Account is temporarily locked.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to change name: {err}")
            return False

    def change_background(self, new_background) -> bool:
        try:
            resp = self.client.post("https://twitter.com/i/api/1.1/account/update_profile_banner.json",
                                    data={"banner": new_background})

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif resp.status_code in (200, 201):
                logger.success(f"{self.account_index} | {self.username} | Changed background.")
            elif "Rate limit exceeded" in resp.text:
                raise Exception("Rate limit exceeded.")
            elif "account is temporarily locked" in resp.text:
                raise Exception("Account is temporarily locked.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to change background: {err}")
            return False

    def change_password(self, current_password: str, new_password: str) -> bool:
        try:
            resp = self.client.post("https://twitter.com/i/api/i/account/change_password.json",
                                    data={
                                        "current_password": current_password,
                                        "password": new_password,
                                        "password_confirmation": new_password
                                    })

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif resp.status_code in (200, 201):
                logger.success(f"{self.account_index} | Successfully changed the password.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to change password: {err}")
            return False

    def change_birthdate(self, new_birthdate: str) -> bool:
        try:
            resp = self.client.post("https://twitter.com/i/api/1.1/account/update_profile.json",
                                    headers={'content-type': 'application/x-www-form-urlencoded'},
                                    data={
                                        "birthdate_day": new_birthdate.split(":")[0],
                                        "birthdate_month": new_birthdate.split(":")[1],
                                        "birthdate_year": new_birthdate.split(":")[2],
                                        "birthdate_visibility": "mutualfollow",
                                        "birthdate_year_visibility": "self"
                                    })

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif resp.status_code in (200, 201):
                logger.success(f"{self.account_index} | {self.username} | Changed birthdate.")
            elif "Rate limit exceeded" in resp.text:
                raise Exception("Rate limit exceeded.")
            elif "account is temporarily locked" in resp.text:
                raise Exception("Account is temporarily locked.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to change birthdate: {err}")
            return False

    def change_location(self, new_location: str) -> bool:
        try:
            resp = self.client.post("https://twitter.com/i/api/1.1/account/update_profile.json",
                                    headers={'content-type': 'application/x-www-form-urlencoded'},
                                    data={"location": new_location})

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif resp.status_code in (200, 201):
                logger.success(f"{self.account_index} | {self.username} | Changed location.")
            elif "Rate limit exceeded" in resp.text:
                raise Exception("Rate limit exceeded.")
            elif "account is temporarily locked" in resp.text:
                raise Exception("Account is temporarily locked.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to change location: {err}")
            return False

    def change_profile_picture(self, new_picture) -> bool:
        try:
            resp = self.client.post("https://twitter.com/i/api/1.1/account/update_profile_image.json",
                                    data={"image": new_picture})

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif resp.status_code in (200, 201):
                logger.success(f"{self.account_index} | {self.username} | Changed profile picture.")
            elif "Rate limit exceeded" in resp.text:
                raise Exception("Rate limit exceeded.")
            elif "account is temporarily locked" in resp.text:
                raise Exception("Account is temporarily locked.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to change profile picture: {err}")
            return False

    def quote_tweet(self, tweet_content: str, tweet_id: str) -> bool:
        try:
            tweet_url = f'https://twitter.com/i/api/graphql/{self.query_ids["tweet"]}/CreateTweet'
            tweet_json = {
                "variables": "{\"tweet_text\":\"" + tweet_content + "\",\"attachment_url\":\"" + tweet_id + "\",\"media\":{\"media_entities\":[],\"possibly_sensitive\":false},\"withDownvotePerspective\":false,\"withReactionsMetadata\":false,\"withReactionsPerspective\":false,\"withSuperFollowsTweetFields\":true,\"withSuperFollowsUserFields\":true,\"semantic_annotation_ids\":[],\"dark_request\":false,\"__fs_dont_mention_me_view_api_enabled\":false,\"__fs_interactive_text_enabled\":false,\"__fs_responsive_web_uc_gql_enabled\":false}"}
            tweet_headers = {'content-type': 'application/json'}

            resp = self.client.post(tweet_url,
                                    headers=tweet_headers,
                                    json=tweet_json,
                                    verify=False)

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif "rest_id" in resp.text:
                logger.success(f"{self.account_index} | {self.username} | Tweeted.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to tweet with a quote: {err}")
            return False

    def send_direct_message(self, message_content: str, recipient_id: str) -> bool:
        try:
            resp = self.client.post("https://api.twitter.com/1.1/direct_messages/events/new.json",
                                    headers={"content-type": "application/json"},
                                    json={"event": {
                                        "type": "message_create", "message_create": {
                                            "target": {
                                                "recipient_id": recipient_id
                                            }, "message_data": {
                                                "text": message_content}
                                        }
                                    }})

            if resp.status_code == 200:
                logger.success(f"{self.account_index} | {self.username} | Sent the message.")
            elif "Rate limit exceeded" in resp.text:
                raise Exception("Rate limit exceeded.")
            elif "account is temporarily locked" in resp.text:
                raise Exception("Account is temporarily locked.")
            elif "You cannot send messages to this user." in resp.text:
                raise Exception("You cannot send messages to this user.")
            elif "You are sending a Direct Message to users that do not follow you." in resp.text:
                raise Exception("You are sending a Direct Message to users that do not follow you.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Failed to send a direct message: {err}")
            return False

    # tweet_image = open(f"/tweet_images/{tweet_image}", "rb")
    def comment_with_picture(self, comment_content: str, comment_image, tweet_id: str) -> bool:
        try:
            payload = {'media_data': comment_image}

            resp = self.client.post("https://upload.twitter.com/1.1/media/upload.json", data=payload)
            media_id = resp.json()['media_id_string']

            resp = self.client.post('https://twitter.com/i/api/graphql/5V_dkq1jfalfiFOEZ4g47A/CreateTweet',
                                    headers={'content-type': 'application/json'},
                                    json={
                                        'variables': {
                                            'tweet_text': comment_content,
                                            'reply': {
                                                'in_reply_to_tweet_id': tweet_id,
                                                'exclude_reply_user_ids': [],
                                            },
                                            'dark_request': False,
                                            'media': {
                                                'media_entities': [
                                                    {
                                                        'media_id': media_id,
                                                        'tagged_users': [],
                                                    },
                                                ],
                                                'possibly_sensitive': False,
                                            },
                                            'semantic_annotation_ids': [],
                                        },
                                        'features': {
                                            'c9s_tweet_anatomy_moderator_badge_enabled': True,
                                            'tweetypie_unmention_optimization_enabled': True,
                                            'responsive_web_edit_tweet_api_enabled': True,
                                            'graphql_is_translatable_rweb_tweet_is_translatable_enabled': True,
                                            'view_counts_everywhere_api_enabled': True,
                                            'longform_notetweets_consumption_enabled': True,
                                            'responsive_web_twitter_article_tweet_consumption_enabled': False,
                                            'tweet_awards_web_tipping_enabled': False,
                                            'responsive_web_home_pinned_timelines_enabled': True,
                                            'longform_notetweets_rich_text_read_enabled': True,
                                            'longform_notetweets_inline_media_enabled': True,
                                            'responsive_web_graphql_exclude_directive_enabled': True,
                                            'verified_phone_label_enabled': False,
                                            'freedom_of_speech_not_reach_fetch_enabled': True,
                                            'standardized_nudges_misinfo': True,
                                            'tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled': True,
                                            'responsive_web_media_download_video_enabled': False,
                                            'responsive_web_graphql_skip_user_profile_image_extensions_enabled': False,
                                            'responsive_web_graphql_timeline_navigation_enabled': True,
                                            'responsive_web_enhance_cards_enabled': False,
                                        },
                                        'queryId': '5V_dkq1jfalfiFOEZ4g47A',
                                    }
                                    ,
                                    )

            if "Could not authenticate you" in resp.text:
                logger.error(f"{self.account_index} | Account locked or token does not work.")
            elif "rest_id" in resp.text:
                logger.success(f"{self.account_index} | {self.username} | Commented.")
            else:
                raise Exception(resp.text)

            return True
        except Exception as err:
            logger.error(f"{self.account_index} | Successfully commented with picture: {err}")
            return False

    # print("   Enter the period of days for which you want to see last messages: ", end='')
    def check_dm(self, days_amount: int) -> bool:
        try:
            dm_time = int(days_amount) * 24 * 60 * 60
            date_now = datetime.now().timestamp()
            user_dm_time = (date_now - dm_time)

            resp = self.client.get('https://twitter.com/i/api/1.1/dm/inbox_initial_state.json',
                                   params={
                                       'nsfw_filtering_enabled': 'false',
                                       'filter_low_quality': 'false',
                                       'include_quality': 'all',
                                       'include_profile_interstitial_type': '1',
                                       'include_blocking': '1',
                                       'include_blocked_by': '1',
                                       'include_followed_by': '1',
                                       'include_want_retweets': '1',
                                       'include_mute_edge': '1',
                                       'include_can_dm': '1',
                                       'include_can_media_tag': '1',
                                       'include_ext_has_nft_avatar': '1',
                                       'include_ext_is_blue_verified': '1',
                                       'include_ext_verified_type': '1',
                                       'include_ext_profile_image_shape': '1',
                                       'skip_status': '1',
                                       'dm_secret_conversations_enabled': 'false',
                                       'krs_registration_enabled': 'true',
                                       'cards_platform': 'Web-12',
                                       'include_cards': '1',
                                       'include_ext_alt_text': 'true',
                                       'include_ext_limited_action_results': 'true',
                                       'include_quote_count': 'true',
                                       'include_reply_count': '1',
                                       'tweet_mode': 'extended',
                                       'include_ext_views': 'true',
                                       'dm_users': 'true',
                                       'include_groups': 'true',
                                       'include_inbox_timelines': 'true',
                                       'include_ext_media_color': 'true',
                                       'supports_reactions': 'true',
                                       'include_ext_edit_control': 'true',
                                       'include_ext_business_affiliations_label': 'true',
                                       'ext': 'mediaColor,altText,mediaStats,highlightedLabel,hasNftAvatar,voiceInfo,birdwatchPivot,superFollowMetadata,unmentionInfo,editControl',
                                   })

            for item in resp.json()['inbox_initial_state']["entries"]:
                try:
                    if 'message' in item:
                        message = item['message']
                        if 'message_data' in message and 'text' in message['message_data']:
                            sms_date = int(message["time"]) / 1000
                            time_to_check = user_dm_time - sms_date

                            if time_to_check < 0:
                                r = self.client.get(f"https://api.twitter.com/1.1/users/show.json?user_id={message['message_data']['sender_id']}")

                                sender = r.json()['screen_name']
                                if sender != self.username:
                                    sms_text = message["message_data"]["text"]
                                    date = datetime.fromtimestamp(sms_date).strftime('%d-%m-%y   %H:%M:%S')

                                    logger.info(f"{self.account_index} | Me: {self.username} | From: {sender} | Date: {date} | Text: {sms_text}")

                except Exception as err:
                    logger.error(f"{self.account_index} | Failed to get direct message: {err}")
            return True

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to check dm: {err}")
            return False

    def check_suspended(self) -> tuple[bool, str]:
        try:
            resp = self.client.get(f"https://api.twitter.com/1.1/users/show.json?screen_name={self.username.replace('@', '')}")

            if "User has been suspended" in resp.text or resp.json()['suspended']:
                logger.warning(f"{self.account_index} | {self.username} | User has been suspended.")
                return True, "ban"
            else:
                logger.success(f"{self.account_index} | {self.username} | Shadow Ban not found")
                return True, ""

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to check if suspended: {err}")
            return False, ""

    def _unfreeze(self) -> bool:
        try:
            logger.info(f"{self.account_index} | Starting to unfreeze the account.")
            two_captcha = utilities.TwoCaptcha(self.account_index, self.config['2captcha_api_key'], self.client, self.proxy)
            captcha_instance = utilities.SolveTwitterCaptcha(self.auth_token, self.ct0, two_captcha, self.account_index, self.user_agent)
            return captcha_instance.solve_captcha(self.proxy)

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to unfreeze: {err}")
            return False

    def start_unfreezer(self):
        try:
            self.username, self.ct0 = utilities.get_account_username(self.account_index, self.client)

            if self.username == "locked":
                return self._unfreeze()
            else:
                logger.success(f"{self.account_index} | This account does not need the unfreeze")
                return True

        except Exception as err:
            logger.error(f"{self.account_index} | Failed to unfreeze: {err}")
            return False
