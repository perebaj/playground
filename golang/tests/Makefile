## create a new migration file. Usage `make migration/create name=<migration_name>`
.PHONY: migration/create
migration/create:
	@echo "Creating a new migration..."
	@go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir postgres/migrations -seq $(name)

## Run integration tests. Usage `make integration-test` or `make integration-test testcase="TestFunctionName"` to run an isolated tests
.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	if [ -n "$(testcase)" ]; then \
		go test ./... -tags integration -timeout 10s -v -run="^$(testcase)$$" ; \
	else \
		go test ./... -tags integration -timeout 10s; \
	fi

## Start the development server
.PHONY: dev/start
dev/start:
	@echo "Starting the development server..."
	@docker-compose up -d

## Stop the development server
.PHONY: dev/stop
dev/stop:
	@echo "Stopping the development server..."
	@docker-compose down
