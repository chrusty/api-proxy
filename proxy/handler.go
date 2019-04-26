package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		errorResponse := &errorResponse{
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
			Message: fmt.Sprintf("Unable to build destination address: %s", err.Error()),
		}
		errorResponse.write(w)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		return
	}

	// Add this to the logger:
	fieldLogger = fieldLogger.WithField("destination", destinationAddress)

	// Read the request body:
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse := &errorResponse{
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
			Message: fmt.Sprintf("Unable to read request body: %s", err.Error()),
		}
		errorResponse.write(w)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		return
	}
	defer req.Body.Close()

	// Prepare a proxy request:
	proxyURL := fmt.Sprintf("http://%s%s", destinationAddress, req.RequestURI)
	proxyRequest, err := http.NewRequest(req.Method, proxyURL, bytes.NewReader(requestBody))
	if err != nil {
		errorResponse := &errorResponse{
			Code:    http.StatusUnprocessableEntity,
			Error:   err.Error(),
			Message: fmt.Sprintf("Unable to build proxy request: %s", err.Error()),
		}
		errorResponse.write(w)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		return
	}

	// Set headers:
	proxyRequest.Header = req.Header
	proxyRequest.Header.Set("Host", req.Host)
	proxyRequest.Header.Set("X-Forwarded-For", req.RemoteAddr)

	// Make the proxy request:
	proxyResponse, err := p.httpClient.Do(proxyRequest)
	if err != nil {
		errorResponse := &errorResponse{
			Code:    http.StatusBadGateway,
			Error:   err.Error(),
			Message: fmt.Sprintf("Unable to proxy request: %s", err.Error()),
		}
		errorResponse.write(w)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		return
	}
	defer proxyResponse.Body.Close()

	// Read the response body:
	responseBody, err := ioutil.ReadAll(proxyResponse.Body)
	if err != nil {
		errorResponse := &errorResponse{
			Code:    http.StatusInternalServerError,
			Error:   err.Error(),
			Message: fmt.Sprintf("Unable to read response body: %s", err.Error()),
		}
		errorResponse.write(w)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		return
	}

	// Get all the response headers:
	for headerKey, headerValue := range proxyResponse.Header {
		fieldLogger.Tracef("Setting response header (%s) = %s", headerKey, headerValue[0])
		w.Header().Set(headerKey, headerValue[0])
	}

	// Forward the response:
	fieldLogger.WithField("status_code", proxyResponse.StatusCode).Debug("Successfully proxied request")
	w.WriteHeader(proxyResponse.StatusCode)
	w.Write(responseBody)
}
