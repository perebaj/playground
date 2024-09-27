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
	poetry run pytest .
