// Copyright 2024 Bitcrush Testing

package connection

import (
	"fmt"
	"os"
	"path/filepath"
)

const tokenFileName = "bon_voyage_cli.token"

var Token string

func tokenFilePath() string {
	confDir, err := os.UserConfigDir()
	if err != nil {
		confDir = "./"
	}
	return filepath.Join(confDir, tokenFileName)
}

func SaveToken() error {
	fmt.Println("Saving token to", tokenFilePath())
	return os.WriteFile(tokenFilePath(), []byte(Token), 0600)
}

func DeleteToken() error {
	if _, err := os.Stat(tokenFilePath()); err != nil {
		return err
	}
	fmt.Println("Deleting token from", tokenFilePath())
	return os.Remove(tokenFilePath())
}

func LoadToken() error {
	t, err := os.ReadFile(tokenFilePath())
	if err != nil {
		return fmt.Errorf("no token found: %v", err)
	}
	token = string(t)
	return nil
}
