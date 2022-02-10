/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/watcher"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify configs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		Verify()
	},
}

func init() {
	configCmd.AddCommand(verifyCmd)
	verifyCmd.PersistentFlags().StringVarP(&c, "conf", "c", conf.UserHome+"/.telescope/conf.d", "Configurations directory")
}

func Verify() {
	cfg, chains := conf.ConfLoad(c)
	for _, chainName := range chains {
		if cfg.Chain[chainName].Info.Github != "" {

			repoInfo := watcher.Parse(cfg.Chain[chainName].Info.Github)
			if repoInfo.Domain == "" || repoInfo.Owner == "" && repoInfo.RepoName == "" {
				fmt.Printf("Can't parse %s as github repository", cfg.Chain[chainName].Info.Github)
			}
		}

		for _, node := range cfg.Chain[chainName].Node {
			if node.Role != "validator" && node.Role != "sentry" {
				fmt.Printf("Can't parse node type %s , chain: %s", node.Role, chainName)
			}
			if node.RPC == "" {
				fmt.Printf("empty RPC for node , chain: %s", chainName)
			}
		}
	}

}
