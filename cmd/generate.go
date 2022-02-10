/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gadost/telescope/conf"
	"github.com/spf13/cobra"
)

var ChainConfBasic = ``
var ChainName string
var Print bool
var err error
var pe int64
var bm int64
var pc int64
var Gh string
var nm bool
var skel = `
[info]
voting_power_changes = 100
blocks_missed_in_a_row = 10
peers_count = 10
github = "https://github.com/REPOOWNER/REPONAME"
[[node]]
role = "validator"
rpc = "http://1.2.3.4:26657"
monitoring_enabled = true
[[node]]
role = "sentry"
rpc = "http://2.3.4.5:26657"
network_monitoring_enabled = true
monitoring_enabled = true
`

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate new <chain> config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if Print {
			fmt.Println(conf.UserHome + "/.telescope/conf.d/" + ChainName + ".toml")
			fmt.Print("\n")
			fmt.Print(skel)
			os.Exit(0)
		}
		BootstrapChainConf()
	},
}

func init() {
	configCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&ChainName, "chain", "c", "", "Chain name required")
	generateCmd.Flags().BoolVarP(&Print, "print-skel", "p", false, "Print skel")
	generateCmd.MarkFlagRequired("chain")
}

func BootstrapChainConf() {
	if _, err := os.Stat(conf.UserHome + "/.telescope/conf.d/" + ChainName + ".toml"); !os.IsNotExist(err) {
		fmt.Printf("config file for chain %s exits", ChainName)
		os.Exit(1)
	}
	for {
		fmt.Fprint(os.Stderr, "Add new chain "+ChainName+" ?(Y/n): ")
		s, _ = r.ReadString('\n')
		s = strings.TrimSuffix(s, "\n")
		if s == "Y" || s == "y" || s == "" {
			break
		} else if s == "N" || s == "n" {
			os.Exit(0)
		}
	}

	for {
		fmt.Fprint(os.Stderr, "Voting power changes step (default 10): ")
		if p, _ := r.ReadString('\n'); p != "" && p != "\n" {
			p = strings.TrimSuffix(p, "\n")
			pe, err = strconv.ParseInt(p, 10, 64)
			if err == nil {
				if pe >= 1 {
					break
				} else {
					fmt.Println("must be >=1")
				}
			} else {
				fmt.Print(err)
			}
		} else if p == "\n" {
			pe = 10
			break
		}
	}
	for {
		fmt.Fprint(os.Stderr, "Blocks missed in a row before alert (default 5): ")
		if p, _ := r.ReadString('\n'); p != "" && p != "\n" {
			p = strings.TrimSuffix(p, "\n")
			bm, err = strconv.ParseInt(p, 10, 64)
			if err == nil {
				if bm >= 1 {
					break
				} else {
					fmt.Println("must be >=1")
				}
			} else {
				fmt.Print(err)
			}
		} else if p == "\n" {
			bm = 10
			break
		}
	}

	for {
		fmt.Fprint(os.Stderr, "Minimum peers connected alert (default 5): ")
		if p, _ := r.ReadString('\n'); p != "" && p != "\n" {
			p = strings.TrimSuffix(p, "\n")
			pc, err = strconv.ParseInt(p, 10, 64)
			if err == nil {
				if pc >= 1 {
					break
				} else {
					fmt.Println("must be >=1")
				}
			} else {
				fmt.Print(err)
			}
		} else if p == "\n" {
			pc = 10
			break
		}
	}

	for {
		fmt.Fprint(os.Stderr, "Github repository Url ( f.e. https://github.com/REPOOWNER/REPONAME): ")
		Gh, _ = r.ReadString('\n')
		if Gh != "" {
			Gh = strings.TrimSuffix(Gh, "\n")
			break
		}
	}

	ChainConfBasic += fmt.Sprintf(`
[info]
voting_power_changes = %d
blocks_missed_in_a_row = %d
peers_count = %d
github = "%s"
`, pe, bm, pc, Gh)

	for {
		fmt.Print("Add new node\n")
		fmt.Fprint(os.Stderr, "Node type (validator/sentry): ")
		t, _ := r.ReadString('\n')
		if t != "" {
			t = strings.TrimSuffix(t, "\n")
		}

		fmt.Fprint(os.Stderr, "Node RPC (f.e. http://1.2.3.4:26657): ")
		rpc, _ := r.ReadString('\n')
		if rpc != "" {
			rpc = strings.TrimSuffix(rpc, "\n")
		}

		fmt.Fprint(os.Stderr, "Enable network monitoring via this node?(Y/n): ")
		n, _ := r.ReadString('\n')
		n = strings.TrimSuffix(n, "\n")
		if n == "Y" || n == "y" || n == "yes" {
			nm = true
		} else {
			nm = false
		}
		ChainConfBasic += fmt.Sprintf(`
[[node]]
role = "%s"
rpc = "%s"
network_monitoring_enabled = %t
monitoring_enabled = true
`, t, rpc, nm)
		fmt.Fprint(os.Stderr, "Add  one more node?(Y/n): ")
		next, _ := r.ReadString('\n')
		next = strings.TrimSuffix(next, "\n")
		if next == "n" || next == "N" || next == "No" || next == "no" {
			Write(ChainName+".toml", ChainConfBasic)
			fmt.Print("Done.")
			os.Exit(0)
		}
	}
}
