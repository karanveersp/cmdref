package cmdref

import (
	"cmdref/prompter"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Action int

const (
	Create Action = iota
	Update
	Remove
	View
	Import
	Exit
)

var Actions = []string{"Create", "Update", "Remove", "View", "Import", "Exit"}

func toAction(s string) (Action, error) {
	switch s {
	case "Create":
		return Create, nil
	case "Update":
		return Update, nil
	case "Remove":
		return Remove, nil
	case "View":
		return View, nil
	case "Import":
		return Import, nil
	case "Exit":
		return Exit, nil
	default:
		return -1, errors.New("unrecognized action " + s)
	}
}

type Command struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	Platform    string `json:"platform"`
	Description string `json:"description"`
}

func (cmd Command) String() string {
	return fmt.Sprintf("\nName: %s\nCommand: %s\nDescription: %s\nPlatform: %s\n",
		cmd.Name, cmd.Command, cmd.Description, cmd.Platform)
}

const CmdDirName = "cmdref"
const CmdFileName = "cmdref.json"

func CmdFilePath() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	cmdStoreDir := filepath.Join(d, CmdDirName)
	err = CreateDirIfNotExists(cmdStoreDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(cmdStoreDir, CmdFileName), nil
}

func UpdateFile(cmdsFilePath string, cmdMap map[string]Command) error {
	var commands []Command
	for _, v := range cmdMap {
		commands = append(commands, v)
	}
	jsonEncoding, err := json.Marshal(commands)
	if err != nil {
		return err
	}
	err = os.WriteFile(cmdsFilePath, jsonEncoding, 0666)
	if err != nil {
		return err
	}
	return nil
}

func UpdateHandler(cmdMap map[string]Command) (map[string]Command, error) {
	entries := keys[Command](cmdMap)

	selection, err := prompter.PromptSelect("Select a command to update", entries)
	if err != nil {
		return nil, err
	}

	fmt.Print(cmdMap[selection])

	cmd, err := createCommandWithName(selection)
	if err != nil {
		return nil, err
	}

	newMap := copyMap(cmdMap)
	newMap[selection] = cmd
	return newMap, nil
}

func DeleteHandler(cmdMap map[string]Command) (map[string]Command, error) {
	entries := keys[Command](cmdMap)
	if len(cmdMap) == 0 {
		fmt.Println("No commands to delete")
		return cmdMap, nil
	}
	selection, err := prompter.PromptSelect("Select command to delete", entries)
	if err != nil {
		return nil, err
	}

	confirm, err := prompter.PromptConfirm(fmt.Sprintf("Are you sure you want to delete '%s'", selection))
	if err != nil {
		return nil, fmt.Errorf("error while prompting delete confirmation - %v\n", err)
	}

	if confirm {
		delete(cmdMap, selection)
	}
	return copyMap(cmdMap), nil

}

func ViewHandler(cmdMap map[string]Command) error {
	if len(cmdMap) == 0 {
		fmt.Println("No existing commands found")
		return nil
	}
	var entries []string
	for name := range cmdMap {
		entries = append(entries, name)
	}
	selectedItem, err := prompter.PromptSelect("Select a command", entries)
	if err != nil {
		return err
	}
	cmd := cmdMap[selectedItem]
	fmt.Print(cmd)
	return nil
}

func ProcessAction(cmdsFilePath string, cmdMap map[string]Command, action Action) (map[string]Command, error) {
	switch action {
	case Create:
		newMap, err := CreateHandler(cmdMap)
		if err != nil {
			return nil, err
		}
		err = UpdateFile(cmdsFilePath, newMap)
		if err != nil {
			return nil, err
		}
		return newMap, nil
	case View:
		err := ViewHandler(cmdMap)
		if err != nil {
			return nil, err
		}
		return cmdMap, nil
	case Update:
		newMap, err := UpdateHandler(cmdMap)
		if err != nil {
			return nil, err
		}
		err = UpdateFile(cmdsFilePath, newMap)
		if err != nil {
			return nil, err
		}
		return newMap, nil
	case Remove:
		newMap, err := DeleteHandler(cmdMap)
		if err != nil {
			return nil, err
		}
		err = UpdateFile(cmdsFilePath, newMap)
		if err != nil {
			return nil, err
		}
		return newMap, err
	default:
		return nil, errors.New("unrecognized action")
	}
}

func GetSelectedAction() (Action, error) {
	action, err := prompter.PromptSelect("Select action", Actions)
	if err != nil {
		return -1, err
	}
	mappedAction, err := toAction(action)
	if err != nil {
		return -1, err
	}
	return mappedAction, nil
}

func CreateHandler(cmdMap map[string]Command) (map[string]Command, error) {
	cmd, err := createCommand()
	if err != nil {
		return nil, err
	}
	newMap := copyMap(cmdMap)
	newMap[cmd.Name] = cmd
	return newMap, nil
}

func ParseCommands(cmdProvider func() ([]byte, error)) (map[string]Command, error) {
	cmdMap := make(map[string]Command)
	cmdData, err := cmdProvider()
	if err != nil {
		return nil, err
	}
	var commands []Command
	err = json.Unmarshal(cmdData, &commands)
	if err != nil {
		return nil, err
	}
	for _, command := range commands {
		cmdMap[command.Name] = command
	}
	return cmdMap, nil
}

func CreateDirIfNotExists(dpath string) error {
	if stat, err := os.Stat(dpath); err == nil && stat.IsDir() {
		return nil // directory exists
	}

	//0755 Commonly used on web servers. The owner can read, write, execute. Everyone else can read and execute but not modify the file.
	//
	//0777 Everyone can read write and execute. On a web server, it is not advisable to use ‘777’ permission for your files and folders, as it allows anyone to add malicious code to your server.
	//
	//0644 Only the owner can read and write. Everyone else can only read. No one can execute the file.
	//
	//0655 Only the owner can read and write, but not execute the file. Everyone else can read and execute, but cannot modify the file.
	err := os.MkdirAll(dpath, 0777)
	if err != nil {
		return err
	}
	return nil
}

func copyMap[T interface{}](m map[string]T) map[string]T {
	newMap := make(map[string]T)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

func keys[T interface{}](m map[string]T) []string {
	var entries []string
	for k := range m {
		entries = append(entries, k)
	}
	return entries
}
func createCommand() (Command, error) {
	name, err := prompter.PromptString("Command name")
	if err != nil {
		return Command{}, err
	}
	return createCommandWithName(name)
}

func createCommandWithName(name string) (Command, error) {
	command, err := prompter.PromptString("Command")
	if err != nil {
		return Command{}, err
	}
	platform, err := prompter.PromptString("Platform")
	if err != nil {
		return Command{}, err
	}
	description, err := prompter.PromptString("Description")
	if err != nil {
		return Command{}, err
	}
	return Command{Name: name, Command: command, Platform: platform, Description: description}, nil
}
