MODULE := backend/gateway
PROTO_DIR := ./docs/proto

wire:
	go run github.com/google/wire/cmd/wire ./cmd/server

generate-proto-rpc: generate-health-rpc

generate-health-rpc:
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=. \
		--go_opt=module=$(MODULE) \
		--go-grpc_out=. \
		--go-grpc_opt=module=$(MODULE) \
		$(PROTO_DIR)/health.proto
