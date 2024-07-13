INTERNAL_REPO_PATH=$(CURDIR)/
MIGRATION_FOLDER=$(INTERNAL_REPO_PATH)migrations
LOCAL_BIN:=$(CURDIR)/bin

PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc
PVZ_PROTO_PATH:="api/proto/pvz/v1"

POSTGRES_OMS_SETUP := user=postgres password=password dbname=oms host=localhost port=5432 sslmode=disable
POSTGRES_TEST_SETUP := user=postgres password=password dbname=test host=localhost port=5432 sslmode=disable

# build docker image
build:
	docker compose build

up-all:
	docker compose up -d postgres
	docker compose up -d zookeeper kafka1 kafka2 kafka3

down:
	docker compose down

# run docker image
up-db:
	docker compose up -d postgres

stop-db:
	docker compose stop postgres

start-db:
	docker compose start postgres

down-db:
	docker compose down postgres

up-service:
	docker compose up -d app

stop-service:
	docker compose stop app

start-service:
	docker compose start app

down-service:
	docker compose down app

# запуск тестового окружения при помощи docker-compose
.PHONY: test-base-up
test-up-all:
	docker compose up -d zookeeper kafka1
	docker compose -f docker-compose.test.yml up -d
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_TEST_SETUP)" up

.PHONY: test-base-down
test-base-down:
	docker compose -f docker-compose.test.yml down

test:
	go test ./... -v

# запуск интеграционных тестов
integration-tests: 
	go test -tags=integration ./tests/integration -v

# запуск Unit-тестов
unit-tests:
	go test -tags=unit ./internal/cli -v
	go test -tags=unit ./internal/module -v

# запуск Suite-тестов
suite-tests:
	go test -tags=suite ./tests/suite -v

# генерация моков
gen-mocks:
	mockgen -source ./internal/cli/init.go -destination=./internal/cli/mocks/cli.go -package=mock_cli
	mockgen -source ./internal/module/init_module.go -destination=./internal/module/mocks/module.go -package=mock_module
	
# запуск скрипта миграций
.PHONY: migration-up-oms
migration-up-oms:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_OMS_SETUP)" up

.PHONY: migration-down-oms
migration-down-oms:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_OMS_SETUP)" down

.PHONY: migration-up-test
migration-up-test:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_TEST_SETUP)" up

.PHONY: migration-down-test
migration-down-test:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_TEST_SETUP)" down

.PHONY: cleanup-db-test
cleanup-db-test:
	psql "postgres://postgres:password@localhost:5432/test?sslmode=disable" -f ./cleanup_test.sql

# Установка всех необходимых зависимостей
.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest

# Вендоринг внешних proto файлов
.vendor-proto: vendor-proto/google/protobuf vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options vendor-proto/validate

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	rm -rf vendor.proto/grpc-ecosystem
		rm -rf vendor.proto/protoc-gen-openapiv2
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.proto/grpc-ecosystem && \
 	cd vendor.proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor.proto/protoc-gen-openapiv2
	mv vendor.proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor.proto/protoc-gen-openapiv2
	rm -rf vendor.proto/grpc-ecosystem

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	rm -rf vendor.proto/google/protobuf
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.proto/protobuf &&\
	cd vendor.proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor.proto/google
	mv vendor.proto/protobuf/src/google/protobuf vendor.proto/google
	rm -rf vendor.proto/protobuf

vendor-proto/google/api:
	rm -rf vendor.proto/google/api
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.proto/googleapis && \
 	cd vendor.proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p vendor.proto/google
	mv vendor.proto/googleapis/google/api vendor.proto/google
	rm -rf vendor.proto/googleapis

vendor-proto/validate:
	rm -rf vendor.proto/validate
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.proto/tmp && \
		cd vendor.proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
	mkdir -p vendor.proto/validate
	mv vendor.proto/tmp/validate vendor.proto/
	rm -rf vendor.proto/tmp

.PHONY: generate
generate: .bin-deps .vendor-proto
	mkdir -p pkg/${PVZ_PROTO_PATH}
	protoc -I api/proto \
		-I vendor.proto \
		${PVZ_PROTO_PATH}/pvz.proto \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=./pkg/${PVZ_PROTO_PATH} --go_opt=paths=source_relative\
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=./pkg/${PVZ_PROTO_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ./pkg/api/proto/pvz/v1  --grpc-gateway_opt  paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=./pkg/api/proto/pvz/v1 \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:pkg/api/proto/pvz/v1" \
		--experimental_allow_proto3_optional

