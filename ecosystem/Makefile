# AdvanceGG Ecosystem Build System
# Builds standalone packages for multiple languages without requiring Go compiler for end users

.PHONY: all clean build-all build-python build-nodejs build-native-libs test package

# Configuration
VERSION := 1.0.0
PLATFORMS := windows-x64 windows-arm64 darwin-x64 darwin-arm64 linux-x64 linux-arm64 linux-armv7
GO_VERSION := 1.21
PYTHON_VERSIONS := 3.7 3.8 3.9 3.10 3.11 3.12

# Directories
BUILD_DIR := build
DIST_DIR := dist
NATIVE_DIR := native
PYTHON_DIR := python
NODEJS_DIR := nodejs

# Default target
all: build-all package

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	rm -rf $(PYTHON_DIR)/build $(PYTHON_DIR)/dist $(PYTHON_DIR)/*.egg-info
	rm -rf $(NODEJS_DIR)/native $(NODEJS_DIR)/node_modules
	@echo "✅ Clean completed"

# Build all packages
build-all: build-native-libs build-python build-nodejs
	@echo "🎉 All packages built successfully!"

# Build native libraries for all platforms
build-native-libs:
	@echo "🔨 Building native libraries for all platforms..."
	mkdir -p $(BUILD_DIR)/$(NATIVE_DIR)
	
	# Build for each platform
	@for platform in $(PLATFORMS); do \
		echo "📦 Building for $$platform..."; \
		$(MAKE) build-native-single PLATFORM=$$platform; \
	done
	
	@echo "✅ Native libraries built for all platforms"

# Build native library for single platform
build-native-single:
	@if [ -z "$(PLATFORM)" ]; then \
		echo "❌ PLATFORM not specified"; \
		exit 1; \
	fi
	
	@echo "🔨 Building native library for $(PLATFORM)..."
	
	# Parse platform
	$(eval OS := $(word 1,$(subst -, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst -, ,$(PLATFORM))))
	
	# Set Go environment
	$(eval GOOS := $(OS))
	$(eval GOARCH := $(if $(filter $(ARCH),x64),amd64,$(ARCH)))
	$(eval EXT := $(if $(filter $(OS),windows),.dll,$(if $(filter $(OS),darwin),.dylib,.so)))
	
	# Build command
	cd .. && GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-buildmode=c-shared \
		-ldflags="-s -w" \
		-o ecosystem/$(BUILD_DIR)/$(NATIVE_DIR)/advancegg-$(PLATFORM)$(EXT) \
		ecosystem/$(PYTHON_DIR)/advancegg_python.go

# Build Python package
build-python: build-native-libs
	@echo "🐍 Building Python package..."
	
	# Copy native libraries to Python package
	mkdir -p $(PYTHON_DIR)/advancegg/native
	cp $(BUILD_DIR)/$(NATIVE_DIR)/* $(PYTHON_DIR)/advancegg/native/
	
	# Build Python package
	cd $(PYTHON_DIR) && python setup.py sdist bdist_wheel
	
	# Copy to dist directory
	mkdir -p $(DIST_DIR)/python
	cp $(PYTHON_DIR)/dist/* $(DIST_DIR)/python/
	
	@echo "✅ Python package built"

# Build Node.js package
build-nodejs: build-native-libs
	@echo "🟨 Building Node.js package..."
	
	# Copy native libraries to Node.js package
	mkdir -p $(NODEJS_DIR)/native
	cp $(BUILD_DIR)/$(NATIVE_DIR)/* $(NODEJS_DIR)/native/
	
	# Install dependencies and build
	cd $(NODEJS_DIR) && npm install
	cd $(NODEJS_DIR) && npm pack
	
	# Copy to dist directory
	mkdir -p $(DIST_DIR)/nodejs
	cp $(NODEJS_DIR)/*.tgz $(DIST_DIR)/nodejs/
	
	@echo "✅ Node.js package built"

# Test packages
test: test-python test-nodejs
	@echo "✅ All tests passed"

# Test Python package
test-python:
	@echo "🧪 Testing Python package..."
	cd $(PYTHON_DIR) && python -m pytest tests/ -v
	cd $(PYTHON_DIR) && python -c "import advancegg; print('✅ Python import test passed')"

# Test Node.js package
test-nodejs:
	@echo "🧪 Testing Node.js package..."
	cd $(NODEJS_DIR) && npm test
	cd $(NODEJS_DIR) && node -e "const AdvanceGG = require('./index.js'); console.log('✅ Node.js import test passed');"

# Package for distribution
package: build-all
	@echo "📦 Creating distribution packages..."
	
	# Create release archives
	mkdir -p $(DIST_DIR)/releases
	
	# Python packages (already built)
	@echo "📦 Python packages ready in $(DIST_DIR)/python/"
	
	# Node.js packages (already built)
	@echo "📦 Node.js packages ready in $(DIST_DIR)/nodejs/"
	
	# Native libraries archive
	cd $(BUILD_DIR) && tar -czf ../$(DIST_DIR)/releases/advancegg-native-libs-$(VERSION).tar.gz $(NATIVE_DIR)/
	
	# Source archive
	cd .. && tar -czf ecosystem/$(DIST_DIR)/releases/advancegg-ecosystem-$(VERSION).tar.gz \
		--exclude='ecosystem/build' \
		--exclude='ecosystem/dist' \
		--exclude='ecosystem/*/node_modules' \
		--exclude='ecosystem/*/build' \
		--exclude='ecosystem/*/dist' \
		ecosystem/
	
	@echo "✅ Distribution packages created in $(DIST_DIR)/"

