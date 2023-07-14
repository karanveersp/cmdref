package cmdref_test

import (
	"path/filepath"
	"testing"

	cmdref "github.com/karanveersp/cmdref/internal"

	"github.com/stretchr/testify/mock"
)

// MockOperations implements the CommandsFileOperator for unit tests.
type MockOperations struct {
	mock.Mock
}

// Load parses the commands file and returns the list of commands.
func (m *MockOperations) Load() ([]cmdref.Command, error) {
	args := m.Called()
	return args.Get(0).([]cmdref.Command), args.Error(1)
}

func (m *MockOperations) LoadExternal(fpath string) ([]cmdref.Command, error) {
	args := m.Called(fpath)
	return args.Get(0).([]cmdref.Command), args.Error(1)
}

// Save writes the list of commands to the commands file.
func (m *MockOperations) Save(commands []cmdref.Command) error {
	args := m.Called(commands)
	return args.Error(0)
}

// GetFilePath returns the absolute path to the commands file.
func (m *MockOperations) GetFilePath() string {
	cmdStoreDir := "testing"
	return filepath.Join(cmdStoreDir, cmdref.CmdFileName)
}

func TestLoadCommands(t *testing.T) {
	// arrange
	mockCommand := cmdref.Command{Name: "test command", Command: "go test", Platform: "go cli", Description: "run tests"}

	mockFileOps := new(MockOperations)
	mockFileOps.On("Load").Return([]cmdref.Command{mockCommand}, nil)

	// act
	cmdMap, err := cmdref.LoadCommands(mockFileOps)

	// assert
	if err != nil {
		t.Errorf("error expected to be nil!")
	}
	expected := map[string]cmdref.Command{
		"test command": mockCommand,
	}
	command, present := cmdMap["test command"]
	if !present {
		t.Errorf("'test command' not present in parsed map")
	}
	if command != expected["test command"] {
		t.Errorf("got: %#v, want: %#v", command, expected["test command"])
	}
}

func TestImportHandlerMerge(t *testing.T) {
	// arrange
	command := cmdref.Command{Name: "test command", Command: "go test", Platform: "go cli", Description: "run tests"}
	cmdMap := make(map[string]cmdref.Command)
	cmdMap[command.Name] = command

	newCommand := cmdref.Command{Name: "cmd2", Command: "go run main.go", Platform: "go cli", Description: "run main.go"}

	mockFilePath := "path/to/import.json"

	mockFileOps := new(MockOperations)
	mockFileOps.On("LoadExternal", mockFilePath).Return([]cmdref.Command{newCommand}, nil)

	// act
	newCmdMap, err := cmdref.ImportHandler(mockFilePath, true, cmdMap, mockFileOps)

	// assert
	if err != nil {
		t.Errorf("error expected to be nil!")
	}
	expected := map[string]cmdref.Command{
		command.Name:    command,
		newCommand.Name: newCommand,
	}
	if !mapsAreEqual(expected, newCmdMap) {
		t.Errorf("Maps are not equal")
	}

}

func TestImportHandlerReplace(t *testing.T) {
	// arrange
	command := cmdref.Command{Name: "test command", Command: "go test", Platform: "go cli", Description: "run tests"}
	cmdMap := make(map[string]cmdref.Command)
	cmdMap[command.Name] = command

	newCommand := cmdref.Command{Name: "cmd2", Command: "go run main.go", Platform: "go cli", Description: "run main.go"}

	mockFilePath := "path/to/import.json"

	mockFileOps := new(MockOperations)
	mockFileOps.On("LoadExternal", mockFilePath).Return([]cmdref.Command{newCommand}, nil)

	// act
	newCmdMap, err := cmdref.ImportHandler(mockFilePath, false, cmdMap, mockFileOps)

	// assert
	if err != nil {
		t.Errorf("error expected to be nil!")
	}
	expected := map[string]cmdref.Command{
		newCommand.Name: newCommand,
	}
	if !mapsAreEqual(expected, newCmdMap) {
		t.Errorf("Maps are not equal")
	}

}

func mapsAreEqual(expected, actual map[string]cmdref.Command) bool {
	// Compare the lengths of the maps
	if len(expected) != len(actual) {
		return false
	}

	// Compare the key-value pairs of the maps
	for key, expectedValue := range expected {
		actualValue, ok := actual[key]
		if !ok || actualValue != expectedValue {
			return false
		}
	}

	return true
}
