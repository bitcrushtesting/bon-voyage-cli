// Copyright 2024 Bitcrush Testing

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "0.0.0"
	Commit  string
	AppName = "bon-voyage-cli"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of " + AppName,
	Long:  `All software has versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(AppName)
		fmt.Println("Version", Version+"-git"+Commit)
		fmt.Println("Copyright Bitcrush Testing (C) 2024")
		os.Exit(0)
	},
}
