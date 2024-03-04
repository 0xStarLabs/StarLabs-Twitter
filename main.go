package main

import (
	"fmt"
	"github.com/zenthangplus/goccm"
	"slices"
	"strings"
	"sync"
	"time"
	"twitter/extra"
	"twitter/grabber"
	"twitter/instance"
	"twitter/utils"
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
	creationDateChan := make(chan string, len(accounts))
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
				CreationDateChan:    creationDateChan,
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
				Process(i, accounts[i], proxies[i], config, userChoice, tasks, txtData, queryIDs, failedAccounts, lockedAccounts, usernames, &mutualSubsUsernames, creationDateChan)
				goroutines.Done()
				extra.RandomSleep(config.Info.PauseBetweenAccounts.Start, config.Info.PauseBetweenAccounts.End)
			}(i)
		}
	}
	goroutines.WaitAllDone()
	close(failedAccounts)
	close(lockedAccounts)
	close(usernames)
	close(creationDateChan)
	failedAccountsSet := make(map[string]struct{})
	lockedAccountsSet := make(map[string]struct{})

	var failedCounter int
	var lockedCounter int

	for account := range failedAccounts {
		failedCounter++
		failedAccountsSet[account] = struct{}{}
		extra.AppendToFile("data/failed_accounts.txt", account)
	}

	for account := range lockedAccounts {
		lockedCounter++
		lockedAccountsSet[account] = struct{}{}
		extra.AppendToFile("data/locked_accounts.txt", account)
	}

	for _, account := range accounts {
		if _, failed := failedAccountsSet[account]; !failed {
			if _, locked := lockedAccountsSet[account]; !locked {
				extra.AppendToFile("data/valid_accounts.txt", account)
			}
		}
	}

	if slices.Contains(userChoice, "Check if account is valid") {
		extra.RewriteChannelToFile("data/my_usernames.txt", usernames)
		extra.RewriteChannelToFile("data/creation_date.txt", creationDateChan)
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

func Process(index int, twitterAccount string, proxy string, config extra.Config, userChoice []string, tasks extra.TasksTargetInfo, txtData extra.TxtTasksData, queryIDs extra.QueryIDs, failedAccounts chan<- string, lockedAccounts chan<- string, usernames chan<- string, mutualSubsUsernames *sync.Map, creationDateChan chan<- string) {
	var report bool

	if len(twitterAccount) == 40 && !strings.Contains(twitterAccount, ":") {
		// auth_token
	} else if strings.Contains(twitterAccount, ":") && len(strings.Split(twitterAccount, ":")) >= 3 && len(strings.Split(twitterAccount, ":")[2]) == 40 {
		// login:pass:auth_token:json_cookies
		twitterAccount = strings.Split(twitterAccount, ":")[2]
	} else if strings.Contains(twitterAccount, ":") && len(strings.Split(twitterAccount, ":")) == 2 {
		// login:pass
		parts := strings.Split(twitterAccount, ":")
		login := parts[0]
		password := parts[1]
		grabberInstance := grabber.Grabber{}
		ok := grabberInstance.InitGrabber(index, login, password, proxy, config, queryIDs)
		if ok == false {
			return
		}
		authToken := grabberInstance.GrabTheCookie()
		err := extra.ReplaceLineInFile("data/accounts.txt", twitterAccount, twitterAccount+":"+authToken)
		if err != nil {
			return
		}
		twitterAccount = authToken
		return
	} else {
		extra.Logger{}.Error("%d | Wrong account format.", index+1)
		return
	}

	twitter := instance.Twitter{}
	ok, reason := twitter.InitTwitter(index+1, twitterAccount, proxy, config, queryIDs)
	if ok == false {
		if reason == "locked" {
			lockedAccounts <- twitterAccount
			return
		}
		failedAccounts <- twitterAccount
		return
	}

	if slices.Contains(userChoice, "Follow") {
		for _, user := range tasks.FollowTargets {
			ok = twitter.Follow(user)
			if ok == false {
				report = true
			}
			extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
		}
	}
	if slices.Contains(userChoice, "Retweet") {
		for _, link := range tasks.RetweetLinks {
			ok = twitter.Retweet(link)
			if ok == false {
				report = true
			}
			extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
		}
	}
	if slices.Contains(userChoice, "Like") {
		for _, link := range tasks.LikeLinks {
			ok = twitter.Like(link)
			if ok == false {
				report = true
			}
			extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
		}
	}
	if slices.Contains(userChoice, "Tweet") {
		if config.Data.Random == "no" {
			ok = twitter.Tweet(txtData.Tweets[index])
		} else {
			ok = twitter.Tweet(utils.GetRandomSliceElement(txtData.Tweets))
		}
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Quote Tweet") {
		for _, link := range tasks.QuoteTweetLinks {
			if config.Data.Random == "no" {
				ok = twitter.TweetQuote(txtData.Tweets[index], link)
			} else {
				ok = twitter.TweetQuote(utils.GetRandomSliceElement(txtData.Tweets), link)
			}
			if ok == false {
				report = true
			}
			extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
		}
	}
	if slices.Contains(userChoice, "Tweet with Picture") {
		if config.Data.Random == "no" {
			ok = twitter.TweetWithPicture(txtData.Tweets[index], txtData.Pictures[index])
		} else {
			ok = twitter.TweetWithPicture(utils.GetRandomSliceElement(txtData.Tweets), utils.GetRandomSliceElement(txtData.Pictures))
		}
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Comment") {
		for _, link := range tasks.CommentTweetLinks {
			if config.Data.Random == "no" {
				ok = twitter.Comment(txtData.Comments[index], link)
			} else {
				ok = twitter.Comment(utils.GetRandomSliceElement(txtData.Comments), link)
			}
			if ok == false {
				report = true
			}
			extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
		}
	}
	if slices.Contains(userChoice, "Comment with Picture") {
		for _, link := range tasks.CommentTweetLinks {
			if config.Data.Random == "no" {
				ok = twitter.CommentWithPicture(txtData.Comments[index], link, txtData.Pictures[index])
			} else {
				ok = twitter.CommentWithPicture(utils.GetRandomSliceElement(txtData.Comments), link, utils.GetRandomSliceElement(txtData.Pictures))
			}
			if ok == false {
				report = true
			}
			extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
		}
	}
	if slices.Contains(userChoice, "Change Description") {
		ok = twitter.ChangeDescription(txtData.Descriptions[index])
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)

	}
	if slices.Contains(userChoice, "Change Username") {
		ok = twitter.ChangeUsername(txtData.Usernames[index])
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)

	}
	if slices.Contains(userChoice, "Change Name") {
		ok = twitter.ChangeName(txtData.Names[index])
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Change Background") {
		if config.Data.Random == "no" {
			ok = twitter.ChangeBackground(txtData.Pictures[index])
		} else {
			ok = twitter.ChangeBackground(utils.GetRandomSliceElement(txtData.Pictures))
		}
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Change Password") {
		newAuthToken := twitter.ChangePassword(txtData.OldPasswords[index], txtData.NewPasswords[index])
		if newAuthToken == "" {
			report = true
		} else {
			err := extra.ReplaceLineInFile("data/accounts.txt", twitterAccount, newAuthToken)
			if err != nil {
				return
			}
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Change Birthdate") {
		ok = twitter.ChangeBirthdate(txtData.Birthdate[index])
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Change Location") {
		ok = twitter.ChangeLocation(txtData.Locations[index])
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Change profile Picture") {
		if config.Data.Random == "no" {
			ok = twitter.ChangeProfilePicture(txtData.Pictures[index])
		} else {
			ok = twitter.ChangeProfilePicture(utils.GetRandomSliceElement(txtData.Pictures))
		}
		if ok == false {
			report = true
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Check if account is valid") {
		var username string
		creationDate := ""
		ok, username, creationDate = twitter.CheckSuspended()
		dateToWrite := fmt.Sprintf("%s:%s", username, creationDate)

		if username != "" {
			creationDateChan <- dateToWrite
		}

		if ok == false {
			report = true
		} else {
			usernames <- username
		}
		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	if slices.Contains(userChoice, "Mutual Subscription") {
		successCount := 0
		mutualSubsUsernames.Range(func(key, value interface{}) bool {
			if successCount >= tasks.MutualSubsFollowersCount {
				return true
			} else if twitter.Username == key.(string) {
				return true
			}
			currentVal := value.(int)
			if currentVal < tasks.MutualSubsFollowersCount {
				newVal := currentVal + 1
				mutualSubsUsernames.Store(key, newVal)
				if twitter.Follow(key.(string)) {
					successCount++
					extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
				} else {
					mutualSubsUsernames.Store(key, newVal-1)
				}
			}
			return true
		})
	}
	if slices.Contains(userChoice, "Vote in the poll") {
		ok = twitter.VotePoll(tasks.Votes.Links[0], tasks.Votes.Choice)
		if ok == false {
			report = true
		}

		extra.RandomSleep(config.Info.PauseBetweenTasks.Start, config.Info.PauseBetweenTasks.End)
	}
	//if slices.Contains(userChoice, "Unfreeze Accounts") {
	//	ok = twitter.Unfreeze()
	//	if ok == false {
	//		report = true
	//	}
	//}

	if report == true {
		failedAccounts <- twitterAccount
	}
}

func MobileProxyWrapper(data MobileProxyData) {
	for i := range data.Indexes {
		Process(i, data.Accounts[i], data.Proxy, data.Config, data.UserChoice, data.Tasks, data.TxtData, data.QueryIDs, data.FailedAccounts, data.LockedAccounts, data.Usernames, data.MutualSubsUsernames, data.CreationDateChan)
		extra.RandomSleep(data.Config.Info.PauseBetweenAccounts.Start, data.Config.Info.PauseBetweenAccounts.End)

		extra.ChangeProxyURL(data.IPChangeLink)
		extra.Logger{}.Success("Successfully changed the IP address of mobile proxies")
		time.Sleep(time.Duration(data.Config.Proxy.ChangeIPPause) * time.Second)

	}
}

type MobileProxyData struct {
	Proxy               string
	IPChangeLink        string
	Indexes             chan int
	Accounts            []string
	Config              extra.Config
	UserChoice          []string
	Tasks               extra.TasksTargetInfo
	TxtData             extra.TxtTasksData
	QueryIDs            extra.QueryIDs
	FailedAccounts      chan string
	LockedAccounts      chan string
	Usernames           chan string
	MutualSubsUsernames *sync.Map
	CreationDateChan    chan string
}
