/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"specgate/internal/display"
	"specgate/internal/report"
	"specgate/internal/settings"
	"specgate/internal/validate"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var strict bool
var outputFormat string

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:     "check <spec>",
	Example: `specgate check oas.json`,
	Short:   "Check an OpenAPI spec for readiness",
	Long: `check evaluates an OpenAPI specification against readiness rules.

If errors are detected, the command exits with a non-zero
status code, allowing it to be used as a quality gate in CI.`,
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

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

		_, configErr := os.Stat("./.specgate.yaml") // check if config exists

		if configErr == nil {
			fmt.Println("Loaded config from .specgate.yaml")
			fmt.Println()
		} else {
			fmt.Println("No SpecGate config found in project root. Created .specgate.yaml")
			settings.CreateConfig()
			fmt.Println("Edit the config, then rerun:\n specgate check")
			os.Exit(3)
		}

		configFile, err := os.ReadFile("./.specgate.yaml")
		var config settings.SpecGateConfig
		err = yaml.Unmarshal(configFile, &config)

		result := &validate.CheckResult{}

		for path, pathItem := range doc.Paths.Map() {
			if pathItem == nil {
				continue
			}
			if pathItem.Get != nil {
				validate.CheckOperation(pathItem.Get, fmt.Sprintf("GET %s", path), result)
			}
			if pathItem.Post != nil {
				validate.CheckOperation(pathItem.Post, fmt.Sprintf("POST %s", path), result)
			}
			if pathItem.Put != nil {
				validate.CheckOperation(pathItem.Put, fmt.Sprintf("PUT %s", path), result)
			}
			if pathItem.Patch != nil {
				validate.CheckOperation(pathItem.Patch, fmt.Sprintf("PATCH %s", path), result)
			}
			if pathItem.Delete != nil {
				validate.CheckOperation(pathItem.Delete, fmt.Sprintf("DELETE %s", path), result)
			}
		}

		if len(doc.Servers) == 0 {
			result.MissingServers = true
		}

		for _, server := range doc.Servers {
			validate.CheckServer(server, result, config.ServerBlockList)
		}

		if strict && outputFormat == "" {
			fmt.Println("STRICT MODE ENABLED")
			fmt.Println()
		}

		if outputFormat != "" && outputFormat != "json" {
			fmt.Fprintf(os.Stderr, "specgate: unsupported format %q\n", outputFormat)
			os.Exit(2)
		}

		if outputFormat == "json" {
			jsonOutput := report.ToJsonFormat(result, strict)
			jsonBytes, err := json.MarshalIndent(jsonOutput, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "specgate: failed to marshal JSON output: %v\n", err)
				os.Exit(2)
			}
			fmt.Println(string(jsonBytes))
		} else {
			display.PrintResults(file, result, strict)
		}

		if result.HasErrors() || (strict && result.HasWarnings()) {
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVar(&strict, "strict", false, "Treat warnings as errors")
	checkCmd.Flags().StringVar(&outputFormat, "format", "", "Output results as json")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
