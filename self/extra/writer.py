from loguru import logger
import threading


def update_new_token(old_token: str, new_token: str, lock: threading.Lock,) -> bool:
    try:
        with lock:
            with open("data/discord_tokens.txt", 'r') as file:
                lines = file.readlines()

            lines = [new_token + '\n' if line.strip() == old_token else line for line in lines]

            with open("data/discord_tokens.txt", 'w') as file:
                file.writelines(lines)

        return True

    except Exception as err:
        logger.error(f"Failed to update new token in discord_tokens.txt: {err}")
        return False
