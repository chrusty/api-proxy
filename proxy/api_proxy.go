package proxy

import (
	"net/http"

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
