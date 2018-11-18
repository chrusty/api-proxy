api-proxy
=========
Routes requests to your API services (by the URL path)

Why?
----
Imagine that you need to run several APIs, and route requests to them based on the URL path:
* api-v1-cruft (/v1/cruft/*)
* api-v2-cruft (/v2/cruft/*)
* api-v1-foo (/v1/foo/*)
* api-v3-bar (/v3/bar/*)

Flags
-----
```
Usage of bin/api-proxy:
  -api_domain string
    	A DNS domain name to append to API hostnames
  -api_hostname_prefix string
    	A DNS domain name to append to API hostnames (default "api")
  -api_hostname_separator string
    	The separator to use when building hostnames (default "-")
  -api_port int
    	The port to proxy reqeusts to on your APIs (default 8080)
  -listen_address string
    	The address to listen to requests on (default ":8080")
  -log_json
    	Emit log messages as JSON instead of text
  -log_level string
    	The level to log at [trace, debug, warning, error] (default "debug")
```