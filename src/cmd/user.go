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
	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:   "user <sub-command>",
	Short: "user command",
	Long:  `user command long version`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("Sub-command missing")
			os.Exit(1)
		}

		switch args[0] {
		case "register":
			if len(args) < 2 {
				fmt.Println("User missing")
				os.Exit(1)
			}
			username := args[1]
			password, err := utils.ReadPassword()
			if err != nil {
				fmt.Printf("Failed to read password: %v", err)
				os.Exit(1)
			}
			if err := connection.Register(username, password); err != nil {
				fmt.Printf("Failed to register account: %v", err)
				os.Exit(1)
			}
			fmt.Println("Successfully registered account:", username)

		case "delete":
			if len(args) < 2 {
				fmt.Println("User missing")
				os.Exit(1)
			}
			d, err := utils.AskForConfirmation("Do you want to delete the account?")
			if err != nil {
				os.Exit(1)
			}
			if !d {
				os.Exit(0)
			}
			username := args[1]
			password, err := utils.ReadPassword()
			if err != nil {
				fmt.Printf("Failed to read password: %v", err)
				os.Exit(1)
			}
			if err := connection.DeleteAccount(username, password); err != nil {
				fmt.Printf("Failed to delete account: %v", err)
				os.Exit(1)
			}
			fmt.Println("Successfully deleted account:", username)

		case "change-password":
			fmt.Println("User delete command used")
			//password, err := utils.ReadPassword()

		case "change-username":
			username := args[2]

			fmt.Println("Username:", username)
			fmt.Print("Enter password: ")

			// password, err := utils.ReadPassword()
			// if err != nil {
			// 	fmt.Printf("Failed to read password: %v", err)
			// 	os.Exit(1)
			// }
		}
	},
}
