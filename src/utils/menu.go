package utils

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// MultiSelectMenu displays an interactive menu and returns selected items
//
// Usage example:
//
//	options := []string{"Option 1", "Option 2", "Option 3"}
//	selected, err := utils.MultiSelectMenu("Select options:", options)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Selected: %v\n", selected)
func MultiSelectMenu(prompt string, options []string) ([]string, error) {
	blue := color.New(color.FgHiBlue)
	helpColor := color.New(color.FgHiYellow, color.Bold)

	// Print empty line for better spacing
	fmt.Println()

	// Print help text
	helpText := []string{
		"Navigation:",
		"  ↑/↓: Move cursor",
		"  Space: Select/Unselect item",
		"  Enter: Confirm selection",
		"",
	}

	for _, line := range helpText {
		helpColor.Println(line)
	}

	// Create the survey prompt
	var selected []string
	prompt = blue.Sprintf(prompt)

	// Configure the survey question
	question := &survey.MultiSelect{
		Message:  prompt,
		Options:  options,
		PageSize: 30,
	}

	// Add custom styling
	surveyOpts := []survey.AskOpt{
		survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = ""
			icons.Question.Format = "blue"
			icons.SelectFocus.Text = blue.Sprintf("→")
			icons.SelectFocus.Format = "blue"
			icons.MarkedOption.Text = blue.Sprintf("[✓]")
			icons.MarkedOption.Format = "blue"
			icons.UnmarkedOption.Text = "[ ]"
			icons.HelpInput.Text = ""
		}),
		survey.WithShowCursor(true),
		survey.WithPageSize(20),
	}

	// Ask the question
	err := survey.AskOne(question, &selected, surveyOpts...)
	if err != nil {
		return nil, fmt.Errorf("menu selection error: %w", err)
	}

	// Print empty line after selection
	fmt.Println()

	return selected, nil
}

// DisplayNumberedMenu shows a menu with numbered options
func DisplayNumberedMenu(title string, options []string) ([]string, error) {
	// Create numbered options
	numberedOptions := make([]string, len(options))
	for i, opt := range options {
		numberedOptions[i] = fmt.Sprintf("%d. %s", i+1, opt)
	}

	// Display the menu with numbers
	return MultiSelectMenu(title, numberedOptions)
}
