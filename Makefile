all: run

build:
	@go build -o bin/blockchain ./cmd/node/main.go

run: build
	@./bin/blockchain

test: 
	@go test -v ./...

.PHONY: proto
# Generate gRPC client and server code
proto:
	protoc --proto_path=proto \
		--go_out=./internal/genproto --go_opt=paths=source_relative \
		--go-grpc_out=./internal/genproto --go-grpc_opt=paths=source_relative \
		proto/*.proto
