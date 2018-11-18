package proxy

// Config contains everything we need to know to run an API proxy:
type Config struct {
	APIPort              int    // The port to proxy reqeusts to on your APIs
	APIDomainName        string // A DNS domain name to append to API hostnames
	APIHostnamePrefix    string // The default prefix for your API hostnames
	APIHostnameSeparator string // The separator to use when building hostnames
	ProxyListenAddress   string // The address to listen to requests on
	LogJSON              bool   // Emit log messages as JSON instead of text?
	LogLevel             string // The level to log at [trace, debug, warning, error]
}
