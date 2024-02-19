include .env

build:
	@go build -o bin/main

prod-db:
	@echo "Starting Prod PostgreSQL..."
	@docker run -d \
		--name prod-postgres \
		-p 5432:5432 \
		-e POSTGRES_USER=${POSTGRES_USER} \
		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
		-e POSTGRES_DB=${POSTGRES_PROD_DB} \
		postgres:latest

run: prod-db build
	@echo "Starting Server at localhost:3000..."
	@./bin/main

prod-clean:
	@echo "Cleaning..."
	@docker stop prod-postgres
	@docker rm prod-postgres
	@rm -rf bin

test-db:
	@echo "Starting Test PostgreSQL..."
	@docker run -d \
		--name test-postgres \
		-p 5433:5432 \
		-e POSTGRES_USER=${POSTGRES_USER} \
		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
		-e POSTGRES_DB=${POSTGRES_TEST_DB} \
		postgres:latest

test: test-db
	@echo "Testing..."
	@go test -v ./tests
	@$(MAKE) test-clean

test-clean:
	@echo "Cleaning..."
	@docker stop test-postgres
	@docker rm test-postgres
