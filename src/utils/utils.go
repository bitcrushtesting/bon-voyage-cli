// Copyright 2024 Bitcrush Testing

package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadPassword() (string, error) {
	fmt.Println("Enter password:")
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}

func AskForConfirmation(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("error reading input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Println("Please enter 'y' or 'n'.")
		}
	}
}
