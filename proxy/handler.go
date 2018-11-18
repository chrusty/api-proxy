package proxy

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// handler handles all HTTP requests:
func (p *Proxy) handler(w http.ResponseWriter, req *http.Request) {

	// Set up a field-logger for this request:
	handlerLoggerFields := logrus.Fields{
		"method": req.Method,
		"path":   req.URL.Path,
	}
	fieldLogger := p.Logger.WithFields(handlerLoggerFields)

	// Determine the destination address:
	destinationAddress, err := p.buildDestinationAddress(req.URL.Path)
	if err != nil {
		errorResponse := newErrorResponse("Unable to build destination address", err)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse.JSON())
		return
	}

	// Add this to the logger:
	fieldLogger = fieldLogger.WithField("destination", destinationAddress)

	// Prepare a proxy request:
	// requestBody := req.Body
	proxyURL := fmt.Sprintf("http://%s%s", destinationAddress, req.RequestURI)
	proxyRequest, err := http.NewRequest(req.Method, proxyURL, nil)
	if err != nil {
		errorResponse := newErrorResponse("Unable to build proxy request", err)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse.JSON())
		return
	}

	// Set headers:
	proxyRequest.Header = req.Header
	proxyRequest.Header.Set("Host", req.Host)
	proxyRequest.Header.Set("X-Forwarded-For", req.RemoteAddr)

	// Make the proxy request:
	proxyResponse, err := p.httpClient.Do(proxyRequest)
	if err != nil {
		errorResponse := newErrorResponse("Unable to proxy request", err)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		w.WriteHeader(http.StatusBadGateway)
		w.Write(errorResponse.JSON())
		return
	}
	defer proxyResponse.Body.Close()

	// Forward the response:
	fieldLogger.WithField("status_code", proxyResponse.StatusCode).Debug("Successfully proxied request")
	w.WriteHeader(proxyResponse.StatusCode)
}
