package display

import (
	"fmt"
	"os"
	"specgate/internal/validate"
	"text/tabwriter"
	"github.com/fatih/color"
)

var W *tabwriter.Writer = tabwriter.NewWriter(os.Stdout, 10, 0, 3, ' ', 0)

func PrintResults(file string, result *validate.CheckResult, strict bool) {

	PrintSummary(file, result, strict)

	if result.HasErrors() {
		printErrors(result, strict)
	}

	if !strict && result.HasWarnings() {
		printWarnings(result, strict)
	}

	W.Flush()

	if !strict && result.HasWarnings() {
		fmt.Println()
		fmt.Println("Run with --strict to treat warnings as errors.")
	} 

	if !result.HasErrors() && !result.HasWarnings() {
		fmt.Println("✅ Spec is ready.")
	}

	if result.HasErrors() {
		os.Exit(1)
	}
}

func printErrors(result *validate.CheckResult, strict bool) {

	error := color.RedString("error\t")

	if len(result.OperationSummaryViolations) > 0 {
		for _, item := range result.OperationSummaryViolations {
			fmt.Fprintf(W, "%sMissing operation summary\t%s\n", error, item)
		}
	}

	if len(result.SuccessResponseViolations) > 0 {
		for _, item := range result.SuccessResponseViolations {
			fmt.Fprintf(W, "%sMissing success responses (2xx)\t%s\n", error, item)
		}
	}

	if len(result.ErrorResponseViolations) > 0 {
		for _, item := range result.ErrorResponseViolations {
			fmt.Fprintf(W, "%sMissing error responses (4xx/5xx/default)\t%s\n", error, item)
		}
	}

	if len(result.SuccessResponseDescriptionViolations) > 0 {
		for _, item := range result.SuccessResponseDescriptionViolations {
			fmt.Fprintf(W, "%sMissing descriptions for success responses (2xx)\t%s\n", error, item)
		}
	}

	if len(result.ErrorResponseDescriptionViolations) > 0 {
		for _, item := range result.ErrorResponseDescriptionViolations {
			fmt.Fprintf(W, "%sMissing descriptions for error responses (4xx/5xx/default)\t%s\n", error, item)
		}
	}

	if result.MissingServers {
		fmt.Fprintf(W, "%sNo servers defined\t\n", error)
	}

	if len(result.ServerPlaceholderViolations) > 0 {
		fmt.Printf("Server URL(s) contain placeholders or non-production URLs:\n")
		fmt.Println()
		for _, item := range result.ServerPlaceholderViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}

	if len(result.ParamDescriptionViolations) > 0 {
		for _, item := range result.ParamDescriptionViolations {
			fmt.Fprintf(W, "%sMissing parameter description\t%s\n", error, item)
		}
	}

	if strict && hasWarnings(result) {
		printWarnings(result, strict)
	}
}

func printWarnings(result *validate.CheckResult, strict bool) {
	var severity string
	if strict {
		severity = color.RedString("error\t")
	} else {
		severity = color.YellowString("warning\t")
	}

	if len(result.OperationDescriptionViolations) > 0 {
		for _, item := range result.OperationDescriptionViolations {
			fmt.Fprintf(W, "%sMissing operation description\t%s\n", severity, item)
		}
	}
	if len(result.OperationTagViolations) > 0 {
		for _, item := range result.OperationTagViolations {
			fmt.Fprintf(W, "%sMissing operation tag\t%s\n", severity, item)
		}
	}
	if len(result.OperationIdViolations) > 0 {
		for _, item := range result.OperationIdViolations {
			fmt.Fprintf(W, "%sMissing operation id\t%s\n", severity, item)
		}
	}

}

func hasWarnings(result *validate.CheckResult) bool {
	return len(result.OperationDescriptionViolations) > 0 ||
		len(result.OperationTagViolations) > 0 ||
		len(result.OperationIdViolations) > 0
}

func PrintSummary(file string, result *validate.CheckResult, strict bool) {
	var ErrorCount int = validate.CountErrors(result)
	var WarningCount int = validate.CountWarnings(result)

	if strict && result.HasWarnings() {
		fmt.Printf("%s - %d errors\n", file, ErrorCount+WarningCount)
		fmt.Println()
	} else {
		fmt.Printf("%s - %d errors, %d warnings\n", file, ErrorCount, WarningCount)
		fmt.Println()
	}
}
