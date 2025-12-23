BUILD_DIR=build
BUILD_NAME=main
BUILD=$(BUILD_DIR)/$(BUILD_NAME)

run: build
	@./$(BUILD)

build : clean
	@go build -o $(BUILD) cmd/api/main.go

clean:
	@rm -rf $(BUILD_DIR)

.PHONY: clean