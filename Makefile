APP_NAME := gofind

.PHONY: build test run bench

build:
	go build -o bin/$(APP_NAME) ./cmd/gofind

test:
	go test ./...

run:
	go run ./cmd/gofind -pattern panic -root .

bench:
	go test -bench=. ./tests/benchmark/...
