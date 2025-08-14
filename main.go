package main

import (
	cmds "aggreGATOR/internal/commands"
	"aggreGATOR/internal/config"
	"aggreGATOR/internal/database"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	state := &cmds.State{Db: dbQueries, Cfg: &cfg}
	commandSet := cmds.DefaultCommands()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Not enough arguments. Usage: gator <command> [args...]\n")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := cmds.Command{Name: cmdName, Args: cmdArgs}

	err = commandSet.Run(state, cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
