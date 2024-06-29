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

#запуск тестового окружения при помощи docker-compose, 
up-test:
	docker compose up test
down-test:
	docker compose down test
test:
	go test ./... 
#запуск интеграционных тестов, 
integration-tests: up-test
	docker compose exec test go test ./tests/integration_tests -v
#запуск Unit-тестов, 
#запуск скрипта миграций, 
#очищение базы от тестовых данных
