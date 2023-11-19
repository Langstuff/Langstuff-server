package extensions

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	node := parseCommand("!learn #a #b")

	if !contains(node.Tags, "a") {
		t.Errorf("Expected 'a' to be in tags")
	}

	if !contains(node.Tags, "b") {
		t.Errorf("Expected 'b' to be in tags")
	}

	if len(node.Tags) != 2 {
		t.Errorf("Return value shoud have length of 2, but has %d", len(node.Tags))
	}

	node = parseCommand("!learn   #hello   #world   ")
	if !contains(node.Tags, "hello") {
		t.Errorf("Expected 'hello' to be in tags")
	}

	if !contains(node.Tags, "world") {
		t.Errorf("Expected 'world' to be in tags")
	}
}

func contains(slice []string, str string) bool {
    for _, s := range slice {
        if s == str {
            return true
        }
    }
    return false
}
