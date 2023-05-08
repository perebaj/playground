.PHONY: dev/start
dev/start:
	docker-compose up -d

.PHONY: dev/stop
dev/stop:
	docker-compose down --remove-orphans

.PHONY: app
app:
	poetry run uvicorn main:app --reload

.PHONY: dev/removeall
dev/removeall:
	docker-compose down --rmi local --volumes --remove-orphans 