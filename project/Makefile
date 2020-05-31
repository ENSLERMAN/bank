.PHONY: build
build:
	go build -v ./main.go

.PHONY: run
run:
	go run ./main.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := run