# на моем линуксе докер запускается c docker compose, без dash
run:
	docker compose up -d --build --force-recreate

run_not_detached:
	docker compose up --build --force-recreate

stop:
	docker compose down

clean:
# возможно на маке не нужны отражающие символы
	docker stop $$(docker ps -a -q) || true
	docker rm $$(docker ps -a -q) || true
	docker rmi $$(docker images -a -q) || true

docker-test-server:
	docker build --no-cache -t server -f Dockerfile .
	docker run --privileged server

proto-generate:
	protoc -I ./proto --go_out=proto/generate \
	--go-grpc_out=proto/generate \
	--grpc-gateway_out=proto/generate \
	shortener.proto

###############################################################################################################

test_post:
	curl --location --request POST 'http://localhost:8085/post' \
	--header 'Content-Type: application/json' \
    --data-raw '{"LinkToHash": "lalala"}'

test_get:
	curl --location --request GET 'http://localhost:8085/get/8ja9JS8xBK' \
	--header 'Content-Type: application/json'

#run:
#	go run cmd/main.go
#
#stop:
#	fuser -k 8080/tcp

#migration:
#	goose create create_table go