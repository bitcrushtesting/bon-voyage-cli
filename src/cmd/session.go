// Copyright 2024 Bitcrush Testing

package cmd

import (
	"bon-voyage-cli/connection"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sessionCmd)
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Create, read, stop, delete and list log sessions",
	Long:  `Session log debug data from the given device into a database.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("Missing sub-command")
			os.Exit(1)
		}
		if err := connection.LoadToken(); err != nil {
			fmt.Println("Authorization token not available", err)
			os.Exit(1)
		}

		switch args[0] {
		case "list":
			sessions, err := connection.SessionList()
			if err != nil {
				fmt.Println("Could not get session list:", err)
				os.Exit(1)
			}
			fmt.Println("Session count:", len(sessions))

			for i, session := range sessions {
				fmt.Println("   -", i, "Session ID:", session.SessionId, "Device ID:", session.DeviceId, "Created:", session.Created)
			}

		case "create":
			if len(args) < 2 {
				fmt.Println("Expected device id after 'create'")
				os.Exit(1)
			}
			deviceId := args[1]
			sessionId, err := connection.SessionCreate(deviceId)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Started session for device '%s' with session id: '%s'\n", deviceId, sessionId)

		case "read":
			if len(args) < 2 {
				fmt.Println("Expected session id after 'read'")
				os.Exit(1)
			}
			sessionId := args[1]
			logs, err := connection.SessionRead(sessionId, nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("-----------------------")
			fmt.Printf("Logs: %d lines \n", len(logs))
			fmt.Println(logs)
			fmt.Println("-----------------------")

		case "stop":
			if len(args) < 2 {
				fmt.Println("Expected session id after 'stop'")
				os.Exit(1)
			}
			sessionId := args[1]
			err := connection.SessionUpdate(sessionId, "stop")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Session stopped successfully")

		case "delete":
			if len(args) < 2 {
				fmt.Println("Expected session id after 'delete'")
				os.Exit(1)
			}
			sessionId := args[1]
			if err := connection.SessionDelete(sessionId); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Session deleted successfully")

		default:
			fmt.Println("Unknown 'session' sub-command:", args[0])
			os.Exit(1)
		}
	},
}
