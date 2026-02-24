NAME := linodectl
BIN_DIR ?= ./bin

.PHONY: install
install:
	go install .

.PHONY: test
test:
	go test -v -race -cover -coverpkg=./... -coverprofile=./coverage.out ./...

.PHONY: generate
generate:
	go generate ./...

DOCS_DIR := ./docs

.PHONY: clean-docs
clean-docs:
	rm -rf $(DOCS_DIR)
	mkdir $(DOCS_DIR)

.PHONY: docs
docs: clean-docs
	go run ./docgen

.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(NAME) .
