package resp

import (
	"encoding/json"
)

type mcpStatus string

const (
	statusSuccess      mcpStatus = "success"
	statusError        mcpStatus = "error"
	statusNotFound     mcpStatus = "not_found"
	statusInvalidInput mcpStatus = "invalid_innput"
)

type mcpSuccessResponse struct {
	Status mcpStatus   `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type mcpErrorResponse struct {
	Status     mcpStatus           `json:"status"`
	Error      mcpErrorDetails     `json:"error,omitempty"`
	Suggestion map[string][]string `json:"data,Ssuggestion"`
}

type mcpErrorDetails struct {
	Message               string `json:"message"`
	InvalidOrMissimgField string `json:"missing_field,omitempty"`
	Hint                  string `json:"hint,omitempty"`
	Details               string `json:"details,omitempty"`
}

func Success(data interface{}) string {
	resp := mcpSuccessResponse{
		Status: statusSuccess,
		Data:   data,
	}
	jsonResp, _ := json.MarshalIndent(resp, "", "  ")
	return string(jsonResp)
}

func InvalidInput(msg string, invalidOrMissimgField string, hint string) string {
	resp := mcpErrorResponse{
		Status: statusInvalidInput,
		Error: mcpErrorDetails{
			Message:               msg,
			InvalidOrMissimgField: invalidOrMissimgField,
			Hint:                  hint,
		},
	}
	jsonResp, _ := json.MarshalIndent(resp, "", "  ")
	return string(jsonResp)
}

func NotFound(msg string, fieldName string, suggestions []string) string {
	resp := mcpErrorResponse{
		Status: statusNotFound,
		Error: mcpErrorDetails{
			Message: msg,
		},
		Suggestion: map[string][]string{
			fieldName: suggestions,
		},
	}
	jsonResp, _ := json.MarshalIndent(resp, "", "  ")
	return string(jsonResp)
}

func InternalError(msg string) string {
	resp := mcpErrorResponse{
		Status: statusError,
		Error: mcpErrorDetails{
			Message: "Internal server error",
			Details: msg,
		},
	}
	jsonResp, _ := json.MarshalIndent(resp, "", "  ")
	return string(jsonResp)
}
