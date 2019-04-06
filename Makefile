APP_NAME="api-proxy"
BUILD_VERSION="latest"

build:
	@echo " => Building bin/api-proxy ..."
	@mkdir -p bin
	@go build -o bin/api-proxy

docker:
	@echo " => Building bin/api-proxy.linux-amd64 ..."
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -o bin/api-proxy.linux-amd64
	@echo " => Building Docker image ${APP_NAME}:${BUILD_VERSION} ..."
	@docker build -t $(APP_NAME):$(BUILD_VERSION) . >/dev/null
