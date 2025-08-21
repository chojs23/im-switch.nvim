BINARY_NAME=im-switch
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

# Go source files
GO_SOURCES := $(wildcard *.go)

.PHONY: build clean install uninstall test force-build

# Build only if Go sources are newer than the binary or if binary doesn't exist
$(BINARY_NAME): $(GO_SOURCES) go.mod
	mkdir -p $(BUILD_DIR)
ifeq ($(shell uname),Darwin)
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
else
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
endif
	cp $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_NAME)

# Alias for the binary target
build: $(BINARY_NAME)

# Force rebuild regardless of timestamps
force-build:
	mkdir -p $(BUILD_DIR)
ifeq ($(shell uname),Darwin)
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
else
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
endif
	cp $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_NAME)

build-release: $(GO_SOURCES) go.mod
	mkdir -p $(BUILD_DIR)
ifeq ($(shell uname),Darwin)
	CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .
else
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .
endif

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/im-switch
	chmod +x $(INSTALL_DIR)/im-switch

uninstall:
	rm -f $(INSTALL_DIR)/im-switch

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

test: build
	@echo "Testing current input source:"
	./$(BUILD_DIR)/$(BINARY_NAME)
	@echo "\nTesting list input sources:"
	./$(BUILD_DIR)/$(BINARY_NAME) -l
	@echo "\nTesting help:"
	./$(BUILD_DIR)/$(BINARY_NAME) -h

run: build
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)
