package main

import (
	"fmt"
	"github.com/zenthangplus/goccm"
	"slices"
	"sync"
	"twitter/extra"
)

func main() {
	for {
		extra.ShowLogo()
		extra.ShowDevInfo()
		extra.ShowMenu()

		options()
	}
}

func options() {
	// read user data, config and txt files
	userChoice := extra.ReadUserChoice(extra.MainMenu, "Input your choice")
	tasks, ok := extra.GetTasksTargetInfo(userChoice)
	if ok == false {
		return
	}
	txtData, ok := extra.GetTxtDataByTasks(userChoice)
	if ok == false {
		return
	}

	queryIDs := extra.ReadQueryIDs()
	config := extra.ReadConfig()
	// read accounts and proxies. set all proxies to "" if it doesn't need
	fmt.Print("\n\033[H\033[2J")
	accounts := extra.ReadTxtFile("accounts", "data/accounts.txt")
	proxies := extra.ReadTxtFile("proxies", "data/proxies.txt")

	var ipChangeLinks []string
	var threads int
	if config.Proxy.MobileProxy == "yes" {
		ipChangeLinks = extra.ReadTxtFile("ip change links", "data/ip_change_links.txt")
		threads = len(ipChangeLinks)
	} else {
		threads, ok = extra.UserInputInteger("How many concurrent threads do you want")
		if ok == false {
			return
		}
	}

	if len(proxies) == 0 {
		if extra.NoProxies() == false {
			return
		} else {
			for i := range accounts {
				proxies[i] = ""
			}
		}
	} else if config.Proxy.MobileProxy == "no" && len(proxies) < len(accounts) {
		newProxies := make([]string, len(accounts))
		for i := range accounts {
			newProxies[i] = proxies[i%len(proxies)]
		}
		proxies = newProxies
	}

	// create channels for account logs: locked, failed and usernames for mutual subscription
	failedAccounts := make(chan string, len(accounts))
	lockedAccounts := make(chan string, len(accounts))
	usernames := make(chan string, len(accounts))
	indexes := make(chan int, len(accounts))
	// save all accounts usernames to map for mutual subscription function
	// every account has 0 subs by default
	var mutualSubsUsernames sync.Map
	for _, key := range txtData.MyUsernames {
		mutualSubsUsernames.Store(key, 0)
	}
	for i := range accounts {
		indexes <- i
	}
	close(indexes)
	// threads
	goroutines := goccm.New(threads)
	fmt.Println()
	// set up account range. use all accounts if accounts range set to 0-0 in config by default
	var start, end int
	if config.Info.AccountRange.Start == 0 && config.Info.AccountRange.End == 0 || config.Info.AccountRange.Start >= config.Info.AccountRange.End {
		start = 0
		end = len(accounts)
	} else {
		start = config.Info.AccountRange.Start - 1
		end = config.Info.AccountRange.End
	}
	// start the flow
	if config.Proxy.MobileProxy == "yes" {
		for i, proxy := range proxies {
			data := MobileProxyData{
				Proxy:               proxy,
				IPChangeLink:        ipChangeLinks[i],
				Indexes:             indexes,
				Accounts:            accounts,
				Config:              config,
				UserChoice:          userChoice,
				Tasks:               tasks,
				TxtData:             txtData,
				QueryIDs:            queryIDs,
				FailedAccounts:      failedAccounts,
				LockedAccounts:      lockedAccounts,
				Usernames:           usernames,
				MutualSubsUsernames: &mutualSubsUsernames,
			}
			goroutines.Wait()
			go func(i int) {
				MobileProxyWrapper(data)
				goroutines.Done()
			}(i)
		}
	} else {
		for i := start; i < end && i < len(accounts); i++ {
			goroutines.Wait()
			go func(i int) {
				Process(i, accounts[i], proxies[i], config, userChoice, tasks, txtData, queryIDs, failedAccounts, lockedAccounts, usernames, &mutualSubsUsernames)
				goroutines.Done()
				extra.RandomSleep(config.Info.PauseBetweenAccounts.Start, config.Info.PauseBetweenAccounts.End)
			}(i)
		}
	}
	goroutines.WaitAllDone()
	close(failedAccounts)
	close(lockedAccounts)
	close(usernames)
	// write down logs about accounts to txt files
	var failedCounter int
	var lockedCounter int
	for account := range failedAccounts {
		failedCounter++
		extra.AppendToFile("data/failed_accounts.txt", account)
	}
	for account := range lockedAccounts {
		lockedCounter++
		extra.AppendToFile("data/locked_accounts.txt", account)
	}
	if slices.Contains(userChoice, "Check if account is valid") {
		extra.RewriteChannelToFile("data/my_usernames.txt", usernames)
	}
	// show statistic
	fmt.Println()
	extra.Logger{}.Info("STATISTICS: %d SUCCESS | %d FAILED | %d LOCKED", end-start-failedCounter-lockedCounter, failedCounter, lockedCounter)
	// exit with Enter
	fmt.Println("Press Enter to exit...")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}
}
