package validate

import (
	"fmt"
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
	OperationIdViolations                []string
	OperationTagViolations               []string
	ParamDescriptionViolations           []string
	MissingServers                       bool
}

func (result *CheckResult) HasErrors() bool {
	return len(result.OperationSummaryViolations) > 0 ||
		len(result.SuccessResponseViolations) > 0 ||
		len(result.ErrorResponseViolations) > 0 ||
		len(result.SuccessResponseDescriptionViolations) > 0 ||
		len(result.ErrorResponseDescriptionViolations) > 0 ||
		result.MissingServers ||
		len(result.ServerPlaceholderViolations) > 0
}

func (result *CheckResult) HasWarnings() bool {
	return len(result.OperationDescriptionViolations) > 0 ||
		len(result.OperationIdViolations) > 0 ||
		len(result.OperationTagViolations) > 0 
}

func CheckServer(server *openapi3.Server, result *CheckResult, blockList []string) {

	for _, url := range blockList {
		if server.URL == url {
			result.ServerPlaceholderViolations = append(result.ServerPlaceholderViolations, server.URL)
			}
		}
}

func checkParam(param *openapi3.Parameter, path string, result *CheckResult) {
	if param == nil {
		return
	}
	if strings.TrimSpace(param.Description) == "" {
		result.ParamDescriptionViolations = append(result.ParamDescriptionViolations, path)
	}
}

func CheckOperation(op *openapi3.Operation, path string, result *CheckResult) {

	if strings.TrimSpace(op.Summary) == "" {
		result.OperationSummaryViolations = append(result.OperationSummaryViolations, path)
	}
	if strings.TrimSpace(op.OperationID) == "" {
		result.OperationIdViolations = append(result.OperationIdViolations, path)
	}
	if len(op.Tags) == 0 {
		result.OperationTagViolations = append(result.OperationTagViolations, path)
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

	for _, paramItem := range op.Parameters {
        if paramItem == nil || paramItem.Value == nil {
            continue
        }
        if paramItem.Value.In == "path" {
            checkParam(paramItem.Value, path, result)
        } else {
            checkParam(paramItem.Value, fmt.Sprintf("%s > query: %s", path, paramItem.Value.Name), result)
        }
	}
}

func CheckPaths(doc *openapi3.T, result *CheckResult) {
	for path, pathItem := range doc.Paths.Map() {
		if pathItem == nil {
			continue
		}
		if pathItem.Get != nil {
			CheckOperation(pathItem.Get, fmt.Sprintf("GET %s", path), result)
		}
		if pathItem.Post != nil {
			CheckOperation(pathItem.Post, fmt.Sprintf("POST %s", path), result)
		}
		if pathItem.Put != nil {
			CheckOperation(pathItem.Put, fmt.Sprintf("PUT %s", path), result)
		}
		if pathItem.Patch != nil {
			CheckOperation(pathItem.Patch, fmt.Sprintf("PATCH %s", path), result)
		}
		if pathItem.Delete != nil {
			CheckOperation(pathItem.Delete, fmt.Sprintf("DELETE %s", path), result)
		}
	}
}
