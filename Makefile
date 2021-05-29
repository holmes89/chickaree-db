build-proto:
	protoc -I=. --go_out=. --go-grpc_out=. ./client.proto 