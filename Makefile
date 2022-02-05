
BUILD_OUT_DIR := "bin/"

API_OUT       := "bin/api"
API_MAIN_FILE := "cmd/api/main.go"

PROTO_ROOT := proto/
RPC_ROOT := rpc/


deps:
	@go install github.com/bufbuild/buf/cmd/buf
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
	@go install github.com/twitchtv/twirp/protoc-gen-twirp


.PHONY: proto-generate ## Compile protobuf to pb and twirp files
proto-generate:
	@echo "\n + Generating pb language bindings\n"
	@buf ls-files
	@buf generate

clean:
	@echo " + Removing generated files\n"
	@rm -rf $(API_OUT) $(MIGRATION_OUT) $(WORKER_OUT) $(RPC_ROOT)


