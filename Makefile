
BUILD_OUT_DIR := "bin/"

API_OUT       := "bin/api"
API_MAIN_FILE := "cmd/api/main.go"

WORKER_OUT       := "bin/worker"
WORKER_MAIN_FILE := "cmd/worker/main.go"

PROTO_ROOT := proto/
RPC_ROOT := rpc/


deps:
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install github.com/twitchtv/twirp/protoc-gen-twirp@latest


.PHONY: proto-generate ## Compile protobuf to pb and twirp files
proto-generate:
	@echo "\n + Generating pb language bindings\n"
	@buf ls-files
	@buf generate

clean:
	@echo " + Removing generated files\n"
	@rm -rf $(API_OUT) $(WORKER_OUT) $(RPC_ROOT)


.PHONY: docker-build-api
docker-build-api:
	@docker build . -f build/docker/Dockerfile.api -t assignment02/api:latest

.PHONY: docker-build-worker
docker-build-worker:
	@docker build . -f build/docker/Dockerfile.worker -t assignment02/worker:latest

.PHONY: go-build-api ## Build the binary file for API server
go-build-api:
	@CGO_ENABLED=0 go build -v -o $(API_OUT) $(API_MAIN_FILE)

.PHONY: go-build-worker ## Build the binary file for worker
go-build-worker:
	@CGO_ENABLED=0 go build -v -o $(WORKER_OUT) $(WORKER_MAIN_FILE)