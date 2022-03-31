package cmdref_test

import (
	"cmdref"
	"testing"
)

func TestParseCommands(t *testing.T) {
	provider := func() ([]byte, error) {
		s := "[{\"name\": \"test command\", \"command\": \"go test\", " +
			"\"platform\": \"go cli\", \"description\": \"run tests\"}]"
		return []byte(s), nil
	}
	cmdMap, _ := cmdref.ParseCommands(provider)
	expected := map[string]cmdref.Command{
		"test command": {Name: "test command", Command: "go test", Platform: "go cli", Description: "run tests"},
	}
	command, prs := cmdMap["test command"]
	if !prs {
		t.Errorf("'test command' not present in parsed map")
	}
	if command != expected["test command"] {
		t.Errorf("got: %#v, want: %#v", command, expected["test command"])
	}

}
