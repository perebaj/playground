## Install dependencies
.PHONY: install
install:
	@echo "Installing dependencies..."
	@poetry install

## enable the virtual environment
.PHONY: venv
venv:
	@echo "Activating virtual environment..."
	@poetry shell

## Run the application
.PHONY: run
run:
	@echo "Running the application..."
	@poetry run python main.py

## save-requirements transcript the poetry.lock file to requirements.txt
.PHONY: save-requirements
save-requirements:
	@echo "Saving requirements..."
	@poetry export --without-hashes --format=requirements.txt > requirements.txt

## Display help for all targets
.PHONY: help
help:
	@awk '/^.PHONY: / { \
		msg = match(lastLine, /^## /); \
			if (msg) { \
				cmd = substr($$0, 9, 100); \
				msg = substr(lastLine, 4, 1000); \
				printf "  ${GREEN}%-30s${RESET} %s\n", cmd, msg; \
			} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
