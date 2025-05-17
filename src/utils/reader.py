import base64
import json
import os
from typing import List
from loguru import logger
import threading
import pandas as pd
from .constants import Account
from .proxy_parser import Proxy


def read_txt_file(file_name: str, file_path: str) -> list:
    with open(file_path, "r") as file:
        items = [line.strip() for line in file]

    logger.success(f"Successfully loaded {len(items)} {file_name}.")
    return items


def split_list(lst, chunk_size=90):
    return [lst[i : i + chunk_size] for i in range(0, len(lst), chunk_size)]


# Mutex for thread-safe Excel file access
excel_mutex = threading.Lock()


def read_accounts_from_excel(
    file_path: str, start_index: int = 1, end_index: int = 0
) -> list:
    """
    Reads account data from an Excel file and returns a list of Account objects.
    It uses the provided start and end indices to determine which accounts to process.

    Args:
        file_path: Path to the Excel file
        start_index: Starting index (1-based, default is 1)
        end_index: Ending index (0 means all accounts, default is 0)

    Returns:
        List of Account objects

    Example:
        accounts = read_accounts_from_excel("./data/accounts.xlsx", 1, 10)
        for account in accounts:
            print(f"Account: {account}")
    """
    # Lock the mutex before accessing the file
    with excel_mutex:
        try:
            # Check if file exists
            if not os.path.exists(file_path):
                raise FileNotFoundError(f"Excel file not found: {file_path}")

            # Read Excel file
            df = pd.read_excel(file_path)

            # Check if required columns exist
            if "AUTH_TOKEN" not in df.columns:
                raise ValueError("Excel file must have an 'AUTH_TOKEN' column")

            accounts = []

            # Process rows until empty AUTH_TOKEN
            for index, row in df.iterrows():
                # Skip rows with empty auth_token
                if (
                    pd.isna(row.get("AUTH_TOKEN", ""))
                    or row.get("AUTH_TOKEN", "") == ""
                ):
                    break

                # Get auth token
                auth_token = row.get("AUTH_TOKEN", "")

                # Get proxy with default empty string
                proxy_str = row.get("PROXY", "")
                if pd.isna(proxy_str):
                    proxy_str = ""

                # Process non-empty proxy through proxy_parser
                if proxy_str:
                    try:
                        proxy_obj = Proxy.from_str(proxy_str)
                        # Convert to user:pass@ip:port format
                        proxy_str = proxy_obj.get_default_format()
                    except Exception as e:
                        # Add row number to error message
                        raise ValueError(
                            f"Failed to parse proxy '{proxy_str}' on row {index + 2}: {str(e)}"
                        )

                # Create account object
                account = Account(auth_token=auth_token, proxy=proxy_str)
                accounts.append(account)

            if not accounts:
                raise ValueError("No valid accounts found in the Excel file. Please check that the file has accounts with valid AUTH_TOKEN values.")

            # If both are 1 and 0 or equal, process all accounts
            if (start_index == 1 and end_index == 0) or start_index == end_index:
                filtered_accounts = accounts
            else:
                # Validate indices
                if start_index < 1:
                    start_index = 1
                if end_index == 0 or end_index > len(accounts):
                    end_index = len(accounts)

                # Return accounts within the range
                filtered_accounts = accounts[start_index - 1 : end_index]

            logger.success(
                f"Successfully loaded {len(filtered_accounts)} accounts from Excel file."
            )
            return filtered_accounts

        except pd.errors.EmptyDataError:
            logger.error("The Excel file is empty")
            return []
        except Exception as e:
            logger.error(f"Error reading Excel file: {e}")
            return []


async def read_pictures(file_path: str) -> List[str]:
    """
    Считывает изображения из указанной папки и кодирует их в base64

    Args:
        file_path: Путь к папке с изображениями

    Returns:
        Список закодированных изображений в формате base64
    """
    encoded_images = []

    # Создаем папку, если она не существует
    os.makedirs(file_path, exist_ok=True)
    logger.info(f"Reading pictures from {file_path}")

    try:
        # Получаем список файлов
        files = os.listdir(file_path)

        if not files:
            logger.warning(f"No files found in {file_path}")
            return encoded_images

        # Обрабатываем каждый файл
        for filename in files:
            if filename.endswith((".png", ".jpg", ".jpeg")):
                # Формируем полный путь к файлу
                image_path = os.path.join(file_path, filename)

                try:
                    with open(image_path, "rb") as image_file:
                        encoded_image = base64.b64encode(image_file.read()).decode(
                            "utf-8"
                        )
                        encoded_images.append(encoded_image)
                except Exception as e:
                    logger.error(f"Error loading image {filename}: {str(e)}")

    except FileNotFoundError:
        logger.error(f"Directory not found: {file_path}")
    except PermissionError:
        logger.error(f"Permission denied when accessing: {file_path}")
    except Exception as e:
        logger.error(f"Error reading pictures from {file_path}: {str(e)}")

    return encoded_images
