package prompter

import (
	"github.com/AlecAivazis/survey/v2"
	"strings"
)

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
