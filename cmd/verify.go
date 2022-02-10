/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify configs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("verify called")
	},
}

func init() {
	configCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVarP(&ChainName, "chain", "c", "", "Chain name")
}
