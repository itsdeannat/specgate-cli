package display

import (
	"fmt"
	"specgate/internal/validate"
)


func PrintResults(result *validate.CheckResult, strict bool) {

	hasErrors := len(result.OperationSummaryViolations) > 0 || 
		len(result.SuccessResponseViolations) > 0 ||
		len(result.ErrorResponseViolations) > 0 ||
		len(result.SuccessResponseDescriptionViolations) > 0 ||
		len(result.ErrorResponseDescriptionViolations) > 0 ||
		result.MissingServers == true ||
		len(result.ServerPlaceholderViolations) > 0 
		// (strict && len(result.OperationDescriptionViolations) > 0)

	if hasErrors { // only prints errors if hasErrors is true
		fmt.Println("ERRORS")
		fmt.Println("------")
		fmt.Println()

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

		if result.MissingServers == true {
			fmt.Printf("No servers defined - add a server URL to your spec\n")
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

		// if strict && len(result.OperationDescriptionViolations) > 0 {
		// 	fmt.Printf("Missing operation descriptions for %d operation(s):\n", len(result.OperationDescriptionViolations))
		// 	fmt.Println()
		// 	for _, item := range result.OperationDescriptionViolations {
		// 		fmt.Println("-", item)
		// 	}
		// 	fmt.Println()
		// }
		

	} else {
		fmt.Println("✅ No errors found.")
		fmt.Println()
	}

	if !strict && len(result.OperationDescriptionViolations) > 0 {
		fmt.Println("WARNINGS")
		fmt.Println("--------")
		fmt.Println()
		fmt.Printf("Missing operation descriptions for %d operation(s):\n", len(result.OperationDescriptionViolations))
		fmt.Println()
		for _, item := range result.OperationDescriptionViolations {
			fmt.Println("-", item)
		}
		fmt.Println()
	}


	if !hasErrors && len(result.OperationDescriptionViolations) == 0 && len(result.SuccessResponseDescriptionViolations) == 0 {
		fmt.Println("✅ Spec is ready.")
	}
}