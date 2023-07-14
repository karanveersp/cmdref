package cmdref

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/karanveersp/cmdref/pkg/prompter"

	"github.com/fatih/color"
	"github.com/karanveersp/store"
)

// CmdDirName is the name of the command directory for config storage.
const CmdDirName = "cmdref"

// CmdFileName is the name of the file for command storage.
const CmdFileName = "cmdref.json"

// CmdFileOperater defines an interface for saving/loading commands from
// the commands file.
// This is needed to create mock operations for unit tests, along with
// a live implementation for the actual file.
type CmdFileOperater interface {
	Load() ([]Command, error)
	LoadExternal(filepath string) ([]Command, error)
	Save(commands []Command) error
	GetFilePath() string
}

// CmdFileOps implements the CommandsFileOperator for live/side-effectful
// operations on the commands file.
type CmdFileOps struct{}

// NewCmdFileOps creates and returns a new file operations struct.
func NewCmdFileOps() CmdFileOps {
	// initialize the application dir for the config store.
	store.Init(CmdDirName)
	return CmdFileOps{}
}

// Load parses the commands file and returns the list of commands.
func (f *CmdFileOps) Load() ([]Command, error) {
	var cmds []Command
	err := store.Load(CmdFileName, &cmds)
	if err != nil {
		return nil, err
	}
	return cmds, nil
}

// LoadExternal parses the commands file and returns the list of commands.
func (f *CmdFileOps) LoadExternal(filepath string) ([]Command, error) {
	var cmds []Command
	provider := func() ([]byte, error) {
		data, err := os.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	cmds, err := parseCommands(provider)
	if err != nil {
		return nil, err
	}
	return cmds, nil
}

// Save writes the list of commands to the commands file.
func (f *CmdFileOps) Save(commands []Command) error {
	return store.Save(CmdFileName, &commands)
}

// GetFilePath returns the absolute path to the commands file.
func (f *CmdFileOps) GetFilePath() string {
	cmdStoreDir := store.GetApplicationDirPath()
	return filepath.Join(cmdStoreDir, CmdFileName)
}

// Action represents any action supported by the app.
type Action int

const (
	// Create action for creating a command
	Create Action = iota
	// Update action for updating an existing command
	Update
	// Remove action for removing a command
	Remove
	// View action for viewing a command
	View
	// Import action for importing commands from another file
	Import
	// Exit action for quitting the app
	Exit
)

// Actions is a list of strings representing all actions.
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

// Command represents a command to store and query.
type Command struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	Platform    string `json:"platform"`
	Description string `json:"description"`
}

func (cmd Command) String() string {
	return fmt.Sprintf("Name: %s\nDescription: %s\nPlatform: %s\nCommand:\n%s\n",
		cmd.Name, cmd.Description, cmd.Platform, cmd.Command)
}

func updateFile(cmdMap map[string]Command, fileOps CmdFileOperater) error {
	var commands []Command
	for _, v := range cmdMap {
		commands = append(commands, v)
	}
	return fileOps.Save(commands)
}

// UpdateHandler handles the update action.
func UpdateHandler(cmdMap map[string]Command) (map[string]Command, error) {
	if len(cmdMap) == 0 {
		fmt.Println("No existing commands.")
		return cmdMap, nil
	}

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

// DeleteHandler handles the delete action.
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

	fmt.Println(cmdMap[selection])

	confirm, err := prompter.PromptConfirm(fmt.Sprintf("Are you sure you want to delete '%s'", selection))
	if err != nil {
		return nil, fmt.Errorf("error while prompting delete confirmation - %v", err)
	}

	if confirm {
		delete(cmdMap, selection)
	}
	return copyMap(cmdMap), nil

}

