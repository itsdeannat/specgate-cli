/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itsdeannat/specgate/internal/llm"
	"github.com/itsdeannat/specgate/internal/report"
	"github.com/itsdeannat/specgate/internal/validate"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

// adviseCmd represents the advise command
var adviseCmd = &cobra.Command{
	Use:   "advise",
	Short: "Get AI-powered advice for improving your spec",
	Long: `advise uses an LLM to suggest summaries and descriptions for operations that are missing them.
	
Suggestions are advisory and human-readable, are not applied automatically,
and do not affect pass/fail status. Review and apply changes manually.`,
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		specBytes, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "specgate: failed to read spec %q: %v\n", file, err)
			os.Exit(2)
		}

		loader := openapi3.NewLoader()

		doc, err := loader.LoadFromFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "specgate: failed to load spec %q: %v\n", file, err)
			os.Exit(2)
		}

		if err := doc.Validate(loader.Context); err != nil {
			fmt.Fprintf(os.Stderr, "specgate: spec validation failed: %v\n", err)
			os.Exit(2)
		}

		result := &validate.CheckResult{}

		validate.CheckPaths(doc, result)

		jsonOutput := report.ToJsonFormat(result, strict)
		jsonBytes, err := json.MarshalIndent(jsonOutput, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "specgate: failed to marshal JSON output: %v\n", err)
			os.Exit(2)
		}

		suggestions, err := llm.SuggestFromReport(jsonBytes, specBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "specgate: failed to generate suggestions: %v\n", err)
			os.Exit(2)
		}

		fmt.Println("Suggested documentation")
		fmt.Println("-----------------------")
		fmt.Println("Suggestions are advisory. Review before applying.")
		fmt.Println(suggestions)

		if result.HasErrors() {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(adviseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adviseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adviseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
