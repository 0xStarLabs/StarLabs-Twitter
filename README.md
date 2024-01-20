#  StarLabs - Twitter 


![Logo](https://i.postimg.cc/fWtrP7hN/twitter.png)

## [SEE ENGLISH VERSION BELOW ](https://github.com/0xStarLabs/StarLabs-Twitter?tab=readme-ov-file#english-version)👇

## 🔗 Links
[![Telegram channel](https://img.shields.io/endpoint?url=https://runkit.io/damiankrawczyk/telegram-badge/branches/master?url=https://t.me/StarLabsTech)](https://t.me/StarLabsTech)
[![Telegram chat](https://img.shields.io/endpoint?url=https://runkit.io/damiankrawczyk/telegram-badge/branches/master?url=https://t.me/StarLabsChat)](https://t.me/StarLabsChat)

🔔 CHANNEL: https://t.me/StarLabsTech

💬 CHAT: https://t.me/StarLabsChat

💰 DONATION EVM ADDRESS: 0x620ea8b01607efdf3c74994391f86523acf6f9e1

📖 FULL TUTORIAL: https://teletype.in/@izanumi/star_labs_twitter


## 🤖 | Функционал:

🟢 Подписка

🟢 Лайк

🟢 Ретвит

🟢 Твиты (с картинкой, с цитатой другого твита)

🟢 Комментарии (с картинкой и без)

🟢 Смена любых данных аккаунта: имя, пароль и тд

🟢 Проверка аккаунта на валидность

🟢 Взаимная подписка между аккаунтами


## 🚀 Installation
```
git clone https://github.com/0xStarLabs/StarLabs-Twitter.git

cd StarLabs-Twitter

go build

# Перед началом работы настройте необходимые модули в файлах config.yaml и /data

twitter.exe
```

## ⚙️ Config

| Name | Description |
| --- | --- |
| max_tasks_retries | Максимальное количество попыток при выполнении задания |
| pause_between_tasks | Пауза между каждым действием |
| pause_between_accounts | Пауза между каждым аккаунтом |
| account_range | Диапазон аккаунтов для работы |
| 2captcha_api_key | Ключ от https://2captcha.com |
| mobile_proxy | Мобильные прокси |
| change_ip_pause | Пауза после смены айпи мобильных прокси |
| random | Брать данные из тхт файлов в случайном порядке |

## 🗂️ Data

Данные в папке data:

| Name | Description |
| --- | --- |
| accounts.txt | Содержит аккаунты твиттер |
| comments.txt | Содержит текст на функций комментариев |
| failed_accounts.txt | Содержит аккаунты которые не смогли выполнить задачу |
| locked_accounts.txt | Содержит аккаунты которые временно заблокированные |
| valid_accounts.txt | Содержит валидные аккаунты |
| tweets.txt | Содержит текст для функции твитов |
| proxies.txt | Содержит прокси в формате user:pass@ip:port |
| my_usernames.txt | Содержит юзернеймы аккаунтов |
| ip_change_links.txt | Содержит ссылки для смены айпи мобильных прокси |

## Дисклеймер
Автоматизация учетных записей пользователей Twitter, также известных как самостоятельные боты, является нарушением Условий обслуживания и правил сообщества Twitter и приведет к закрытию вашей учетной записи (аккаунтов). Рекомендуется осмотрительность. Я не буду нести ответственность за ваши действия. Прочтите об Условиях обслуживания Twitter и Правилах сообщества.

Это программное обеспечение было написано как доказательство концепции того, что учетные записи Twitter могут быть автоматизированы и могут выполнять действия, выходящие за рамки обычных пользователей Twitter, чтобы Twitter мог вносить изменения. Авторы  освобождаются от любой ответственности, которую может повлечь за собой ваше использование.

## ENGLISH VERSION:

## 🤖 | Features :

🟢 Subscribe

🟢 Like

🟢 Retweet

🟢 Tweets (with picture, with quote of another tweet)

🟢 Comments (with and without picture)

🟢 Change any account details: name, password, etc.

🟢 Checking account validity

🟢 Mutual subscription between accounts

## 🚀 Installation
```
git clone https://github.com/0xStarLabs/StarLabs-Twitter.git

cd StarLabs-Twitter

go build

# Before you start, configure the required modules in config.yaml and /data files

twitter.exe
```

## ⚙️ Config

| Name | Description |
| --- | --- |
| max_tasks_retries | Maximum number of attempts to complete a task |
| pause_between_tasks | pause between each action |
| pause_between_accounts | pause between each account |
| account_range | range of accounts to work |
| 2captcha_api_key | Key from https://2captcha.com |
| mobile_proxy | mobile proxies |
| change_ip_pause | pause after changing the mobile proxy's ip |
| random | take data from txt files in random order |


## 🗂️ Data

The data is in the data folder:

| Name | Description |
| --- | --- |
| accounts.txt | Contains Twitter accounts |
| comments.txt | Contains text on comment functions |
| failed_accounts.txt | Contains accounts that failed to complete the task |
| locked_accounts.txt | Contains accounts that are temporarily locked |
| valid_accounts.txt | Contains valid accounts |
| tweets.txt | Contains the text for the tweets function |
| proxies.txt | Contains proxies in the format user:pass@ip:port |
| my_usernames | Contains account usernames |
| ip_change_links | Contains links to change mobile proxy IPs |


## Disclaimer

The automation of User Twitter accounts also known as self-bots is a violation of Twitter Terms of Service & Community guidelines and will result in your account(s) being terminated. Discretion is adviced. I will not be responsible for your actions. Read about Twitter Terms Of service and Community Guidelines

Twitter StarLabs was written as a proof of concept that Twitter accounts can be automated and can perform actions beyond the scope of regular Twitter Users so that Twitter can make changes. The Twitter StarLabs authors are released of any liabilities which your usage may entail.

