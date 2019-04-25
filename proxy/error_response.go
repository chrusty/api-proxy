package proxy

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *errorResponse) JSON() (jsonBytes []byte) {
	jsonBytes, _ = json.Marshal(e)
	return jsonBytes
}

func (e *errorResponse) write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	w.Write(e.JSON())
}
