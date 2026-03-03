GO ?= /usr/local/go/bin/go
BINARY ?= sshconfctl

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X 'github.com/chalvinwz/sshconfctl/internal/cmd.Version=$(VERSION)'

.PHONY: build release test fmt run clean

build:
	$(GO) build -o $(BINARY) ./cmd/sshconfctl

release:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/sshconfctl

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

run:
	$(GO) run ./cmd/sshconfctl

clean:
	rm -f $(BINARY)
