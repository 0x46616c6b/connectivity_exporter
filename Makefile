BUILD_DIR = build

GO       = go
GOX      = gox
GOX_ARGS = -output="$(BUILD_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="linux/amd64 linux/386 linux/arm linux/arm64 darwin/amd64 freebsd/amd64 freebsd/386 windows/386 windows/amd64"

.PHONY: build
build:
	$(GO) build -o $(BUILD_DIR)/connectivity_exporter .

.PHONY: clean
clean:
	rm -R $(BUILD_DIR)/* || true

.PHONY: test
test:
	$(GO) test ./...

.PHONY: release-build
release-build:
	@go get -u github.com/mitchellh/gox
	@$(GOX) $(GOX_ARGS) github.com/0x46616c6b/connectivity_exporter