package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Range represents a min-max range of values
type Range struct {
	Start int `yaml:",flow"`
	End   int `yaml:",flow"`
}

// UnmarshalYAML implements custom unmarshaling for Range from [min, max] format
func (r *Range) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var arr []int
	if err := unmarshal(&arr); err != nil {
		return err
	}
	if len(arr) != 2 {
		return fmt.Errorf("range must have exactly 2 values, got %d", len(arr))
	}
	r.Start = arr[0]
	r.End = arr[1]
	return nil
}

// Settings represents the settings section of the config
type Settings struct {
	Threads              int   `yaml:"THREADS"`
	Retries              int   `yaml:"RETRIES"`
	PauseBetweenRetries  Range `yaml:"PAUSE_BETWEEN_RETRIES"`
	PauseBetweenAccounts Range `yaml:"PAUSE_BETWEEN_ACCOUNTS"`
	AccountsRange        Range `yaml:"ACCOUNTS_RANGE"`
}

type MutualSubscription struct {
	FollowersForEveryAccount Range `yaml:"FOLLOWERS_FOR_EVERY_ACCOUNT"`
}

// Config represents the root configuration structure
type Config struct {
	Settings           Settings           `yaml:"SETTINGS"`
	MutualSubscription MutualSubscription `yaml:"MUTUAL_SUBSCRIPTION"`
}

// ReadConfig reads and parses the YAML configuration file
//
// Usage example:
//
//	config, err := utils.ReadConfig("./data/config.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Access config values
//	retries := config.Settings.Retries
//	minPause := config.Settings.PauseBetweenRetries.Start
//	maxPause := config.Settings.PauseBetweenRetries.End
//
//	// Access account range
//	startAccount := config.Settings.AccountsRange.Start
//	endAccount := config.Settings.AccountsRange.End
func ReadConfig(filePath string) (*Config, error) {
	// Read the config file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validate the config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// validateConfig checks if the configuration values are valid
func validateConfig(config *Config) error {
	// Check retries
	if config.Settings.Retries < 1 {
		return fmt.Errorf("retries must be at least 1, got %d", config.Settings.Retries)
	}

	// Validate ranges
	ranges := []struct {
		name  string
		value Range
	}{
		{"PAUSE_BETWEEN_RETRIES", config.Settings.PauseBetweenRetries},
		{"PAUSE_BETWEEN_ACCOUNTS", config.Settings.PauseBetweenAccounts},
		{"ACCOUNTS_RANGE", config.Settings.AccountsRange},
		{"FOLLOWERS_FOR_EVERY_ACCOUNT", config.MutualSubscription.FollowersForEveryAccount},
	}

	for _, r := range ranges {
		if r.value.Start > r.value.End {
			return fmt.Errorf("%s range start (%d) cannot be greater than end (%d)",
				r.name, r.value.Start, r.value.End)
		}
		if r.value.Start < 0 {
			return fmt.Errorf("%s range start cannot be negative, got %d",
				r.name, r.value.Start)
		}
	}

	return nil
}
