.PHONY: dev build preview format lint test tidy install

BINARY := mezamero
CONFIG ?= config.example.yaml
ADDR ?= :8080

# Fetch and verify modules (similar to npm install)
install:
	go mod download
	go mod verify

dev:
	go run ./cmd/mezamero -config $(CONFIG) -addr $(ADDR)

build:
	CGO_ENABLED=0 go build -trimpath -o $(BINARY) ./cmd/mezamero

preview: build
	./$(BINARY) -config $(CONFIG) -addr $(ADDR)

format:
	gofmt -w .
	@command -v goimports >/dev/null 2>&1 && goimports -w . || true

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || \
		( echo "golangci-lint not installed; try: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" >&2; exit 1 )

test:
	go test ./...

tidy:
	go mod tidy
