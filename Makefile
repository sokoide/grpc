
.PHONY: gen client server

all: gen server client

server: gen cmd/server/*.go
	go build ./cmd/server

client: gen cmd/client/*.go
	go build ./cmd/client

gen: proto/hello.proto
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/hello.proto

