.PHONY: *

run:
	@echo ">>> running app ..."
	@go run main.go
	@echo ">>> ... done"

build:
	@echo ">>> building app ..."
	go build -o fem-go-crud
	@echo ">>> ... done"

lint: lint-fp lint-tp

# first-party linters
lint-fp:
	@echo ">>> running first-party linters ..."
	@go fmt ./...
	@go vet ./...
	@echo ">>> ... done"

# third-party linters
lint-tp:
	@echo ">>> running third-party linters ..."
	@go tool gofumpt -w .
	@go tool errcheck ./...
	@go tool staticcheck ./...
	@echo ">>> ... done"

test:
	@echo ">>> testing ..."
	@go test ./...
	@echo ">>> ... done"

docker:
	@echo ">>> starting docker ..."
	@docker compose up --build --force-recreate
	@echo ">>> ... done"

docker-test:
	@echo ">>> starting docker for tests ..."
	@docker compose -f docker-compose.test.yaml up --build --force-recreate
	@echo ">>> ... done"
