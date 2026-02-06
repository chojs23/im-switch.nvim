BINARY_NAME=im-switch
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin
WSL_WIN_OUT ?= /mnt/c/Tools/im-switch.exe

# Detect operating system
UNAME_S := $(shell uname -s 2>/dev/null || echo "Windows")
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    MKDIR := cmd /c "if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)"
    COPY := copy
    RM := del /q
    RMDIR := rmdir /s /q
    PATH_SEP := \\
    EXEC_PREFIX := .\\
    EXEC_EXT := .exe
else
    DETECTED_OS := $(UNAME_S)
    MKDIR := mkdir -p $(BUILD_DIR)
    COPY := cp
    RM := rm -f
    RMDIR := rm -rf
    PATH_SEP := /
    EXEC_PREFIX := ./
    EXEC_EXT :=
endif

# Go source files
GO_SOURCES := $(wildcard *.go)
.PHONY: build clean install uninstall test force-build build-release run build-wsl-win test-wsl-win

# Build only if Go sources are newer than the binary or if binary doesn't exist
$(BINARY_NAME): $(GO_SOURCES) go.mod
	$(MKDIR)
ifeq ($(DETECTED_OS),Darwin)
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	$(COPY) $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_NAME)
else ifeq ($(DETECTED_OS),Windows)
	cmd /c "set CGO_ENABLED=0 && go build -o $(BUILD_DIR)/$(BINARY_NAME).exe ."
	$(COPY) $(BUILD_DIR)$(PATH_SEP)$(BINARY_NAME).exe .$(PATH_SEP)$(BINARY_NAME).exe
else
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	$(COPY) $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_NAME)
endif

# Alias for the binary target
build: $(BINARY_NAME)

# Force rebuild regardless of timestamps
force-build:
	$(MKDIR)
ifeq ($(DETECTED_OS),Darwin)
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	$(COPY) $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_NAME)
else ifeq ($(DETECTED_OS),Windows)
	cmd /c "set CGO_ENABLED=0 && go build -o $(BUILD_DIR)/$(BINARY_NAME).exe ."
	$(COPY) $(BUILD_DIR)$(PATH_SEP)$(BINARY_NAME).exe .$(PATH_SEP)$(BINARY_NAME).exe
else
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	$(COPY) $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_NAME)
endif

build-release: $(GO_SOURCES) go.mod
	$(MKDIR)
ifeq ($(DETECTED_OS),Darwin)
	CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .
else ifeq ($(DETECTED_OS),Windows)
	cmd /c "set CGO_ENABLED=0 && go build -ldflags=\"-s -w\" -o $(BUILD_DIR)/$(BINARY_NAME).exe ."
else
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .
endif

build-wsl-win: $(GO_SOURCES) go.mod
	@mkdir -p "$(dir $(WSL_WIN_OUT))"
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o "$(WSL_WIN_OUT)" .
	@echo "Built Windows binary at $(WSL_WIN_OUT)"

test-wsl-win: build-wsl-win
	"$(WSL_WIN_OUT)"
	"$(WSL_WIN_OUT)" -l
	"$(WSL_WIN_OUT)" -h

install: build
	$(COPY) $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/im-switch
	chmod +x $(INSTALL_DIR)/im-switch

uninstall:
	$(RM) $(INSTALL_DIR)/im-switch

clean:
	$(RMDIR) $(BUILD_DIR) 2>nul || echo Directory cleaned
ifeq ($(DETECTED_OS),Windows)
	$(RM) $(BINARY_NAME).exe 2>nul || echo Binary cleaned
else
	$(RM) $(BINARY_NAME)
endif

test: build
	@echo "Testing current input source:"
ifeq ($(DETECTED_OS),Windows)
	$(EXEC_PREFIX)$(BUILD_DIR)$(PATH_SEP)$(BINARY_NAME).exe
	@echo.
	@echo "Testing list input sources:"
	$(EXEC_PREFIX)$(BUILD_DIR)$(PATH_SEP)$(BINARY_NAME).exe -l
	@echo.
	@echo "Testing help:"
	$(EXEC_PREFIX)$(BUILD_DIR)$(PATH_SEP)$(BINARY_NAME).exe -h
else
	@if $(EXEC_PREFIX)$(BUILD_DIR)/$(BINARY_NAME) --status >/dev/null 2>&1; then \
		$(EXEC_PREFIX)$(BUILD_DIR)/$(BINARY_NAME); \
		echo "\nTesting list input sources:"; \
		$(EXEC_PREFIX)$(BUILD_DIR)/$(BINARY_NAME) -l; \
		echo "\nTesting help:"; \
		$(EXEC_PREFIX)$(BUILD_DIR)/$(BINARY_NAME) -h; \
	else \
		echo "No active Linux input method backend detected; skipping runtime IM smoke tests."; \
		echo "Use 'make test-wsl-win' in WSL when validating Windows IME flow."; \
		echo "\nTesting help:"; \
		$(EXEC_PREFIX)$(BUILD_DIR)/$(BINARY_NAME) -h; \
	fi
endif

run: build
ifeq ($(DETECTED_OS),Windows)
	$(EXEC_PREFIX)$(BUILD_DIR)$(PATH_SEP)$(BINARY_NAME).exe $(ARGS)
else
	$(EXEC_PREFIX)$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)
endif
