help:
	@just --list --unsorted

lint:
	golangci-lint run --fix

build:
	go build -v ./...

test: build
	./kredens help 1>/dev/null
	./kredens list 1>/dev/null
	@echo "seems to work"

release:
	goreleaser --snapshot --clean
