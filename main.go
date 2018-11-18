package main

import (
	"flag"

	"github.com/chrusty/api-proxy/proxy"
)

var (
	config = &proxy.Config{}
)

// Load command line flags into a default config:
func init() {
	flag.IntVar(&config.APIPort, "api_port", 8080, "The port to proxy reqeusts to on your APIs")
	flag.StringVar(&config.APIDomainName, "api_domain", "", "A DNS domain name to append to API hostnames")
	flag.StringVar(&config.APIHostnamePrefix, "api_hostname_prefix", "api", "A DNS domain name to append to API hostnames")
	flag.StringVar(&config.APIHostnameSeparator, "api_hostname_separator", "-", "The separator to use when building hostnames")
	flag.StringVar(&config.ProxyListenAddress, "listen_address", ":8080", "The address to listen to requests on")
	flag.BoolVar(&config.LogJSON, "log_json", false, "Emit log messages as JSON instead of text")
	flag.StringVar(&config.LogLevel, "log_level", "debug", "The level to log at [trace, debug, warning, error]")
	flag.Parse()
}

// Make a new API proxy and start it up:
func main() {

	// Make a new API proxy with our config:
	apiProxy, err := proxy.New(config)
	if err != nil {
		panic(err)
	}

	// Serve requests:
	apiProxy.Start()
}
