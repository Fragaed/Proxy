.PHONY: build test docker-build run lint

build:
	go build -o app ./cmd/main.go

test:
	go test ./...

docker-build:
	docker-compose up -d

run:
	docker-compose run --rm app ./app

lint:
	golangci-lint run
