.PHONY: fmt
fmt:
	poetry run black .
	poetry run isort .

.PHONY: lint
lint:
	poetry run flake8 .

.PHONY: test
test:
	poetry run pytest . --ignore=tests/integration

.PHONY: integration-test
integration-test:
	poetry run pytest tests/integration

.PHONY: poetry-config
poetry-config:
	poetry config repositories.gcp https://us-central1-python.pkg.dev/blup-432616/libtest/simple
	poetry config http-basic.gcp oauth2accesstoken $$(gcloud auth print-access-token)