# Publish packages
publish: publish-python publish-nodejs
	@echo "🚀 All packages published!"

# Publish Python package to PyPI
publish-python:
	@echo "🐍 Publishing Python package to PyPI..."
	cd $(PYTHON_DIR) && twine upload dist/*
	@echo "✅ Python package published"

# Publish Node.js package to npm
publish-nodejs:
	@echo "🟨 Publishing Node.js package to npm..."
	cd $(NODEJS_DIR) && npm publish
	@echo "✅ Node.js package published"

# Development setup
dev-setup:
	@echo "🛠️ Setting up development environment..."
	
	# Check Go installation
	@go version || (echo "❌ Go not installed. Please install Go $(GO_VERSION)+"; exit 1)
	
	# Check Python installation
	@python3 --version || (echo "❌ Python not installed"; exit 1)
	
	# Check Node.js installation
	@node --version || (echo "❌ Node.js not installed"; exit 1)
	
	# Install Python development dependencies
	cd $(PYTHON_DIR) && pip install -e .[dev]
	
	# Install Node.js development dependencies
	cd $(NODEJS_DIR) && npm install
	
	@echo "✅ Development environment ready"

# Install packages locally for testing
install-local: build-all
	@echo "📥 Installing packages locally for testing..."
	
	# Install Python package
	cd $(PYTHON_DIR) && pip install -e .
	
	# Install Node.js package globally
	cd $(NODEJS_DIR) && npm install -g .
	
	@echo "✅ Packages installed locally"

# Uninstall local packages
uninstall-local:
	@echo "🗑️ Uninstalling local packages..."
	pip uninstall -y advancegg || true
	npm uninstall -g advancegg || true
	@echo "✅ Local packages uninstalled"

# Check package integrity
check: build-all
	@echo "🔍 Checking package integrity..."
	
	# Check Python package
	cd $(PYTHON_DIR) && python setup.py check
	cd $(PYTHON_DIR) && twine check dist/*
	
	# Check Node.js package
	cd $(NODEJS_DIR) && npm audit
	cd $(NODEJS_DIR) && npm run lint || true
	
	@echo "✅ Package integrity check completed"

# Generate documentation
docs:
	@echo "📚 Generating documentation..."
	
	# Python docs
	cd $(PYTHON_DIR) && python -m sphinx docs/ docs/_build/
	
	# Node.js docs
	cd $(NODEJS_DIR) && npm run docs || true
	
	@echo "✅ Documentation generated"

# Show help
help:
	@echo "AdvanceGG Ecosystem Build System"
	@echo "================================"
	@echo ""
	@echo "Targets:"
	@echo "  all              - Build all packages (default)"
	@echo "  clean            - Clean build artifacts"
	@echo "  build-all        - Build all packages"
	@echo "  build-python     - Build Python package"
	@echo "  build-nodejs     - Build Node.js package"
	@echo "  build-native-libs - Build native libraries for all platforms"
	@echo "  test             - Run all tests"
	@echo "  package          - Create distribution packages"
	@echo "  publish          - Publish packages to repositories"
	@echo "  dev-setup        - Setup development environment"
	@echo "  install-local    - Install packages locally for testing"
	@echo "  uninstall-local  - Uninstall local packages"
	@echo "  check            - Check package integrity"
	@echo "  docs             - Generate documentation"
	@echo "  help             - Show this help"
	@echo ""
	@echo "Environment Variables:"
	@echo "  VERSION          - Package version (default: $(VERSION))"
	@echo "  PLATFORMS        - Target platforms (default: $(PLATFORMS))"
	@echo ""
	@echo "Examples:"
	@echo "  make build-all                    # Build all packages"
	@echo "  make build-native-single PLATFORM=linux-x64  # Build for specific platform"
	@echo "  make test                         # Run tests"
	@echo "  make package                      # Create distribution packages"
