package prompter

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

func PromptString(msg string) (string, error) {
	prompt := promptui.Prompt{
		Label: msg,
	}
	value, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return value, nil
}

func PromptConfirm(msg string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
	}
	value, err := prompt.Run()
	fmt.Println("Confirm result: " + value)
	fmt.Println(err != nil)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return false, err
	}
	return value != "N", nil
}

func PromptSelect(msg string, options []string) (string, error) {
	prompt := promptui.Select{
		Label: msg,
		Items: options,
	}
	_, choice, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return choice, nil
}
