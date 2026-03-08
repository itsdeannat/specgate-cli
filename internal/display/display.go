package display

import (
	"fmt"
	"os"
	"specgate/internal/validate"
)

func PrintResults(result *validate.CheckResult, strict bool) {

	if result.HasErrors() {
		fmt.Println("ERRORS")
		fmt.Println("------")
		fmt.Println()
		printErrors(result, strict)
		fmt.Println()
	} else {
		fmt.Println("✅ No errors found.")
		fmt.Println()
	}

		if !strict && result.HasWarnings() {
		fmt.Println("WARNINGS")
		fmt.Println("--------")
		fmt.Println()
		printWarnings(result)
		fmt.Println()
		fmt.Println("Run with --strict to treat warnings as errors.")
		fmt.Println()
	}

	if !result.HasErrors() && !result.HasWarnings() {
		fmt.Println("✅ Spec is ready.")
	}

	if result.HasErrors() {
		os.Exit(1)
	}
}

func printErrors(result *validate.CheckResult, strict bool) {

	if len(result.OperationSummaryViolations) > 0 {
		fmt.Printf("Missing operation summaries for %d operation(s):\n", len(result.OperationSummaryViolations))
		fmt.Println()
		for _, item := range result.OperationSummaryViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}

	if len(result.SuccessResponseViolations) > 0 {
		fmt.Printf("Missing success responses (2xx) for %d operation(s):\n", len(result.SuccessResponseViolations))
		fmt.Println()
		for _, item := range result.SuccessResponseViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}

	if len(result.ErrorResponseViolations) > 0 {
		fmt.Printf("Missing error responses (4xx/5xx/default) for %d operation(s):\n", len(result.ErrorResponseViolations))
		fmt.Println()
		for _, item := range result.ErrorResponseViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}

	if len(result.SuccessResponseDescriptionViolations) > 0 {
		fmt.Printf("Missing descriptions for success responses (2xx) in %d operation(s):\n", len(result.SuccessResponseDescriptionViolations))
		fmt.Println()
		for _, item := range result.SuccessResponseDescriptionViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}

	if len(result.ErrorResponseDescriptionViolations) > 0 {
		fmt.Printf("Missing descriptions for error responses (4xx/5xx/default) in %d operation(s):\n", len(result.ErrorResponseDescriptionViolations))
		fmt.Println()
		for _, item := range result.ErrorResponseDescriptionViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}

	if result.MissingServers {
		fmt.Println("No servers defined: ")
		fmt.Println()
		fmt.Println("- Add a server URL and description to your spec.")
		fmt.Println()
	}

	if len(result.ServerPlaceholderViolations) > 0 {
		fmt.Printf("Server URL(s) contain placeholders:\n")
		fmt.Println()
		for _, item := range result.ServerPlaceholderViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}

	if strict && hasWarnings(result) {
		printWarnings(result)
	}
}

func printWarnings(result *validate.CheckResult) {
	if len(result.OperationDescriptionViolations) > 0 {
		fmt.Printf("Missing operation descriptions for %d operation(s):\n", len(result.OperationDescriptionViolations))
		fmt.Println()
		for _, item := range result.OperationDescriptionViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}
}

func hasWarnings(result *validate.CheckResult) bool {
	return len(result.OperationDescriptionViolations) > 0
}