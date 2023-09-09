.PHONY: run
run:
	@echo "=> Running app..."
	@go run ./...

.PHONY: install
install:
	@echo "=> Installing dependencies..."
	@go mod tidy

.PHONY: test
test:
	echo "=> Running tests..."
	@go test ./...
