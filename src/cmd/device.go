// Copyright 2024 Bitcrush Testing

package cmd

import (
	"bon-voyage-cli/connection"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deviceCmd)
}

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Interact with devices connected to the server",
	Long:  `Login device description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Device command used")
		if len(args) < 1 {
			os.Exit(1)
		}
		if err := connection.LoadToken(); err != nil {
			fmt.Println("Token not available", err)
			os.Exit(1)
		}

		switch args[0] {
		case "list":
			devices, err := connection.DeviceList()
			if err != nil {
				fmt.Println("Could not get device list:", err)
				os.Exit(1)
			}
			fmt.Println("Device count:", len(devices))
			if len(devices) > 0 {
				fmt.Println(devices)
			}

		case "socket":
			if len(args) < 2 {
				fmt.Println("Expected device uuid after socket")
				os.Exit(1)
			}
			deviceId := args[1]
			if err := connection.DeviceSocket(deviceId); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		default:
			fmt.Println("Unknown device subcommand:", args[2])
			os.Exit(1)
		}
	},
}
