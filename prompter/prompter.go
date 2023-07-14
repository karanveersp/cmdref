package prompter

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// PromptString asks the user for a string input with the given message.
func PromptString(msg string) (string, error) {
	result := ""
	if !strings.HasSuffix(msg, ":") {
		msg += ":"
	}
	prompt := &survey.Input{
		Message: msg,
	}
	err := survey.AskOne(prompt, &result, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}
	return result, nil
}

// PromptConfirm asks the user a y/n question with the given message.
func PromptConfirm(msg string) (bool, error) {

	value := false
	prompt := &survey.Confirm{
		Message: msg,
	}
	err := survey.AskOne(prompt, &value, survey.WithValidator(survey.Required))
	if err != nil {
		return false, err
	}
	return value, nil
}

// PromptSelect asks the user to select an option from the given message and options.
func PromptSelect(msg string, options []string) (string, error) {
	choice := ""
	prompt := &survey.Select{
		Message: msg,
		Options: options,
	}
	err := survey.AskOne(prompt, &choice, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}
	return choice, nil
}
