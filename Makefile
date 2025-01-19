.PHONY: build clean package deploy all help

PACKAGE_NAME := u.alfredworkflow
BINARY_NAME := workflow/u
WORKFLOW_DIR := workflow

# Default target
default: package

# Build target
build:
	@echo "Building binary for Alfred workflow..."
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME) *.go
	@chmod +x $(BINARY_NAME)

# Clean target
clean:
	@echo "Cleaning up build artifacts..."
	rm -f $(PACKAGE_NAME)
	rm -rf $(WORKFLOW_DIR)

# Package target
package: build
	@echo "Packaging Alfred workflow..."
	cp ./info.plist $(WORKFLOW_DIR)
	cp ./icon.png $(WORKFLOW_DIR)
	cd $(WORKFLOW_DIR) && zip -r ../$(PACKAGE_NAME) *

# Deploy target
deploy:
	@echo "Deploying Alfred workflow..."
	@echo "Fetching git tags..."

	git fetch --tags

	@echo "Determining the latest tag..."

	LATEST_TAG=$$(git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null || echo ""); \
	if [ -z "$$LATEST_TAG" ]; then \
	    echo "No existing tags. Setting the initial tag to v0.0.1"; \
	    NEW_TAG=v0.0.1; \
	else \
	    echo "Latest tag: $$LATEST_TAG"; \
	    NEW_TAG=$$(echo $$LATEST_TAG | awk -F. -v OFS=. 'NF==1 {print $$1".0.1"} NF==2 {print $$1"."$$2".1"} NF>=3 {$$NF++; print}'); \
	fi; \
	if [ -n "$(TAG)" ]; then NEW_TAG=$(TAG); fi; \
	echo "New tag: $$NEW_TAG"

	@echo "Creating and pushing new tag..."
	git tag $$NEW_TAG; \
	git push origin $$NEW_TAG

# All target
all: clean build package

# Help target
help:
	@echo "Available targets:"
	@echo "  build       - Build the Alfred workflow binary"
	@echo "  clean       - Remove build artifacts"
	@echo "  package     - Package the workflow into a .alfredworkflow file"
	@echo "  deploy      - Deploy a new release with a git tag"
	@echo "                Optionally pass TAG=<new_tag> to set a custom version"
	@echo "  all         - Clean, build, and package"