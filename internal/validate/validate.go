package validate

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type CheckResult struct {
	OperationSummaryViolations           []string
	OperationDescriptionViolations       []string
	SuccessResponseViolations            []string
	SuccessResponseDescriptionViolations []string
	ErrorResponseViolations              []string
	ErrorResponseDescriptionViolations   []string
	ServerPlaceholderViolations          []string
	MissingServers                       bool
}


func CheckServer(server *openapi3.Server, result *CheckResult) {

    if strings.Contains(server.URL, "example.com") || strings.Contains(server.URL, "localhost") {
        result.ServerPlaceholderViolations = append(result.ServerPlaceholderViolations, server.URL)
    }
}

func CheckOperation(op *openapi3.Operation, path string, result *CheckResult) {

	if strings.TrimSpace(op.Summary) == "" {
		result.OperationSummaryViolations = append(result.OperationSummaryViolations, path)
	}
	if strings.TrimSpace(op.Description) == "" {
		result.OperationDescriptionViolations = append(result.OperationDescriptionViolations, path)
	}

	has2xx := false
	hasError := false
	has2xxWithDescription := false
	hasErrorWithDescription := false

	if op.Responses != nil {
		for code, ref := range op.Responses.Map() {
			if ref == nil || ref.Value == nil {
				continue
			}

			if strings.HasPrefix(code, "2") {
				has2xx = true

				if ref.Value.Description != nil && strings.TrimSpace(*ref.Value.Description) != "" {
					has2xxWithDescription = true
				}
			}

			if code == "default" || strings.HasPrefix(code, "4") || strings.HasPrefix(code, "5") {
				hasError = true

				if ref.Value.Description != nil && strings.TrimSpace(*ref.Value.Description) != "" {
					hasErrorWithDescription = true
				}
			}
		}
	}

	if !has2xx {
		result.SuccessResponseViolations = append(result.SuccessResponseViolations, path)
	} else if !has2xxWithDescription {
		result.SuccessResponseDescriptionViolations = append(result.SuccessResponseDescriptionViolations, path)
	}

	if !hasError {
		result.ErrorResponseViolations = append(result.ErrorResponseViolations, path)
	} else if !hasErrorWithDescription {
		result.ErrorResponseDescriptionViolations = append(result.ErrorResponseDescriptionViolations, path)
	}
}
