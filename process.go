package main

import (
	"slices"
	"strings"
	"sync"
	"time"
	"twitter/extra"
	"twitter/instance"
	"twitter/utils"
)

func Process(index int, twitterAccount string, proxy string, config extra.Config, userChoice []string, tasks extra.TasksTargetInfo, txtData extra.TxtTasksData, queryIDs extra.QueryIDs, failedAccounts chan<- string, lockedAccounts chan<- string, usernames chan<- string, mutualSubsUsernames *sync.Map) {
	var report bool

	if len(twitterAccount) == 40 && !strings.Contains(twitterAccount, ":") {
		// auth_token
	} else if strings.Contains(twitterAccount, ":") && len(strings.Split(twitterAccount, ":")) >= 3 && len(strings.Split(twitterAccount, ":")[2]) == 40 {
		// login:pass:auth_token:json_cookies
		twitterAccount = strings.Split(twitterAccount, ":")[2]
		//} else if strings.Contains(twitterAccount, ":") && len(strings.Split(twitterAccount, ":")) == 2 {
		//	// login:pass
		//	parts := strings.Split(twitterAccount, ":")
		//	login := parts[0]
		//	password := parts[1]
		//	grabberInstance := grabber.Grabber{}
		//	ok := grabberInstance.InitGrabber(index, login, password, proxy, config, queryIDs)
		//	if ok == false {
		//		return
		//	}
		//	grabberInstance.GrabTheCookie()
		//	return
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
		ok, username = twitter.CheckSuspended()
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
		Process(i, data.Accounts[i], data.Proxy, data.Config, data.UserChoice, data.Tasks, data.TxtData, data.QueryIDs, data.FailedAccounts, data.LockedAccounts, data.Usernames, data.MutualSubsUsernames)
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
}
