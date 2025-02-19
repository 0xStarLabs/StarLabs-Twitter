package main

import (
	"log"
	"starlabs-twitter/src/model"
	"starlabs-twitter/src/utils"
)

func main() {
	for {
		utils.ShowLogo()
		utils.ShowDevInfo()

		options()
	}
}


func options() {
	userSelectedOptions, err := utils.MultiSelectMenu("Select options:", utils.MainMenuOptions)
	if err != nil {
		log.Fatal(err)
	}

	config, err := utils.ReadConfig("data/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	accounts, err := utils.ReadAccountsFromExcel("data/accounts.xlsx", config.Settings.AccountsRange.Start, config.Settings.AccountsRange.End)
	if err != nil {
		log.Fatal(err)
	}

	taskData := model.GetTasksTargetInfo(userSelectedOptions)

	fileData, err := model.GetTxtDataByTasks(userSelectedOptions)
	if err != nil {
		utils.Logger{Prefix: "Start"}.Error("Failed to load task files: %v", err)
		return
	}

	if utils.Contains(userSelectedOptions, "Mutual Subscription") {
		model.MutualSubscription(accounts, *config)
	} else {
		model.Start(accounts, userSelectedOptions, *config, taskData, fileData)
	}
}
