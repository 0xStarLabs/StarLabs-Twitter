package model

import (
	"fmt"
	"starlabs-twitter/src/utils"
	"time"

	"github.com/0xStarLabs/TwitterAPI/client"
	"github.com/0xStarLabs/TwitterAPI/models"
	twitterUtils "github.com/0xStarLabs/TwitterAPI/utils"
)

// TwitterResponse represents the response from Twitter API calls
type TwitterResponse struct {
	Success bool
	Error   error
	Status  models.ActionStatus
}

// Status constants for Twitter API responses
const (
	StatusAuthError = models.StatusAuthError
	StatusLocked    = models.StatusLocked
)

func ProcessAccount(accountIndex int, account utils.Account, userSelectedOptions []string, config utils.Config, taskData utils.TasksTargetInfo, fileData utils.TxtTasksData, stats *Statistics) {
	// Initialize logger for this account
	logger := utils.Logger{Prefix: fmt.Sprintf("Account %d", accountIndex)}
	logger.Info("Starting to process account")

	// Initialize Twitter client
	twitter, err := initTwitterClient(account, config)
	if err != nil {
		logger.Error("Failed to initialize Twitter client: %v", err)
		utils.UpdateAccountStatus(utils.AccountsFilePath, account, "UNKNOWN")
		stats.IncrementFailed()
		return
	}

	// Validate account
	info, resp := twitter.IsValid()
	if info.Username != "" {
		utils.UpdateAccountInExcel(utils.AccountsFilePath, account, map[string]string{
			"USERNAME": info.Username,
		})
	}
	if !resp.Success {
		reason := handleTwitterError(logger, TwitterResponse{
			Success: resp.Success,
			Error:   resp.Error,
			Status:  resp.Status,
		})

		if info.Suspended {
			utils.UpdateAccountStatus(utils.AccountsFilePath, account, "SUSPENDED")
			stats.IncrementSuspended()
		} else {
			utils.UpdateAccountStatus(utils.AccountsFilePath, account, reason)
			switch reason {
			case "LOCKED":
				stats.IncrementLocked()
			case "AUTH_ERROR", "UNKNOWN":
				stats.IncrementFailed()
			}
		}
		return
	}
	logger.Success("Account validated successfully: @%s", info.Username)
	utils.UpdateAccountStatus(utils.AccountsFilePath, account, "VALID")

	// Process each selected option
	allActionsSuccessful := true
	for _, option := range userSelectedOptions {
		actionSuccessful := false

		// Special case for Check Valid
		if option == "Check Valid" {
			continue // Skip to next option in main loop
		}

		for attempt := 1; attempt <= config.Settings.Retries; attempt++ {
			var success bool
			var err error

			switch option {
			case "Follow":
				for _, username := range taskData.FollowTargets {
					logger.Info("Attempting to follow @%s", username)
					resp := twitter.Follow(username)
					success = resp.Success
					err = resp.Error
				}

			case "Unfollow":
				for _, username := range taskData.UnfollowTargets {
					logger.Info("Attempting to unfollow @%s", username)
					resp := twitter.Unfollow(username)
					success = resp.Success
					err = resp.Error
				}
			case "Retweet":
				for _, link := range taskData.RetweetLinks {
					logger.Info("Attempting to retweet %s", link)
					resp := twitter.Retweet(link)
					success = resp.Success
					err = resp.Error
				}
			case "Like":
				for _, link := range taskData.LikeLinks {
					logger.Info("Attempting to like %s", link)
					resp := twitter.Like(link)
					success = resp.Success
					err = resp.Error
				}
			case "Quote Tweet":
				for _, link := range taskData.QuoteTweetLinks {
					logger.Info("Attempting to quote tweet %s", link)
					resp := twitter.Tweet(fileData.Tweets[accountIndex-1], &client.TweetOptions{
						QuoteTweetURL: link,
					})
					success = resp.Success
					err = resp.Error
				}
			case "Quote Tweet with Picture":
				for _, link := range taskData.QuoteTweetLinks {
					logger.Info("Attempting to quote tweet with picture %s", link)
					resp := twitter.Tweet(fileData.Tweets[accountIndex-1], &client.TweetOptions{
						MediaBase64: fileData.Pictures[accountIndex-1],
						QuoteTweetURL: link,
					})
					success = resp.Success
					err = resp.Error
				}
			case "Comment":
				for _, link := range taskData.CommentTweetLinks {
					logger.Info("Attempting to comment on %s", link)
					resp := twitter.Comment(fileData.Comments[accountIndex-1], link, nil)
					success = resp.Success
					err = resp.Error
				}
			case "Comment with Picture":
				for _, link := range taskData.CommentTweetLinks {
					logger.Info("Attempting to comment on %s", link)
					resp := twitter.Comment(fileData.Comments[accountIndex-1], link, &client.CommentOptions{
						MediaBase64: fileData.Pictures[accountIndex-1],
					})
					success = resp.Success
					err = resp.Error
				}
			case "Tweet":
				logger.Info("Attempting to tweet %s", fileData.Tweets[accountIndex-1])
				resp := twitter.Tweet(fileData.Tweets[accountIndex-1], nil)
				success = resp.Success
				err = resp.Error
			case "Tweet with Picture":
				logger.Info("Attempting to tweet with picture")
				resp := twitter.Tweet(fileData.Tweets[accountIndex-1], &client.TweetOptions{
					MediaBase64: fileData.Pictures[accountIndex-1],
				})
				success = resp.Success
				err = resp.Error
			case "Vote on Poll":
				logger.Info("Attempting to vote on poll %s", taskData.Votes.Links[accountIndex-1])
				resp := twitter.VotePoll(taskData.Votes.Links[accountIndex-1], taskData.Votes.Choice)
				success = resp.Success
				err = resp.Error
			}
			// Handle response
			if success {
				logger.Success("Action '%s' completed successfully", option)
				actionSuccessful = true
				break
			} else {
				logger.Error("Action failed: %v", err)
				if attempt < config.Settings.Retries {
					sleepDuration := time.Duration(config.Settings.PauseBetweenRetries.Start) * time.Second
					time.Sleep(sleepDuration)
					logger.Warning("Retrying attempt %d/%d", attempt+1, config.Settings.Retries)
				}
			}
		}
		if !actionSuccessful {
			allActionsSuccessful = false
		}
	}

	// Update final statistics
	if allActionsSuccessful {
		stats.IncrementSuccess()
	} else {
		stats.IncrementFailed()
	}
	stats.IncrementProcessed()
}

func initTwitterClient(account utils.Account, config utils.Config) (*client.Twitter, error) {
	// Create new account instance
	twitterAccount := client.NewAccount(
		account.AuthToken,
		"", // CSRF token (optional)
		account.Proxy,
	)
	twitterConfig := models.NewConfig()
	twitterConfig.MaxRetries = config.Settings.Retries
	twitterConfig.LogLevel = twitterUtils.LogLevelDebug

	// Initialize Twitter client
	return client.NewTwitter(twitterAccount, twitterConfig)
}

func handleTwitterError(logger utils.Logger, resp TwitterResponse) string {
	switch resp.Status {
	case StatusAuthError:
		return "AUTH_ERROR"
	case StatusLocked:
		return "LOCKED"
	default:
		logger.Error("Unknown error: %v", resp.Error)
		return "UNKNOWN"
	}
}
