// Copyright 2024 Bitcrush Testing

package cmd

import (
	"bon-voyage-cli/connection"
	"bon-voyage-cli/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the 'bon-voyage' server",
	Long:  `All software has versions. This is bon-voyage-cli's`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("User missing")
			os.Exit(1)
		}
		username := args[0]
		fmt.Println("Username:", username)

		password, err := utils.ReadPassword()
		if err != nil {
			fmt.Printf("Failed to read password: %v", err)
			os.Exit(1)
		}
		if err := connection.Login(username, password); err != nil {
			fmt.Println("Could not log in", err)
			os.Exit(1)
		}
		fmt.Println("Successfully logged in account:", username)
	},
}
