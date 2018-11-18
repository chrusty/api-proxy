FROM alpine:3.8

COPY ./bin/api-proxy.linux-amd64 /api-proxy

ENV PROXY_API_DOMAIN=""
ENV PROXY_API_HOSTNAME_PREFIX="api"
ENV PROXY_API_HOSTNAME_SEPARATOR="-"
ENV PROXY_API_PORT="8080"
ENV PROXY_LISTEN_ADDRESS=":8080"
ENV PROXY_LOG_JSON="false"
ENV PROXY_LOG_LEVEL="info"

EXPOSE 8080

ENTRYPOINT /api-proxy -api_domain="${PROXY_API_DOMAIN}" -api_hostname_prefix="${PROXY_API_HOSTNAME_PREFIX}" -api_hostname_separator="${PROXY_API_HOSTNAME_SEPARATOR}" -api_port="${PROXY_API_PORT}" -listen_address="${PROXY_LISTEN_ADDRESS}" -log_json="${PROXY_LOG_JSON}" -log_level="${PROXY_LOG_LEVEL}"
