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

	fmt.Printf("%s - %d errors, %d warnings\n", file, validate.CountErrors(result), validate.CountWarnings(result))
	fmt.Println()

	if result.HasErrors() {
		printErrors(result, strict)
	} else {
		fmt.Println("✅ No errors found.")
		fmt.Println()
	}

	if !strict && result.HasWarnings() {
		printWarnings(result)
		W.Flush()
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
		for _, item := range result.OperationSummaryViolations {
			fmt.Fprintf(W, "%sMissing operation summary\t%s\n", color.RedString("error\t"), item)
		}
	}

	if len(result.SuccessResponseViolations) > 0 {
		for _, item := range result.SuccessResponseViolations {
			fmt.Fprintf(W, "%sMissing success responses (2xx)\t%s\n", color.RedString("error\t"), item)
		}
	}

	if len(result.ErrorResponseViolations) > 0 {
		for _, item := range result.ErrorResponseViolations {
			fmt.Fprintf(W, "%sMissing error responses (4xx/5xx/default)\t%s\n", color.RedString("error\t"), item)
		}
	}

	if len(result.SuccessResponseDescriptionViolations) > 0 {
		for _, item := range result.SuccessResponseDescriptionViolations {
			fmt.Fprintf(W, "%sMissing descriptions for success responses (2xx)\t%s\n", color.RedString("error\t"), item)
		}
	}

	if len(result.ErrorResponseDescriptionViolations) > 0 {
		for _, item := range result.ErrorResponseDescriptionViolations {
			fmt.Fprintf(W, "%sMissing descriptions for error responses (4xx/5xx/default)\t%s\n", color.RedString("error\t"), item)
		}
	}

	if result.MissingServers {
		fmt.Fprintf(W, "%sNo servers defined\t\n", color.RedString("error\t"))
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
			fmt.Fprintf(W, "%sMissing parameter description\t%s\n", color.RedString("error\t"), item)
		}
	}

	if strict && hasWarnings(result) {
		printWarnings(result)
	}
}

func printWarnings(result *validate.CheckResult) {
	if len(result.OperationDescriptionViolations) > 0 {
		for _, item := range result.OperationDescriptionViolations {
			fmt.Fprintf(W, "%sMissing operation description\t%s\n", color.YellowString("warning\t"), item)
		}
	}
	if len(result.OperationTagViolations) > 0 {
		for _, item := range result.OperationTagViolations {
			fmt.Fprintf(W, "%sMissing operation tag\t%s\n", color.YellowString("warning\t"), item)
		}
	}
	if len(result.OperationIdViolations) > 0 {
		for _, item := range result.OperationIdViolations {
			fmt.Fprintf(W, "%sMissing operation id\t%s\n", color.YellowString("warning\t"), item)
		}
	}
	
}

func hasWarnings(result *validate.CheckResult) bool {
	return len(result.OperationDescriptionViolations) > 0 ||
	len(result.OperationTagViolations) > 0 || 
	len(result.OperationIdViolations) > 0 
}
