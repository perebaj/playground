.PHONY: fmt
fmt:
	poetry run black .
	poetry run isort .

.PHONY: lint
lint:
	poetry run flake8 .
	poetry run mypy .

.PHONY: test
test:
	poetry run pytest . --ignore=tests/integration

.PHONY: integration-test
integration-test:
	poetry run pytest tests/integration
