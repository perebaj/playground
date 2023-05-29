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

.PHONY: runk6
runk6:
	k6 run --vus 20 --duration 10s load_test.js