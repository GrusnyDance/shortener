proto-generate:
	protoc -I ./proto --go_out=proto/generate \
	--go-grpc_out=proto/generate \
	--grpc-gateway_out=proto/generate \
	shortener.proto

test_post:
	curl --location --request POST 'http://localhost:8085/post' \
	--header 'Content-Type: application/json' \
    --data-raw '{"data": "lalala"}'

test_get:
	curl --location --request GET 'http://localhost:8085/get/lalala' \
	--header 'Content-Type: application/json'

run:
	go run cmd/main.go

stop:
	fuser -k 8080/tcp

#migration:
#	goose create create_table go