BUILD_DIR=./build
BINARY_NAME=translator
TRANSLATOR_BIN=$(BUILD_DIR)/$(BINARY_NAME)
MAIN_FILE=./cmd/translator/main.go

.PHONY: run build clean test lint static migrate

run:
	$(TRANSLATOR_BIN) --config='config.yml'

build: clean
	go build -o $(TRANSLATOR_BIN) $(MAIN_FILE)

clean:
	@rm -rf $(BUILD_DIR)

test:
	go test -v ./...

lint:
	golangci-lint run

static:
	staticcheck ./...

migrate:
	migrate -database 'postgresql://postgres:postgres@localhost:5432/dictionary' -path ./db/migrations up