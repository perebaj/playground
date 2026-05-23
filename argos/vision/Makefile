.PHONY: help install test lint format build run list clean
.DEFAULT_GOAL := help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies via uv sync
	uv sync

test: ## Run tests. Use `make test testcase=<pattern>` to filter
	uv run pytest $(if $(testcase),-k $(testcase) -v)

lint: ## Run ruff in check mode (does not modify files)
	uv run ruff check .
	uv run ruff format --check .

format: ## Apply ruff auto-fixes and formatting
	uv run ruff check --fix .
	uv run ruff format .

build: ## Build the celebrity embedding cache. Use `make build dataset=hf1000` to switch
	uv run argos-vision build $(if $(dataset),--dataset $(dataset))

run: ## Run face matching against the cached database. Use `make run dataset=hf1000` to switch
	uv run argos-vision match $(if $(dataset),--dataset $(dataset))

list: ## List datasets and which ones have caches built
	uv run argos-vision list

clean: ## Remove .venv, cached embeddings, and Python/pytest caches
	rm -rf .venv data __pycache__ */__pycache__ .pytest_cache

