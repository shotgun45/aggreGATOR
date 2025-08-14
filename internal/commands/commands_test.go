package commands

import (
	"testing"
)

func TestRegisterAndRunCommand(t *testing.T) {
	cmds := &Commands{Handlers: make(map[string]func(*State, Command) error)}
	called := false
	handler := func(s *State, c Command) error {
		called = true
		if c.Name != "test" {
			t.Errorf("expected command name 'test', got '%s'", c.Name)
		}
		return nil
	}
	cmds.Register("test", handler)
	state := &State{}
	cmd := Command{Name: "test", Args: []string{}}
	if err := cmds.Run(state, cmd); err != nil {
		t.Errorf("Run returned error: %v", err)
	}
	if !called {
		t.Error("handler was not called")
	}
}

func TestRunUnknownCommand(t *testing.T) {
	cmds := &Commands{Handlers: make(map[string]func(*State, Command) error)}
	state := &State{}
	cmd := Command{Name: "unknown", Args: []string{}}
	if err := cmds.Run(state, cmd); err == nil {
		t.Error("expected error for unknown command, got nil")
	}
}
