BINARY_NAME=im-switch
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

.PHONY: build clean install uninstall test

build:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	cp $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_NAME)

build-release:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/im-switch
	chmod +x $(INSTALL_DIR)/im-switch

uninstall:
	rm -f $(INSTALL_DIR)/im-switch

clean:
	rm -rf $(BUILD_DIR)

test: build
	@echo "Testing current input source:"
	./$(BUILD_DIR)/$(BINARY_NAME)
	@echo "\nTesting list input sources:"
	./$(BUILD_DIR)/$(BINARY_NAME) -l
	@echo "\nTesting help:"
	./$(BUILD_DIR)/$(BINARY_NAME) -h

run: build
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)
