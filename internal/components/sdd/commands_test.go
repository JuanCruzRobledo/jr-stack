package sdd

import "testing"

func TestOpenCodeCommandsIncludesCoreWorkflow(t *testing.T) {
	commands := OpenCodeCommands()
	if len(commands) < 4 {
		t.Fatalf("OpenCodeCommands() length = %d, want at least 4", len(commands))
	}

	// OPSX: first command is now opsx-explore (replaced sdd-init).
	if commands[0].Name != "opsx-explore" {
		t.Fatalf("first command = %q, want %q", commands[0].Name, "opsx-explore")
	}
}
