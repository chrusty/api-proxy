package proxy

import (
	"fmt"
	"strings"
)

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
