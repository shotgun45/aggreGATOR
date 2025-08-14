package commands

import (
	"errors"
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

func TestHandlerReturnsError(t *testing.T) {
	cmds := &Commands{Handlers: make(map[string]func(*State, Command) error)}
	handler := func(s *State, c Command) error {
		return errors.New("handler error")
	}
	cmds.Register("fail", handler)
	state := &State{}
	cmd := Command{Name: "fail", Args: []string{}}
	if err := cmds.Run(state, cmd); err == nil {
		t.Error("expected error from handler, got nil")
	}
}

func TestMultipleCommandRegistration(t *testing.T) {
	cmds := &Commands{Handlers: make(map[string]func(*State, Command) error)}
	calls := make(map[string]bool)
	cmds.Register("one", func(s *State, c Command) error { calls["one"] = true; return nil })
	cmds.Register("two", func(s *State, c Command) error { calls["two"] = true; return nil })
	state := &State{}
	cmd1 := Command{Name: "one", Args: []string{}}
	cmd2 := Command{Name: "two", Args: []string{}}
	cmds.Run(state, cmd1)
	cmds.Run(state, cmd2)
	if !calls["one"] || !calls["two"] {
		t.Error("not all handlers were called")
	}
}
