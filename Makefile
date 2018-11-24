.DEFAULT_GOAL := all
PKGS := $(shell go list ./cmd/ ./pkg/* | grep -v /vendor)
BINARY := alien-invasion
GO_BIN_DIR := $(GOPATH)/bin
DEP := $(GO_BIN_DIR)/dep

# Install dep if we don't have it at the moment
$(DEP):
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Make sure the dependencies have been initialised
Gopkg.toml: $(DEP)
	dep init

clean:
	rm -f $(BINARY)

# Install any necessary dependencies
vendor: Gopkg.toml
	dep ensure

# Build our binary
alien-invasion: cmd/main.go
	go build -o $(BINARY) cmd/main.go

all: clean vendor test alien-invasion

# Run all of our tests
test: clean vendor
	go test -v $(PKGS)

.PHONY: clean vendor

