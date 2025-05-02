BINARY_NAME=myapp

# Definition Go
GO=go
GOCMD=$(GO)
port=8080
.PHONY: test #biar tidak dianggap comant

# Flag Go
GO_FLAGS=-v

# Running binary
build:
	$(GOCMD) build $(GO_FLAGS) -o $(BINARY_NAME) ./blockchain_server/main.go

# running main apps
run:
	$(GOCMD) run ./blockchain_server/main.go --port=$(port)

# running all unit test
test:
	$(GOCMD) test -v ./...

# Format kode Go
fmt:
	$(GOCMD) fmt ./...