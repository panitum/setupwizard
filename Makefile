APP_NAME ?= setupwizard
BUILD_DIR ?= bin

.PHONY: all build-all windows windows-amd64 windows-arm64 darwin darwin-arm64 darwin-amd64 linux linux-amd64 linux-arm64 clean help

all: build-all

help:
	@echo "Доступные команды:"
	@echo "  make build-all      - Собрать под все платформы (Windows, macOS ARM/Intel, Linux)"
	@echo "  make windows        - Собрать под Windows (amd64, arm64)"
	@echo "  make darwin         - Собрать под macOS (Apple Silicon и Intel)"
	@echo "  make darwin-arm64   - Собрать под macOS ARM64 (Apple Silicon)"
	@echo "  make darwin-amd64   - Собрать под macOS AMD64 (Intel)"
	@echo "  make linux          - Собрать под Linux (amd64, arm64)"
	@echo "  make clean          - Очистить директорию с собранными бинарниками"

build-all: windows darwin linux

# Windows
windows: windows-amd64 windows-arm64

windows-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .

windows-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe .

# macOS (Darwin)
darwin: darwin-arm64 darwin-amd64

darwin-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .

darwin-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .

# Linux
linux: linux-amd64 linux-arm64

linux-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .

linux-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .

# Clean
clean:
	rm -rf $(BUILD_DIR)
