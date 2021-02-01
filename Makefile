.PHONY: lint
lint:
	golint ./...

.PHONY: build
build:
	go build -o ./bin/social ./cmd/social/main.go

.PHONY: run
run:
	go run cmd/social/main.go