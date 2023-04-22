proto-generate:
	protoc -I ./proto --go_out=proto/generate \
	--go-grpc_out=proto/generate \
	--grpc-gateway_out=proto/generate \
	hello.proto

