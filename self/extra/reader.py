from loguru import logger
import base64
import os


def read_txt_file(file_name: str, file_path: str) -> list:
    with open(file_path, "r") as file:
        items = [line.strip() for line in file]

    logger.success(f"Successfully loaded {len(items)} {file_name}.")
    return items


def get_txt_data(tasks: dict) -> tuple[dict, bool]:
    try:
        change_data = {}

        if any(tasks.get(task) not in [False, {}] for task in ['tweet', 'tweet with picture', 'quote tweet']):
            change_data['tweet'] = read_txt_file("tweets", "data/tweets.txt")

        if any(tasks.get(task) not in [False, {}] for task in ['comment', 'comment with picture']):
            change_data['comment'] = read_txt_file("comments", "data/comments.txt")

        if any(tasks.get(task) not in [False, {}] for task in ['tweet with picture', 'comment with picture', 'change background', 'change profile picture']):
            encoded_images = []

            for filename in os.listdir("data/pictures"):
                if filename.endswith((".png", ".jpg", ".jpeg")):
                    file_path = os.path.join("data/pictures", filename)

                    with open(file_path, 'rb') as image_file:
                        encoded_image = base64.b64encode(image_file.read()).decode('utf-8')
                        encoded_images.append(encoded_image)

            change_data['pictures'] = encoded_images

        if tasks['change password']:
            change_data['current password'] = read_txt_file("current passwords", "data/self/current_passwords.txt")
            change_data['new password'] = read_txt_file("new passwords", "data/self/new_passwords.txt")

        if tasks['mutual subscription']['start']:
            change_data['mutual subscription'] = read_txt_file("my usernames", "data/my_usernames.txt")

        change_data['description'] = read_txt_file("description", "data/self/description.txt") if tasks.get('change description') else None
        change_data['username'] = read_txt_file("usernames", "data/self/usernames.txt") if tasks.get('change username') else None
        change_data['name'] = read_txt_file("names", "data/self/names.txt") if tasks.get('change name') else None
        change_data['birthdate'] = read_txt_file("birthdate", "data/self/birthdate.txt") if tasks.get('change birthdate') else None
        change_data['location'] = read_txt_file("location", "data/self/location.txt") if tasks.get('change location') else None

        change_data['grabber'] = read_txt_file("accounts", "data/twitters_for_grabber.txt") if tasks.get('grabber') else None

        return change_data, True

    except Exception as err:
        logger.error(f"Failed to get change data: {err}")
        return {}, False
