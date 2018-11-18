FROM alpine:3.8
COPY ./bin/api-proxy.linux-amd64 /api-proxy
EXPOSE 8080
ENTRYPOINT ["/api-proxy"]
