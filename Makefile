# Edito Makefile

.PHONY: build clean test install config-tool example-config example-plugin

# Build the main edito binary
build:
	go build -o edito ./main.go

# Build the configuration compiler tool
config-tool:
	go build -o edito-config ./cmd/edito-config/

# Build example configuration
example-config: config-tool
	cd example-config/.config/edito && \
	../../../edito-config config.go

# Build example plugin
example-plugin:
	cd example-plugin/go-mode && \
	go build -buildmode=plugin -o go-mode.so go-mode.go

# Build everything
all: build config-tool example-config example-plugin

# Clean build artifacts
clean:
	rm -f edito edito-config
	rm -f example-config/.config/edito/*.so
	rm -f example-plugin/*/*.so

# Run tests
test:
	go test ./...

# Install edito system-wide
install: build config-tool
	sudo cp edito /usr/local/bin/
	sudo cp edito-config /usr/local/bin/

# Development setup
dev-setup:
	mkdir -p ~/.config/edito
	mkdir -p ~/.local/share/edito/plugins
	mkdir -p ~/.cache/edito

# Copy example configuration to user directory
install-example-config: example-config dev-setup
	cp example-config/.config/edito/config.go ~/.config/edito/config.go
	cd ~/.config/edito && edito-config config.go

# Copy example plugin to user directory
install-example-plugin: example-plugin dev-setup
	cp example-plugin/go-mode/go-mode.so ~/.local/share/edito/plugins/

# Complete development setup
dev-install: all dev-setup install-example-config install-example-plugin
	@echo "Development setup complete!"
	@echo "Configuration: ~/.config/edito/"
	@echo "Plugins: ~/.local/share/edito/plugins/"
	@echo "Run 'edito filename.txt' to start editing"

# Help
help:
	@echo "Available targets:"
	@echo "  build            - Build main edito binary"
	@echo "  config-tool      - Build configuration compiler"
	@echo "  example-config   - Build example configuration"
	@echo "  example-plugin   - Build example plugin"
	@echo "  all              - Build everything"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run tests"
	@echo "  install          - Install system-wide"
	@echo "  dev-setup        - Create user directories"
	@echo "  dev-install      - Complete development setup"
	@echo "  help             - Show this help"