help:
	@cat Makefile | grep '^[a-z]'

build:
	go build -o bin/api ./src

run:
	go run ./src
