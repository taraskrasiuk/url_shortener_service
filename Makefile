BIN_PATH="./bin/app"

build:
	GOOS=darwin go build -o $(BIN_PATH) ./cmd/app/
run: build
	./$(BIN_PATH)
test:
	go test ./...
