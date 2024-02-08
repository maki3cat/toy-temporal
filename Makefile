
pb:
	protoc --go_out=. --go_opt=module=github.com/maki3cat/toy-temporal --go-grpc_out=. --go-grpc_opt=module=github.com/maki3cat/toy-temporal api/workflow.proto
run:
	go run ./server
