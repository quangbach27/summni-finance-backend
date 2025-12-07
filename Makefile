SERVICE := sumni-finance-backend

.PHONY: test dev stop lint

test:
	@./scripts/test.sh .e2e.env

dev:
	DEBUG=$(DEBUG) docker compose up --build $(SERVICE) -d	

logs:
	docker logs -f $(SERVICE)

stop:
	docker compose down $(SERVICE)
	
lint:
	golangci-lint run