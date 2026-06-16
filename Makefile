.PHONY: build build-all build-windows build-linux build-darwin clean run

APP_NAME    := gh-proxy-go
DIST_DIR    := dist
GO          := go
LD_FLAGS    := -s -w

build: build-linux build-windows

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="$(LD_FLAGS)" -o $(DIST_DIR)/$(APP_NAME)_windows_amd64.exe .

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LD_FLAGS)" -o $(DIST_DIR)/$(APP_NAME)_linux_amd64 .

build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags="$(LD_FLAGS)" -o $(DIST_DIR)/$(APP_NAME)_linux_arm64 .

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="$(LD_FLAGS)" -o $(DIST_DIR)/$(APP_NAME)_darwin_amd64 .

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="$(LD_FLAGS)" -o $(DIST_DIR)/$(APP_NAME)_darwin_arm64 .

build-all: build-windows build-linux build-linux-arm64 build-darwin build-darwin-arm64

clean:
	rm -rf $(DIST_DIR)/*

run:
	$(GO) run .
