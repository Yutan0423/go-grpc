protoc:
	@protoc -I. --go_out=. --go-grpc_out=. proto/*.proto

client:
	@go run client/main.go

server:
	@go run server/main.go