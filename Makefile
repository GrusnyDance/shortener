# на моем линуксе докер запускается c docker compose, без dash
run_pg:
ifeq ($(shell uname), Linux)
	docker compose up server_postgres --build --force-recreate
else
	docker-compose up server_postgres --build --force-recreate
endif

run_cache:
ifeq ($(shell uname), Linux)
	docker compose up server_cache --build --force-recreate
else
	docker-compose up server_cache --build --force-recreate
endif

test:
	go test --race ./internal/service ./internal/storage/cache ./pkg/hasher

stop:
ifeq ($(shell uname), Linux)
	docker compose down
else
	docker-compose down
endif

clean:
# возможно на маке не нужны отражающие символы
	docker stop $$(docker ps -a -q) || true
	docker rm $$(docker ps -a -q) || true
	docker rmi $$(docker images -a -q) || true

###############################################################################################################

proto-generate:
	protoc -I ./proto --go_out=proto/generate \
	--go-grpc_out=proto/generate \
	--grpc-gateway_out=proto/generate \
	shortener.proto

mock:
	mockgen -source=internal/entities/repository.go \
	-destination=internal/storage/postgres/repository/tests/mock_instance.go

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