NAME := linodectl

install:
	cd cmd/linodectl && go install .

test:
	go test -v -race -cover -coverpkg=./... -coverprofile=./coverage.out ./...

generate:
	go generate ./...

DOCS_DIR := ./docs

clean-docs:
	rm -rf $(DOCS_DIR)
	mkdir $(DOCS_DIR)

docs: clean-docs
	go run cmd/docgen/main.go

build:
	go build -o $(linodectl) cmd/linodectl/main.go
