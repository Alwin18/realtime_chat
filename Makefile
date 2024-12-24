.PHONY:
	build run

build:
	go build -o bin/main cmd/main.go

run: build
	./bin/main

hot-reload:
	nodemon --exec go run cmd/main.go --signal SIGTERM

all: build run