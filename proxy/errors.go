package proxy

import (
	"encoding/json"
	"fmt"
)

type errorResponse struct {
	Code    string `json:"code,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func newErrorResponse(message string, err error) *errorResponse {
	return &errorResponse{
		Error:   fmt.Sprintf("%s", err),
		Message: message,
	}
}

func (er *errorResponse) JSON() (jsonBytes []byte) {
	jsonBytes, _ = json.Marshal(er)
	return jsonBytes
}
