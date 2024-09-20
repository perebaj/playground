version=$(shell git rev-parse --short HEAD)
image_name=service1:${version}
imagedev_name=service1:dev
devrun=docker-compose run --rm service1

.PHONY: dev/image
dev/image:
	@echo "Building image"
	@docker build -t ${imagedev_name} -f Dockerfile.dev .

.PHONY: dev/%
dev/%: dev/image
	@echo "Running on container dev"
	@$(devrun) make ${*}

dev/integration:
	@echo "Running tests"
	docker-compose up service1 -d
	sleep 3
	docker-compose exec service1 /bin/sh -c "make integration-test"
	docker-compose down

.PHONY: save-requirements
save-requirements:
	@echo "Saving requirements.txt"
	@poetry export -f requirements.txt --output requirements.txt

.PHONY: image
image:
	@echo "Building image"
	@docker build -t ${image_name} .

.PHONY: run
run:
	@echo "Running container"
	docker run -d -p 80:80 ${image_name}

.PHONY: runapp
runapp:
	@echo "Running app"
	@poetry run fastapi run main.py --port 8000 --host 0.0.0.0

.PHONY: integration-test
integration-test:
	@echo "Running tests"
	@python -m pytest tests/integration

.PHONY: unit-test
unit-test:
	@echo "Running tests"
	@python -m pytest --ignore=tests/integration

.PHONY: dev/integration-test
dev/integration-test:
	@echo "Running tests"
	docker-compose run --rm -d service1
	docker-compose run --rm integration-test
	docker-compose down