// ViewHandler handles viewing commands.
func ViewHandler(cmdMap map[string]Command) error {
	if len(cmdMap) == 0 {
		fmt.Println("No existing commands found")
		return nil
	}
	var entries []string
	entryToCommandName := make(map[string]string)
	for name := range cmdMap {
		entryName := fmt.Sprintf("%s - %s", cmdMap[name].Platform, name)
		entries = append(entries, entryName)
		entryToCommandName[entryName] = name
	}
	sort.Strings(entries)
	selectedItem, err := prompter.PromptSelect("Select a command", entries)
	if err != nil {
		return err
	}
	cmd := cmdMap[entryToCommandName[selectedItem]]

	fmt.Printf("Name: %s\nDescription: %s\nPlatform: %s\nCommand:\n", cmd.Name, cmd.Description, cmd.Platform)
	color.Yellow(cmd.Command)
	fmt.Println()
	return nil
}

// ProcessAction takes the command file path, command map and an action to invoke the action.
func ProcessAction(cmdMap map[string]Command, action Action, fileOps CmdFileOperater) (map[string]Command, error) {
	switch action {
	case Create:
		newMap, err := CreateHandler(cmdMap)
		if err != nil {
			return nil, err
		}
		err = updateFile(newMap, fileOps)
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
		err = updateFile(newMap, fileOps)
		if err != nil {
			return nil, err
		}
		return newMap, nil
	case Remove:
		newMap, err := DeleteHandler(cmdMap)
		if err != nil {
			return nil, err
		}
		err = updateFile(newMap, fileOps)
		if err != nil {
			return nil, err
		}
		return newMap, nil
	case Import:
		f, err := prompter.PromptString("Enter the absolute file path for the commands you want to import")
		if err != nil {
			return nil, err
		}
		isMerge, err := prompter.PromptConfirm("Do you want to merge the new commands with existing commands?")
		if err != nil {
			return nil, err
		}
		cmdMap, err = ImportHandler(f, isMerge, cmdMap, fileOps)
		if err != nil {
			return nil, err
		}
		err = updateFile(cmdMap, fileOps)
		if err != nil {
			return nil, err
		}
		return cmdMap, nil
	default:
		return nil, errors.New("unrecognized action")
	}
}

// ImportHandler loads the commands from the given path, and either merges or
// replaces the command map with new commands.
func ImportHandler(fpath string, isMerge bool, cmdMap map[string]Command, fileOps CmdFileOperater) (map[string]Command, error) {
	newCmds, err := fileOps.LoadExternal(fpath)
	if err != nil {
		return nil, err
	}
	if isMerge {
		for _, cmd := range newCmds {
			cmdMap[cmd.Name] = cmd
		}
		return cmdMap, nil
	}
	// new map, instead of merge
	newMap := make(map[string]Command)
	for _, cmd := range newCmds {
		newMap[cmd.Name] = cmd
	}
	return newMap, nil
}

// GetSelectedAction returns the action the user chose from the prompt.
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

// CreateHandler handles the creation of a new command.
func CreateHandler(cmdMap map[string]Command) (map[string]Command, error) {
	cmd, err := createCommand()
	if err != nil {
		return nil, err
	}
	newMap := copyMap(cmdMap)
	newMap[cmd.Name] = cmd
	return newMap, nil
}

// LoadCommands reads existing commands from the commands file.
func LoadCommands(fileOps CmdFileOperater) (map[string]Command, error) {
	cmdMap := make(map[string]Command)
	commands, err := fileOps.Load()
	if err != nil {
		return nil, err
	}
	for _, command := range commands {
		cmdMap[command.Name] = command
	}
	return cmdMap, nil
}

// CreateDirIfNotExists creates the given directory if it doesn't exist.
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

// parseCommands is a helper function that uses the provider to get
// json bytes to unmarshall into commands.
func parseCommands(cmdProvider func() ([]byte, error)) ([]Command, error) {
	var commands []Command
	cmdJSON, err := cmdProvider()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(cmdJSON, &commands)
	if err != nil {
		return nil, err
	}
	return commands, nil
}
