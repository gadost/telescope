/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate new <chain> config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")
	},
}

var ChainName string

func init() {
	configCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&ChainName, "chain", "r", "", "Chain name required")
	generateCmd.MarkFlagRequired("chain")
}
