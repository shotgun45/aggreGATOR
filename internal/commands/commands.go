package commands

import (
	"aggreGATOR/internal/config"
	"aggreGATOR/internal/database"
	"aggreGATOR/internal/rssfeed"
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
	cmds.Register("agg", HandlerAgg)
	cmds.Register("addfeed", HandlerAddFeed)
	cmds.Register("feeds", HandlerFeeds)
	cmds.Register("follow", HandlerFollow)
	cmds.Register("following", HandlerFollowing)
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

func HandlerAgg(s *State, cmd Command) error {
	ctx := context.Background()
	feed, err := rssfeed.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("addfeed requires name and url arguments")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	userName := s.Cfg.CurrentUserName
	user, err := s.Db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("could not find current user: %v", err)
	}
	params := database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: user.ID,
	}
	feed, err := s.Db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}
	fmt.Printf("Feed created: ID=%v Name=%s Url=%s UserID=%v\n", feed.ID, feed.Name, feed.Url, feed.UserID)

	id := uuid.New()
	now := time.Now()
	ff, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %v", err)
	}
	fmt.Printf("You are now following feed '%s' as user '%s'\n", ff.FeedName, ff.UserName)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeedsWithUser(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("- %s\n  url: %s\n  created by: %s\n", feed.Name, feed.Url, feed.UserName)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("follow requires a feed url argument")
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not find feed with url %s: %v", url, err)
	}
	userName := s.Cfg.CurrentUserName
	user, err := s.Db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("could not find current user: %v", err)
	}
	id := uuid.New()
	now := time.Now()
	ff, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %v", err)
	}
	fmt.Printf("Followed feed '%s' as user '%s'\n", ff.FeedName, ff.UserName)
	return nil
}

func HandlerFollowing(s *State, cmd Command) error {
	userName := s.Cfg.CurrentUserName
	user, err := s.Db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("could not find current user: %v", err)
	}
	follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %v", err)
	}
	if len(follows) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}
	fmt.Println("Feeds you are following:")
	for _, f := range follows {
		fmt.Printf("* %s\n", f.FeedName)
	}
	return nil
}
