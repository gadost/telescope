/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "verify or generate configs",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
