.PHONY: default build_and_up down install build

default:
	docker-compose up -d sammy

down:
	docker-compose down -t 5

build_and_up:
	docker-compose up --build sammy

install:
	go install -v ./...

build:
	go build -o  -v
