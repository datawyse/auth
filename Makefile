-include .env

#VERSION                 := $(shell git describe --tags)
GOCMD                   := go
GOTEST                  := $(GOCMD) test
BUILD                   := $(shell git rev-parse --short HEAD)
PROJECTNAME             := $(shell basename "$(PWD)")
IMAGE_NAME              := dataflow/dataflow
STDERR                  := /tmp/.$(PROJECTNAME)-stderr.txt # Redirect error output to a file, so we can show it in development mode.
PID                     := /tmp/.$(PROJECTNAME)-api-server.pid # PID file will store the server process id when it's running on development mode
SERVER_OUT              := "bin/dataflow"
ENTRYPOINT              := "scripts/image.sh"
PKG                     := "main.go"
SERVER_PKG_BUILD        := "${PKG}"
EXPORT_RESULT           ?= false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all test build vendor

all: help

run:
	@echo "starting server ..."
	@go run main.go serve

test:
	@echo "running test"
	@go test -v ./...

clean:
	@echo "cleaning dependencies"
	@go mod tidy

build: ## Build the binary file for server
	@echo "building ..."
	@go clean -modcache
	@go build -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

watch:
	@echo "running in watch mode ..."
	@ulimit -n 1000 #increase the file watch limit, might required on MacOS
	@reflex -s -r '\.go$$' make run

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)