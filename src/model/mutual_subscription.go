package model

import (
	"fmt"
	"math/rand/v2"
	"starlabs-twitter/src/utils"
	"sync"
	"time"

	"github.com/0xStarLabs/TwitterAPI/client"
	"github.com/zenthangplus/goccm"
)

type AccountInstance struct {
	Account utils.Account
	Twitter *client.Twitter
}

func MutualSubscription(accounts []utils.Account, config utils.Config) {
	// Initialize goroutine manager
	goroutines := goccm.New(config.Settings.Threads)

	// Initialize statistics
	stats := NewStatistics(len(accounts))
	utils.Logger{Prefix: "Start"}.Info("Processing %d accounts with %d threads", len(accounts), config.Settings.Threads)

	// Filter accounts first
	filteredAccounts := filterAccounts(accounts, config)
	if len(filteredAccounts) == 0 {
		utils.Logger{Prefix: "MutualSubscription"}.Error("No valid accounts found")
		return
	}

	// Create subscription plan for each account
	subscriptionPlan := make(map[string][]string)
	for _, acc := range filteredAccounts {
		// Generate random number of followers for this account
		followersCount := utils.RandomInt(
			config.MutualSubscription.FollowersForEveryAccount.Start,
			config.MutualSubscription.FollowersForEveryAccount.End,
		)

		// Get random accounts to follow this account (excluding self)
		availableFollowers := make([]AccountInstance, 0)
		for _, potential := range filteredAccounts {
			if potential.Account.Username != acc.Account.Username {
				availableFollowers = append(availableFollowers, potential)
			}
		}

		// Shuffle available followers
		ShuffleAccounts(availableFollowers)

		// Take only required number of followers
		if followersCount > len(availableFollowers) {
			followersCount = len(availableFollowers)
		}

		followers := availableFollowers[:followersCount]

		// Add to subscription plan
		for _, follower := range followers {
			subscriptionPlan[follower.Account.Username] = append(
				subscriptionPlan[follower.Account.Username],
				acc.Account.Username,
			)
		}
	}

	// Execute the subscription plan
	for _, acc := range filteredAccounts {
		goroutines.Wait()
		go func(account AccountInstance) {
			defer goroutines.Done()
			logger := utils.Logger{Prefix: fmt.Sprintf("Account @%s", account.Account.Username)}

			// Get list of accounts to follow
			toFollow := subscriptionPlan[account.Account.Username]
			if len(toFollow) == 0 {
				logger.Info("No accounts to follow")
				stats.IncrementProcessed()
				stats.IncrementSuccess() // Count as success if there's nothing to do
				return
			}

			logger.Info("Will follow %d accounts", len(toFollow))

			allSuccess := true // Track overall success across all follows
			// Follow each account
			for _, username := range toFollow {
				followSuccess := false // Track success for this specific follow
				var lastError error

				// Try following with retries
				for attempt := 1; attempt <= config.Settings.Retries; attempt++ {
					logger.Info("Following @%s (attempt %d/%d)", username, attempt, config.Settings.Retries)
					resp := account.Twitter.Follow(username)

					if resp.Success {
						logger.Success("Successfully followed @%s", username)
						followSuccess = true
						break // Success, no need to retry
					} else {
						lastError = resp.Error
						logger.Warning("Failed to follow @%s (attempt %d/%d): %v",
							username, attempt, config.Settings.Retries, resp.Error)

						if attempt < config.Settings.Retries {
							sleepDuration := time.Duration(config.Settings.PauseBetweenRetries.Start) * time.Second
							time.Sleep(sleepDuration)
							logger.Info("Retrying...")
						}
					}
				}

				// Only log final error if all retries failed
				if !followSuccess {
					logger.Error("Failed to follow @%s after %d attempts: %v",
						username, config.Settings.Retries, lastError)
					allSuccess = false
				}

				// Sleep between follows
				utils.RandomSleep(
					config.Settings.PauseBetweenAccounts.Start,
					config.Settings.PauseBetweenAccounts.End,
				)
			}

			// Update statistics after processing all follows
			stats.IncrementProcessed()
			if allSuccess {
				stats.IncrementSuccess()
			} else {
				stats.IncrementFailed()
			}
		}(acc)
	}

	goroutines.WaitAllDone()

	// Print statistics
	total, processed, success, failed, locked, suspended := stats.GetStats()
	utils.PrintStatistics(total, processed, success, failed, locked, suspended)

	fmt.Print("Press Enter to exit...")
	fmt.Scanln()
}

func filterAccounts(accounts []utils.Account, config utils.Config) []AccountInstance {
	// Initialize goroutine manager
	goroutines := goccm.New(config.Settings.Threads)

	utils.Logger{Prefix: "Start"}.Info("Filtering %d accounts with %d threads", len(accounts), config.Settings.Threads)

	// Use a mutex to protect the slice
	var mutex sync.Mutex
	filteredAccounts := []AccountInstance{}

	// Process accounts
	for i, account := range accounts {
		goroutines.Wait()
		go func(acc utils.Account, index int) {
			defer goroutines.Done()

			for attempt := 1; attempt <= config.Settings.Retries; attempt++ {
				// Initialize Twitter client
				twitter, err := initTwitterClient(account, config)
				if err != nil {
					continue
				}

				// Validate account
				info, resp := twitter.IsValid()
				if resp.Success && info.Username != "" {
					// Create AccountInstance with the correct username
					instance := AccountInstance{
						Account: utils.Account{
							AuthToken: acc.AuthToken,
							Proxy:     acc.Proxy,
							Username:  info.Username, // Store the actual Twitter username
							Status:    acc.Status,
						},
						Twitter: twitter,
					}
					utils.UpdateAccountInExcel(account.AuthToken, instance.Account, map[string]string{
						"USERNAME": instance.Account.Username,
						"STATUS":   instance.Account.Status,
					})
					// Thread-safe append
					mutex.Lock()
					filteredAccounts = append(filteredAccounts, instance)
					mutex.Unlock()
					break
				}

				// Random sleep between accounts as configured
				utils.RandomSleep(
					config.Settings.PauseBetweenAccounts.Start,
					config.Settings.PauseBetweenAccounts.End,
				)
			}
		}(account, i)
	}

	goroutines.WaitAllDone()

	// Print statistics table
	utils.Logger{Prefix: "Filter"}.Info("Valid accounts after filtering: %d/%d", len(filteredAccounts), len(accounts))

	return filteredAccounts
}

// ShuffleAccounts randomly shuffles a slice of AccountInstance
func ShuffleAccounts(accounts []AccountInstance) {
	rand.Shuffle(len(accounts), func(i, j int) {
		accounts[i], accounts[j] = accounts[j], accounts[i]
	})
}
