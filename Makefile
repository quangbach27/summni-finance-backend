SERVICE := sumni-finance-backend

.PHONY: test dev

test:
	@./scripts/test.sh .e2e.env

dev:
	DEBUG=$(DEBUG) docker compose up --build $(SERVICE) -d	

stop:
	docker compose down $(SERVICE)
	