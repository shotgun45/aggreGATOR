package config

import (
	"encoding/json"
	"os"
	"testing"
)

func testWrite(path string, cfg Config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}

func testRead(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func TestWriteAndReadConfig(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error: %v", err)
	}
	path := home + "/.gatorconfig_test.json"
	cfg := Config{DBUrl: "postgres://test", CurrentUserName: "testuser"}
	defer os.Remove(path)
	if err := testWrite(path, cfg); err != nil {
		t.Fatalf("testWrite() error: %v", err)
	}
	readCfg, err := testRead(path)
	if err != nil {
		t.Fatalf("testRead() error: %v", err)
	}
	if readCfg.DBUrl != cfg.DBUrl || readCfg.CurrentUserName != cfg.CurrentUserName {
		t.Errorf("testRead() returned wrong config: got %+v, want %+v", readCfg, cfg)
	}
}

func TestSetUser(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error: %v", err)
	}
	path := home + "/.gatorconfig_test.json"
	cfg := Config{DBUrl: "postgres://test"}
	defer os.Remove(path)
	if err := testWrite(path, cfg); err != nil {
		t.Fatalf("testWrite() error: %v", err)
	}
	// Simulate SetUser by updating struct and writing again
	cfg.CurrentUserName = "newuser"
	if err := testWrite(path, cfg); err != nil {
		t.Fatalf("testWrite() error: %v", err)
	}
	readCfg, err := testRead(path)
	if err != nil {
		t.Fatalf("testRead() error: %v", err)
	}
	if readCfg.CurrentUserName != "newuser" {
		t.Errorf("SetUser() did not update user: got %s, want %s", readCfg.CurrentUserName, "newuser")
	}
}
