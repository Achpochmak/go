INTERNAL_REPO_PATH=$(CURDIR)/
MIGRATION_FOLDER=$(INTERNAL_REPO_PATH)migrations

POSTGRES_OMS_SETUP := user=postgres password=password dbname=oms host=localhost port=5432 sslmode=disable
POSTGRES_TEST_SETUP := user=postgres password=password dbname=test host=localhost port=5432 sslmode=disable

# build docker image
build:
	docker compose build

up-all:
	docker compose up -d postgres app

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
test-base-up:
	docker compose -f docker-compose.test.yml up -d

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
