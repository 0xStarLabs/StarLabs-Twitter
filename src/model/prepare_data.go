package model

import (
	"errors"
	"starlabs-twitter/src/utils"
	"strconv"
)

func GetTxtDataByTasks(tasks []string) (utils.TxtTasksData, error) {
	changeData := utils.TxtTasksData{}
	err := errors.New("no data")

	for _, task := range tasks {
		if task == "Tweet" || task == "Quote Tweet" || task == "Tweet with Picture" {
			data := utils.ReadTxtFile("tweets", "data/tweet_text.txt")
			if len(data) == 0 {
				return changeData, err
			}
			changeData.Tweets = data
		}

		if task == "Comment" || task == "Comment with Picture" {
			data := utils.ReadTxtFile("comments", "data/comment_text.txt")
			if len(data) == 0 {
				return changeData, err
			}
			changeData.Comments = data
		}

		if task == "Tweet with Picture" || task == "Comment with Picture" {
			data, ok := utils.ReadPictures("data/images")
			if ok == false {
				return changeData, err
			}
			changeData.Pictures = data
		}
	}

	return changeData, nil
}

func GetTasksTargetInfo(tasks []string) utils.TasksTargetInfo {
	tasksInfo := utils.TasksTargetInfo{}

	for _, task := range tasks {
		if task == "Follow" {
			tasksInfo.FollowTargets = utils.UserInputToSlice("Paste the usernames you want to follow")
		}
		if task == "Unfollow" {
			tasksInfo.UnfollowTargets = utils.UserInputToSlice("Paste the usernames you want to unfollow")
		}
		if task == "Retweet" {
			tasksInfo.RetweetLinks = utils.UserInputToSlice("Paste your link for retweet")
		}
		if task == "Like" {
			tasksInfo.LikeLinks = utils.UserInputToSlice("Paste your link for like")
		}
		if task == "Quote Tweet" {
			tasksInfo.QuoteTweetLinks = utils.UserInputToSlice("Paste your link for quote tweet")
		}
		if task == "Comment" || task == "Comment with Picture" {
			tasksInfo.CommentTweetLinks = utils.UserInputToSlice("Paste your link for the comment")
		}

		// if task == "Mutual Subscription" {
		// 	number, ok := UserInputInteger("How many followers should each account have")
		// 	if ok == false {
		// 		return tasksInfo, false
		// 	}
		// 	tasksInfo.MutualSubsFollowersCount = number
		// }

		if task == "Vote in the poll" {
			tasksInfo.Votes.Links = utils.UserInputToSlice("Paste your link to the poll")
			answer, _ := utils.UserInputInteger("Enter the sequential number of the answer ( 1 2 3 etc.)")
			tasksInfo.Votes.Choice = strconv.Itoa(answer)
		}
	}
	return tasksInfo
}
