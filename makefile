VERSION=$(shell git describe --tags)

build:
	go build -o bin/$(shell basename $(CURDIR)) -ldflags "-s -w -X main.VERSION=$(VERSION)" -trimpath .

install:
	go install -ldflags "-s -w -X main.VERSION=$(VERSION)" -trimpath .
