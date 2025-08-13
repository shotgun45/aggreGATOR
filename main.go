package main

import (
	"aggreGATOR/internal/config"
	"fmt"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	err = cfg.SetUser("shotgun45")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting user: %v\n", err)
		os.Exit(1)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config: %+v\n", cfg)
}
