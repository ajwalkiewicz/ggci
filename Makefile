APP_NAME := ggci
CMD_PATH := ./cmd/ggci
BIN_DIR := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

VERSION := $(shell git describe --tags --always --dirty)

.PHONY: build clean version

build:
	mkdir -p $(BIN_DIR)
	go build \
		-ldflags "-X github.com/ajwalkiewicz/ggci/internal/app.version=$(VERSION)" \
		-o $(BIN_PATH) \
		$(CMD_PATH)
	@echo "Built $(APP_NAME) version $(VERSION) at $(BIN_PATH)"

version:
	@echo $(VERSION)

clean:
	rm -rf $(BIN_DIR)