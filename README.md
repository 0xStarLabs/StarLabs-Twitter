# StarLabs Twitter Bot 2.1 🌟
A powerful Python-based Twitter automation tool with multithreading support and comprehensive statistics tracking.

## 📚 Documentation & Tutorials
> ### [📖 English Tutorial](https://star-labs.gitbook.io/star-labs/twitter/eng)
> ### [📖 Russian Tutorial](https://star-labs.gitbook.io/star-labs/twitter/ru)

## ✨ Features
- 📊 Real-time statistics display
- 🎨 Beautiful CLI interface with gradient display
- 🔄 Automatic retries with configurable attempts
- 🔧 Configurable execution settings
- 📝 Excel-based account management
- 🚀 Multiple account support with optional shuffle
- 📱 Telegram integration for reporting
- 🛠️ Wide range of Twitter actions:
  - Follow/Unfollow users
  - Like tweets
  - Retweet posts
  - Post tweets with/without images
  - Comment on tweets with/without images
  - Quote tweets with/without images
  - Account validation

## 📋 Requirements
- Python 3.11.6 or higher
- Excel file with Twitter accounts
- Valid Twitter authentication tokens
- (Optional) Proxies for account management

## 🔧 Installation
1. Clone the repository:
```bash
git clone https://github.com/0xStarLabs/StarLabs-Twitter
cd StarLabs-Twitter
```

2. Install the requirements:
```bash
pip install -r requirements.txt
```

3. Configure your accounts and settings (see Configuration section)

4. Run the bot:
```bash
python main.py
```

## 📁 Project Structure
```
StarLabs-Twitter/
├── data/
│   ├── accounts.xlsx        # Twitter accounts data
│   ├── tweet_text.txt       # Tweets content
│   ├── comment_text.txt     # Comments for interactions
│   └── images/              # Images for media tweets
├── src/
│   ├── model/               # Core Twitter functionality
│   │   ├── twitter/         # Twitter API handlers
│   │   ├── instance.py      # Account instance management
│   │   ├── prepare_data.py  # Data preparation utilities
│   │   └── start.py         # Main execution flow
│   └── utils/               # Utility functions
│       ├── telegram_logger.py # Telegram integration
│       ├── reader.py        # File and data readers
│       ├── output.py        # CLI output formatting
│       └── config.py        # Configuration manager
├── process.py               # Main process handler
├── config.yaml              # Configuration settings
└── main.py                  # Entry point
```

## 📝 Configuration

### 1. Account Setup (accounts.xlsx)
Your Excel file should have the following columns:
```
AUTH_TOKEN | PROXY | USERNAME | STATUS
```
- **AUTH_TOKEN**: Your Twitter auth_token (required)
- **PROXY**: Proxy in format user:pass@ip:port (optional)
- **USERNAME**: Will be auto-filled by the bot
- **STATUS**: Account status, auto-updated by the bot

### 2. Configuration (config.yaml)
The bot's behavior is controlled through the config.yaml file:

```yaml
SETTINGS:
  THREADS: 1                      # Number of parallel threads
  ATTEMPTS: 5                     # Retry attempts for failed actions
  ACCOUNTS_RANGE: [0, 0]          # Account range to process (0,0 = all)
  EXACT_ACCOUNTS_TO_USE: []       # Specific accounts to use (e.g., [1, 4, 6])
  SHUFFLE_ACCOUNTS: true          # Randomize account processing order
  PAUSE_BETWEEN_ATTEMPTS: [3, 10] # Random pause between retries (seconds)
  RANDOM_PAUSE_BETWEEN_ACCOUNTS: [3, 10]  # Pause between accounts (seconds)
  RANDOM_PAUSE_BETWEEN_ACTIONS: [3, 10]   # Pause between actions (seconds)
  RANDOM_INITIALIZATION_PAUSE: [3, 10]    # Initial pause for accounts
  
  # Telegram notifications
  SEND_TELEGRAM_LOGS: false
  SEND_ONLY_SUMMARY: false
  TELEGRAM_BOT_TOKEN: "your_token_here"
  TELEGRAM_USERS_IDS: [your_user_id]

FLOW:
  SKIP_FAILED_TASKS: false        # Continue after task failures

TWEETS:
  RANDOM_TEXT_FOR_TWEETS: false   # Use random text from file
  RANDOM_PICTURE_FOR_TWEETS: true # Use random pictures from folder

COMMENTS:
  RANDOM_TEXT_FOR_COMMENTS: false # Use random text from file
  RANDOM_PICTURE_FOR_COMMENTS: true # Use random pictures for comments
```

### 3. Content Files
- **tweet_text.txt**: One tweet text per line
- **comment_text.txt**: One comment text per line
- **images/**: Place your .jpg or .png images for media tweets/comments

## 🚀 Usage
1. Prepare your account data in accounts.xlsx
2. Configure bot settings in config.yaml
3. Run the bot:
```bash
python main.py
```
4. Select an option from the menu:
   - ⭐️ Start farming
   - 🔧 Edit config
   - 👋 Exit

5. Choose tasks to perform:
   - Follow
   - Like
   - Retweet
   - Comment
   - Comment with image
   - Tweet
   - Tweet with image
   - Quote
   - Quote with image
   - Unfollow
   - Check Valid
   - Exit

6. For each task, the bot will prompt for necessary input such as usernames or tweet URLs

## 📊 Statistics
The bot tracks detailed statistics for each run:
- Total accounts processed
- Success/failure rates by task
- Individual account results
- Task-specific performance metrics

Optional Telegram reporting can send detailed statistics at the end of execution.

## 🌐 Support
- GitHub: https://github.com/0xStarLabs
- Telegram: https://t.me/StarLabsTech
- Chat: https://t.me/StarLabsChat

## ⚠️ Disclaimer
This tool is for educational purposes only. Use at your own risk and in accordance with Twitter's Terms of Service.