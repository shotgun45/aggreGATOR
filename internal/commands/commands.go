package commands

import (
	"aggreGATOR/internal/config"
	"aggreGATOR/internal/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}

func DefaultCommands() *Commands {
	cmds := &Commands{Handlers: make(map[string]func(*State, Command) error)}
	cmds.Register("login", HandlerLogin)
	cmds.Register("register", HandlerRegister)
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerUsers)
	return cmds
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login requires a username argument")
	}
	username := cmd.Args[0]
	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user '%s' does not exist", username)
	}
	err = s.Cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User set to '%s'\n", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("register requires a username argument")
	}
	name := cmd.Args[0]
	// Check if user exists
	_, err := s.Db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("user '%s' already exists", name)
	}
	id := uuid.New()
	now := time.Now()
	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}
	user, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	err = s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("failed to set current user: %v", err)
	}
	fmt.Printf("User '%s' created!\n", name)
	fmt.Printf("User data: %+v\n", user)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("reset failed: %v", err)
	}
	fmt.Println("All users deleted successfully.")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}
	current := s.Cfg.CurrentUserName
	for _, u := range users {
		if u.Name == current {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}
	return nil
}
