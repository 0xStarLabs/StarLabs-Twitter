package model

import (
	"fmt"
	"starlabs-twitter/src/utils"

	"github.com/zenthangplus/goccm"
)

func Start(accounts []utils.Account, userSelectedOptions []string, config utils.Config, taskData utils.TasksTargetInfo, fileData utils.TxtTasksData) {
	// Initialize goroutine manager
	goroutines := goccm.New(config.Settings.Threads)

	// Initialize statistics
	stats := NewStatistics(len(accounts))

	utils.Logger{Prefix: "Start"}.Info("Processing %d accounts with %d threads", len(accounts), config.Settings.Threads)

	// Process accounts
	for i, account := range accounts {
		goroutines.Wait()
		go func(acc utils.Account, index int) {
			defer goroutines.Done()

			// Process the account
			ProcessAccount(index+1, acc, userSelectedOptions, config, taskData, fileData, stats)

			// Random sleep between accounts as configured
			utils.RandomSleep(
				config.Settings.PauseBetweenAccounts.Start,
				config.Settings.PauseBetweenAccounts.End,
			)

		}(account, i)
	}

	// Wait for all goroutines to complete
	goroutines.WaitAllDone()

	// Print statistics table
	total, processed, success, failed, locked, suspended := stats.GetStats()
	utils.PrintStatistics(total, processed, success, failed, locked, suspended)

	// Wait for user input before exiting
	fmt.Print("Press Enter to exit...")
	fmt.Scanln()
}
