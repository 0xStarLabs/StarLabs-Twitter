# Twitter Bot 2.0.0 🚀

A powerful and flexible Twitter automation tool with multiple features and parallel processing capabilities.

## 🌟 Features

- ✨ Multi-threaded processing
- 🔄 Automatic retries
- 📊 Beautiful statistics display
- 🔐 Proxy support
- 📝 Excel-based account management
- 🎯 Multiple action support:
  - Follow/Unfollow users
  - Tweet with/without images
  - Quote tweets
  - Like posts
  - Comment on tweets
  - Vote in polls
  - Mutual subscription system

## 📋 Requirements

- Go 1.21 or higher
- Excel file with Twitter accounts
- Valid Twitter authentication tokens
- (Optional) Proxies

## 🚀 Installation

1. Clone the repository:
```bash
git clone https://github.com/0xStarLabs/StarLabs-Twitter
```
```bash
cd StarLabs-Twitter
```
2. Install dependencies:
```bash
go mod tidy
```
3. Run the bot:
```bash
go run main.go
```

## 📁 Project Structure
StarLabs-Twitter/
├── data/
│ ├── accounts.xlsx # Twitter accounts data
│ ├── config.yaml # Configuration settings
│ ├── comment_text.txt # Comments for interactions
│ ├── tweet_text.txt # Tweets content
│ └── images/ # Images for media tweets (jpg or png)
└── src/
├── model/ # Core business logic
└── utils/ # Utility functions

## 📝 Configuration

The configuration is done in the `data/config.yaml` file.

### 1. accounts.xlsx Structure
Required columns:
| AUTH_TOKEN | PROXY             | USERNAME | STATUS |
|------------|-------------------|----------|--------|
| token1     | user:pass@ip:port | user1    | VALID  |
| token2     | user:pass@ip:port | user2    | VALID  |

- `AUTH_TOKEN`: Your Twitter auth_token (required)
- `PROXY`: Proxy in format `user:pass@ip:port` (optional)
- `USERNAME`: Will be auto-filled by the bot
- `STATUS`: Account status, auto-updated by the bot

### 2. config.yaml Settings
```yaml
SETTINGS:
  THREADS: 5                      # Number of parallel threads
  RETRIES: 3                      # Retry attempts for failed actions
  PAUSE_BETWEEN_RETRIES: [1, 3]   # Random pause between retries (seconds)
  PAUSE_BETWEEN_ACCOUNTS: [1, 3]  # Random pause between accounts (seconds)
  ACCOUNTS_RANGE: [0, 0]          # Account range to process (0,0 = all)

MUTUAL_SUBSCRIPTION:
  FOLLOWERS_FOR_EVERY_ACCOUNT: [1, 2]  # Random followers per account
```

### 3. Content Files
- `tweet_text.txt`: One tweet text per line
- `comment_text.txt`: One comment text per line
- `images/`: Place your .jpg or .png images here

## 🎮 Usage

1. Prepare your files:
   - Fill `accounts.xlsx` with auth tokens and proxies
   - Configure `config.yaml` with desired settings
   - Add content to text files and images folder

2. Run the bot:
```bash
go run main.go
```

3. Select actions from the menu:
   - Use ↑/↓ to navigate
   - Space to select/unselect
   - Enter to confirm
   - Multiple actions can be selected

## 🔥 Available Actions

- **Check Valid**: Validates accounts and updates their status
- **Follow**: Follow target users
- **Unfollow**: Unfollow target users
- **Tweet**: Post text updates
- **Tweet with Picture**: Share images with text
- **Quote Tweet**: Quote and comment on tweets
- **Like**: Like target tweets
- **Comment**: Comment on tweets
- **Vote on Poll**: Participate in Twitter polls
- **Mutual Subscription**: Create follower network between accounts

## 📊 Statistics

The bot provides real-time statistics in a beautiful table:
```
┌────────────┬───────┐
│ STATISTICS │ VALUE │
├────────────┼───────┤
│ Total      │ 5     │
│ Processed  │ 5     │
│ Success    │ 3     │
│ Failed     │ 1     │
│ Locked     │ 1     │
│ Suspended  │ 0     │
└────────────┴───────┘
```

## 🤝 Support

- GitHub: https://github.com/0xStarLabs
- Telegram: https://t.me/StarLabsTech
- Chat: https://t.me/StarLabsChat

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is for educational purposes only. Use at your own risk and in accordance with Twitter's Terms of Service.
