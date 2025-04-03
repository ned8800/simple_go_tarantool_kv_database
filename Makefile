COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
COVERIGNORE_FILE := cover_ignore.txt
APP_NAME := tarantoolHttp.a

UNAME := $(shell uname -s)
ifeq ($(UNAME), Linux)
    OPEN_CMD = xdg-open
else ifeq ($(UNAME), Darwin)
    OPEN_CMD = open
else
    OPEN_CMD = start
endif

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

build:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(APP_NAME)$(if $(filter windows,$(GOOS)),.exe,) ./main.go

build-windows:
	@make build GOOS=windows GOARCH=amd64

build-linux:
	@make build GOOS=linux GOARCH=amd64

build-darwin:
	@make build GOOS=darwin GOARCH=arm64

run: build
	@./$(APP_NAME)

clean:
	@rm -f $(APP_NAME)

.PHONY: build run clean
