listen:
	go run ./cmd/test
	
codegen:
	protoc --proto_path=proto --go-grpc_out=proto --go-grpc_opt=paths=source_relative --go_out=proto --go_opt=paths=source_relative proto/broker.proto