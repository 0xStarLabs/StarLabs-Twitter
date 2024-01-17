package extra

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Config struct {
	Info struct {
		TwoCaptchaKey           string `yaml:"2captcha_api_key"`
		MaxTasksRetries         int    `yaml:"max_tasks_retries"`
		AutoUnfreeze            string `yaml:"auto_unfreeze"`
		PauseBetweenTasksRaw    string `yaml:"pause_between_tasks"`
		PauseBetweenAccountsRaw string `yaml:"pause_between_accounts"`
		PauseBetweenTasks       struct {
			Start int
			End   int
		} `yaml:"-"`
		PauseBetweenAccounts struct {
			Start int
			End   int
		} `yaml:"-"`
		AccountRange struct {
			Start int
			End   int
		} `yaml:"-"`
		AccountRangeRaw string `yaml:"account_range"`
	} `yaml:"info"`
	Proxy struct {
		MobileProxy   string `yaml:"mobile_proxy"`
		ChangeIPPause int    `yaml:"change_ip_pause"`
	} `yaml:"proxy"`
	Data struct {
		Random string `yaml:"random"`
	} `yaml:"data"`
}

func ReadConfig() Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	accountRangeParts := strings.Split(config.Info.AccountRangeRaw, "-")
	if len(accountRangeParts) == 2 {
		var start, end int
		_, errStart := fmt.Sscanf(accountRangeParts[0], "%d", &start)
		_, errEnd := fmt.Sscanf(accountRangeParts[1], "%d", &end)
		if errStart == nil && errEnd == nil {
			config.Info.AccountRange.Start = start
			config.Info.AccountRange.End = end
		} else {
			Logger{}.Error("Failed to read account_range from config:", errStart, errEnd)
		}
	} else {
		Logger{}.Error("Wrong account_range format")
	}
	pauseBetweenTasksParts := strings.Split(config.Info.PauseBetweenTasksRaw, "-")
	if len(pauseBetweenTasksParts) == 2 {
		var start, end int
		_, errStart := fmt.Sscanf(pauseBetweenTasksParts[0], "%d", &start)
		_, errEnd := fmt.Sscanf(pauseBetweenTasksParts[1], "%d", &end)
		if errStart == nil && errEnd == nil {
			config.Info.PauseBetweenTasks.Start = start
			config.Info.PauseBetweenTasks.End = end
		} else {
			Logger{}.Error("Failed to read pause_between_tasks from config:", errStart, errEnd)
		}
	} else {
		Logger{}.Error("Wrong pause_between_tasks format")
	}

	// Парсинг PauseBetweenAccounts
	pauseBetweenAccountsParts := strings.Split(config.Info.PauseBetweenAccountsRaw, "-")
	if len(pauseBetweenAccountsParts) == 2 {
		var start, end int
		_, errStart := fmt.Sscanf(pauseBetweenAccountsParts[0], "%d", &start)
		_, errEnd := fmt.Sscanf(pauseBetweenAccountsParts[1], "%d", &end)
		if errStart == nil && errEnd == nil {
			config.Info.PauseBetweenAccounts.Start = start
			config.Info.PauseBetweenAccounts.End = end
		} else {
			Logger{}.Error("Failed to read pause_between_accounts from config:", errStart, errEnd)
		}
	} else {
		Logger{}.Error("Wrong pause_between_accounts format")
	}

	return config
}

func ReadQueryIDs() QueryIDs {
	data, err := os.ReadFile("extra/query_ids.yaml")
	if err != nil {
		panic(err)
	}
	// Unmarshal YAML content into struct
	var config QueryIDs
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}

type QueryIDs struct {
	BearerToken string `yaml:"bearer_token"`
	LikeID      string `yaml:"like"`
	RetweetID   string `yaml:"retweet"`
	TweetID     string `yaml:"tweet"`
}
