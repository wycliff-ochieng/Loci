SERVICE_NAME ?= loci

.PHONY: lint test build run docker-build

lint:
	@echo "Running gofmt + go vet"
	@UNFORMATTED=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "These files are not gofmt-formatted:"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi
	go vet ./...

test:
	go test -v ./...

build:
	mkdir -p bin
	go build -v -o bin/$(SERVICE_NAME) ./cmd/main.go

run:
	@echo "Running $(SERVICE_NAME) locally..."
	go run ./cmd/main.go

docker-build:
	docker build -t $(SERVICE_NAME):latest .