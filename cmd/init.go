/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/itsdeannat/specgate/internal/settings"

	"github.com/spf13/cobra"
)

var force bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a SpecGate config file",
	Long: `Create a default SpecGate configuration file in the current directory. 
	
This file defines validation rules and quality gates used by SpecGate commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat("./.specgate.yaml")

		if err == nil && !force {
			fmt.Println("A SpecGate config file already exists in this directory.\nUse --force to overwrite it.")
			return
		}

		settings.CreateConfig()

		if force {
			fmt.Println("Config file overwritten.")
		} else {
			fmt.Println("Config file created.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&force, "force", false, "Overwrite an existing config file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
