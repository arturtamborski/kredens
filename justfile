help:
	@just --list --unsorted

build:
	go build -v ./...

test:
	go build -v ./...

	./kredens help 1>/dev/null
	./kredens list 1>/dev/null
	@echo "seems to work"
