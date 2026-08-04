BINARY := devarchitect
CMD := ./cmd/devarchitect

.PHONY: build test fmt vet lint check clean run

build:
	go build -o bin/$(BINARY) $(CMD)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

# Runs everything CI runs, so you can catch problems before pushing.
check: fmt vet test

run: build
	./bin/$(BINARY) $(ARGS)

clean:
	rm -rf bin
