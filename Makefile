.PHONY: build run test docker

build:
	go build -o bin/ttldb ./cmd/octopus-server

run:
	go run ./cmd/octopus-server

test:
	go test -v ./...

docker:
	docker-compose build

docker-run:
	docker-compose up -d