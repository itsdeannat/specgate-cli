/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rulesCmd represents the rules command
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Display the rules used to validate OpenAPI specs",
	Long: `rules displays the error and warning rules that SpecGate enforces 
when validating an OpenAPI specification.

Error rules will fail the check and return a non-zero exit code.
Warning rules are reported but do not fail the check unless --strict is enabled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ERRORS")
		fmt.Println("------")
		fmt.Println("- Operations must have a summary")
		fmt.Println("- Operations must include at least one 2xx response")
		fmt.Println("- Operations must include at least one error response (4xx, 5xx, default)")
		fmt.Println("- Success responses must include descriptions")
		fmt.Println("- Error responses must include descriptions")
		fmt.Println("- Server URLs must be present and cannot contain placeholders like 'example.com' or 'localhost'")
		fmt.Println("- A servers object must be included")
		fmt.Println()

		fmt.Println("WARNINGS")
		fmt.Println("--------")
		fmt.Println("- Operations should include a description")
		fmt.Println()

		fmt.Println("STRICT MODE") 
		fmt.Println("-----------")
		fmt.Println("Run with --strict to treat warning rules as errors.")
	},
}

func init() {
	rootCmd.AddCommand(rulesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rulesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rulesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
