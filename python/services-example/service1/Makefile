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

.PHONY: poetry-config
poetry-config:
	poetry config repositories.gcp https://us-python.pkg.dev/truepay-tools/python/simple/
	poetry config http-basic.gcp oauth2accesstoken $$(gcloud auth print-access-token)
