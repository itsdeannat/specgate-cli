package report

import (
	"github.com/itsdeannat/specgate/internal/validate"
)

type JsonFormat struct {
    Ready bool `json:"ready"`
	Strict bool `json:"strict"`
	Errors map[string][]string `json:"errors,omitempty"`
	Warnings map[string][]string `json:"warnings,omitempty"`
}

func ToJsonFormat(result *validate.CheckResult, strict bool) *JsonFormat {
	jsonResult := &JsonFormat{
		
		Ready: !result.HasErrors() && !result.HasWarnings(),
		Strict: strict,
		Errors: make(map[string][]string),
		Warnings: make(map[string][]string),
	}

	if len(result.OperationSummaryViolations) > 0 {
		jsonResult.Errors["missing_operation_summaries"] = result.OperationSummaryViolations
	}
	if len(result.SuccessResponseViolations) > 0 {
		jsonResult.Errors["missing_success_responses"] = result.SuccessResponseViolations
	}
	if len(result.ErrorResponseViolations) > 0 {
		jsonResult.Errors["missing_error_responses"] = result.ErrorResponseViolations
	}
	if len(result.SuccessResponseDescriptionViolations) > 0 {
		jsonResult.Errors["missing_success_response_descriptions"] = result.SuccessResponseDescriptionViolations
	}
	if len(result.ErrorResponseDescriptionViolations) > 0 {
		jsonResult.Errors["missing_error_response_descriptions"] = result.ErrorResponseDescriptionViolations
	}
	if result.MissingServers {
		jsonResult.Errors["missing_servers"] = []string{"No servers defined - add a server URL to your spec"}
	} 
	if len(result.ServerPlaceholderViolations) > 0 {
		jsonResult.Errors["server_placeholders"] = result.ServerPlaceholderViolations
	}
	if len(result.OperationDescriptionViolations) > 0 {
		if strict {
			jsonResult.Errors["missing_operation_descriptions"] = result.OperationDescriptionViolations
		} else {
			jsonResult.Warnings["missing_operation_descriptions"] = result.OperationDescriptionViolations
		}
	}

	return jsonResult
}