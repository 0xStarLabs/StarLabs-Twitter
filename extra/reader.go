package extra

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadTxtFile(fileName string, filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		return []string{}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var items []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		items = append(items, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return []string{}
	}

	Logger{}.Info("Successfully loaded %d %s.", len(items), fileName)
	return items
}

func ReadPictures(path string) ([]string, bool) {
	var encodedImages []string

	files, err := os.ReadDir(path)
	if err != nil {
		Logger{}.Error("Failed to read pictures dir: %s", err)
		return encodedImages, false
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".jpeg") {
			filePath := filepath.Join(path, file.Name())

			imageFile, err := os.ReadFile(filePath)
			if err != nil {
				Logger{}.Error("Failed to read the picture: %s", err)
			}

			encodedImage := base64.StdEncoding.EncodeToString(imageFile)
			encodedImages = append(encodedImages, encodedImage)
		}
	}

	return encodedImages, true
}

func ReadUserChoice(tasksList []string, textToShow string) []string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\n" + textToShow + ": ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	userInput = strings.ReplaceAll(userInput, ",", " ")
	var chosenTasks []string

	for _, task := range strings.Fields(userInput) {
		taskIndex, err := strconv.Atoi(strings.TrimSpace(task))
		if err != nil || taskIndex < 1 || taskIndex > len(tasksList) {
			Logger{}.Error("Invalid input: ", task)
			continue
		}
		chosenTasks = append(chosenTasks, tasksList[taskIndex-1])
	}

	return chosenTasks
}

func GetTxtDataByTasks(tasks []string) (TxtTasksData, bool) {
	changeData := TxtTasksData{}

	for _, task := range tasks {
		if task == "Tweet" || task == "Quote Tweet" || task == "Tweet with Picture" {
			data := ReadTxtFile("tweets", "data/tweets.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.Tweets = data
		}

		if task == "Comment" || task == "Comment with Picture" {
			data := ReadTxtFile("comments", "data/comments.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.Comments = data
		}

		if task == "Tweet with Picture" || task == "Comment with Picture" || task == "Change Background" || task == "Change profile Picture" {
			data, ok := ReadPictures("data/pictures")
			if ok == false {
				return changeData, false
			}
			changeData.Pictures = data
		}

		if task == "Change Password" {
			data := ReadTxtFile("current passwords", "data/self/current_passwords.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.OldPasswords = data

			data = ReadTxtFile("new passwords", "data/self/new_passwords.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.NewPasswords = data
		}

		if task == "Change Description" {
			data := ReadTxtFile("descriptions", "data/self/description.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.Descriptions = data
		}

		if task == "Change Username" {
			data := ReadTxtFile("usernames", "data/self/usernames.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.Usernames = data
		}

		if task == "Change Name" {
			data := ReadTxtFile("names", "data/self/names.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.Names = data
		}

		if task == "Change Birthdate" {
			data := ReadTxtFile("birthdate", "data/self/birthdate.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.Birthdate = data
		}

		if task == "Change Location" {
			data := ReadTxtFile("locations", "data/self/location.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.Locations = data
		}

		if task == "Mutual Subscription" {
			data := ReadTxtFile("usernames", "data/my_usernames.txt")
			if len(data) == 0 {
				return changeData, false
			}
			changeData.MyUsernames = data
		}

	}

	return changeData, true
}

func GetTasksTargetInfo(tasks []string) (TasksTargetInfo, bool) {
	tasksInfo := TasksTargetInfo{}

	for _, task := range tasks {
		if task == "Follow" {
			tasksInfo.FollowTargets = userInputToSlice("Paste the usernames you want to follow")
		}
		if task == "Retweet" {
			tasksInfo.RetweetLinks = userInputToSlice("Paste your link for retweet")
		}
		if task == "Like" {
			tasksInfo.LikeLinks = userInputToSlice("Paste your link for like")
		}
		if task == "Quote Tweet" {
			tasksInfo.QuoteTweetLinks = userInputToSlice("Paste your link for quote tweet")
		}
		if task == "Comment" || task == "Comment with Picture" {
			tasksInfo.CommentTweetLinks = userInputToSlice("Paste your link for the comment")
		}
		if task == "Check account messages" {
			number, ok := UserInputInteger("Enter the number of days since the last tweet")
			if ok == false {
				return tasksInfo, false
			}
			tasksInfo.CheckMessagesLastDays = number
		}
		if task == "Mutual Subscription" {
			number, ok := UserInputInteger("How many followers should each account have")
			if ok == false {
				return tasksInfo, false
			}
			tasksInfo.MutualSubsFollowersCount = number
		}
	}
	return tasksInfo, true
}

func userInputToSlice(textToShow string) []string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\n" + textToShow + ": ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	userInput = strings.ReplaceAll(userInput, ",", " ")

	return strings.Fields(userInput)
}

func UserInputInteger(textToShow string) (int, bool) {
	var input int

	fmt.Print(textToShow + ": ")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	input, err := strconv.Atoi(userInput)
	if err != nil {
		Logger{}.Error("Wrong input.")
		return input, false
	}
	return input, true
}

func NoProxies() bool {
	var userChoice int
	fmt.Println("No proxies were detected. Do you want to continue without proxies? (1 or 2)")
	fmt.Println("[1] Yes")
	fmt.Println("[2] No")
	fmt.Print(">> ")
	_, err := fmt.Scan(&userChoice)
	if err != nil {
		Logger{}.Error("Wrong input. Enter the number.")
		panic(err)
	}

	return userChoice == 1
}

type TxtTasksData struct {
	Tweets       []string
	Comments     []string
	Pictures     []string
	OldPasswords []string
	NewPasswords []string
	Descriptions []string
	Usernames    []string
	Names        []string
	Birthdate    []string
	Locations    []string
	MyUsernames  []string
}

type TasksTargetInfo struct {
	FollowTargets            []string
	RetweetLinks             []string
	LikeLinks                []string
	QuoteTweetLinks          []string
	CommentTweetLinks        []string
	CheckMessagesLastDays    int
	MutualSubsFollowersCount int
}
