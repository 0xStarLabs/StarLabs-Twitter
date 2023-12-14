from pyfiglet import figlet_format
from termcolor import cprint
from colorama import init
import sys
import os


def show_menu(menu_items: list):
    os.system("")
    print()
    counter = 0
    for item in menu_items:
        counter += 1

        if counter == len(menu_items):
            print('' + '[' + '\033[34m' + f'{counter}' + '\033[0m' + ']' + f' {item}\n')
        else:
            print('' + '[' + '\033[34m' + f'{counter}' + '\033[0m' + ']' + f' {item}')


MENU_ITEMS = [
    "Follow",
    "Retweet",
    "Like",
    "Tweet",
    "Tweet with picture",
    "Quote Tweet",
    "Comment",
    "Comment with picture",
    "Change Description",
    "Change Username",
    "Change Name",
    "Change Background",
    "Change Password",
    "Change Birthdate",
    "Change Location",
    "Change profile picture",
    "Send direct message",
    # "Mass DM",
    "Check account messages",
    "Check if account is valid",
    "Unfreeze Accounts",
    "Mutual Subscription"
]


def show_logo():
    os.system("cls")
    init(strip=not sys.stdout.isatty())
    print("\n")
    logo = figlet_format("STAR LABS", font="banner3")
    cprint(logo, 'light_cyan')
    print("")


def show_dev_info():
    print("\033[36m" + "VERSION: " + "\033[34m" + "1.2" + "\033[34m")
    print("\033[36m"+"DEV: " + "\033[34m" + "https://t.me/StarLabsTech" + "\033[34m")
    print("\033[36m"+"GitHub: " + "\033[34m" + "https://github.com/0xStarLabs/StarLabs-Twitter" + "\033[34m")
    print("\033[36m" + "DONATION EVM ADDRESS: " + "\033[34m" + "0x620ea8b01607efdf3c74994391f86523acf6f9e1" + "\033[0m")
    print()

