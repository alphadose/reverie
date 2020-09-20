PROJECTNAME := $(shell basename "$(PWD)")
PACKAGES := $(shell go list ./... | grep -v vendor)

# Go related variables.
PROJECTROOT := $(shell pwd)
GOBIN := $(PROJECTROOT)/bin
GOFILES := $(PROJECTROOT)/*.go

# Shell script related variables.
BUILDIR := $(PROJECTROOT)/scripts

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: default
default: help

## install: Install missing dependencies
install:
	@printf "🔨 Installing project dependencies to vendor\n" 
	@GOBIN=$(GOBIN) go get ./... && go mod vendor
	@printf "👍 Done\n"

## build: Build the project binary
build:
	@printf "🔨 Building binary $(GOBIN)/$(PROJECTNAME)\n" 
	@go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)
	@printf "👍 Done\n"

## tools: Install development tools
tools:
	@$(BUILDIR)/install_fresh.sh
	@$(BUILDIR)/install_golint.sh

## start: Start in development mode with hot-reload enabled
start: tools
	@$(PROJECTROOT)/bin/fresh

## clean: Clean build files
clean:
	@printf "🔨 Cleaning build cache\n" 
	@go clean $(PACKAGES)
	@printf "👍 Done\n"
	@-rm $(GOBIN)/$(PROJECTNAME) 2>/dev/null

## fmt: Format entire codebase
fmt:
	@printf "🔨 Formatting\n" 
	@gofmt -s -w .
	@printf "👍 Done\n"

## vet: Vet entire codebase
vet:
	@printf "🔨 Vetting\n" 
	@go vet $(PACKAGES)
	@printf "👍 Done\n"

## lint: Check codebase for style mistakes
lint:
	@printf "🔨 Linting\n"
	@golint -set_exit_status $(PACKAGES)
	@printf "👍 Done\n"

## test: Run tests
test:
	@printf "🔨 Testing\n"
	@go test -race -coverprofile=coverage.txt -covermode=atomic
	@printf "👍 Done\n"

## help: Display this help
help: Makefile
	@printf "\n Reverie: The dark side of EzFlo\n\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@printf ""
