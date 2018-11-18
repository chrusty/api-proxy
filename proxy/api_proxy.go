package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// Proxy is a self-contained API proxy:
type Proxy struct {
	Config     *Config
	Logger     *logrus.Logger
	httpClient *http.Client
}

// New returns a configured API proxy:
func New(config *Config) (*Proxy, error) {

	// Make a new Logrus logger:
	logger := logrus.New()
	if config.LogJSON {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	// Attempt to parse the configured log-level:
	parsedLogLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logger.WithError(err).WithField("configured_log_level", config.LogLevel).Fatal("Unable to parse the configured log-level")
	}
	logger.SetLevel(parsedLogLevel)

	return &Proxy{
		Config:     config,
		Logger:     logger,
		httpClient: &http.Client{},
	}, nil
}

// Start tells the proxy to begin accepting requests:
func (p *Proxy) Start() {
	p.Logger.WithField("listen_address", p.Config.ProxyListenAddress).Info("Starting the proxy ...")

	// Attempt to listen and serve:
	http.HandleFunc("/", p.handler)
	if err := http.ListenAndServe(p.Config.ProxyListenAddress, nil); err != nil {
		p.Logger.WithError(err).Fatal("Error starting HTTP server")
	}

}

// buildDestinationAddress works out the destination API address from a given URL:
func (p *Proxy) buildDestinationAddress(requestURL string) (proxyAddress string, err error) {

	// Break up the path:
	pathComponents := strings.Split(requestURL, "/")

	// Make sure we have enough path components to build a sensible destination address:
	if len(pathComponents) < 2 {
		return "", fmt.Errorf("Can't build an address from less than 2 path components (eg '/v3/units')")
	}

	// Build the basic hostname:
	apiHostname := fmt.Sprintf("%s%s%s%s%s", p.Config.APIHostnamePrefix, p.Config.APIHostnameSeparator, pathComponents[1], p.Config.APIHostnameSeparator, pathComponents[2])

	// Add the DNS domain (if it has been configured):
	if p.Config.APIDomainName != "" {
		apiHostname = fmt.Sprintf("%s.%s", apiHostname, p.Config.APIDomainName)
	}

	// Add the port on the way out:
	return fmt.Sprintf("%s:%d", apiHostname, p.Config.APIPort), nil
}

// handler handles all HTTP requests:
func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {

	// Set up a field-logger for this request:
	handlerLoggerFields := logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}
	fieldLogger := p.Logger.WithFields(handlerLoggerFields)

	// Determine the destination address:
	destinationAddress, err := p.buildDestinationAddress(r.URL.Path)
	if err != nil {
		errorResponse := newErrorResponse("Unable to build destination address", err)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse.JSON())
		return
	}

	// Add this to the logger:
	fieldLogger = fieldLogger.WithField("destination", destinationAddress)

	// Make the proxy request:
	r.Host = destinationAddress
	proxyResponse, err := p.httpClient.Do(r)
	if err != nil {
		errorResponse := newErrorResponse("Unable to proxy request", err)
		fieldLogger.WithError(err).Warn(errorResponse.Message)
		w.WriteHeader(http.StatusBadGateway)
		w.Write(errorResponse.JSON())
		return
	}

	// Forward the response:
	w.WriteHeader(proxyResponse.StatusCode)
}
