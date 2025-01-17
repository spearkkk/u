.PHONY: build install clean

WORKFLOW_DIR=$(HOME)/Library/Application Support/Alfred/Alfred.alfredpreferences/workflows/utilities

build:
	GOOS=darwin GOARCH=amd64 go build -o workflow/utilities *.go
	chmod +x workflow/utilities

package: build
	cd workflow && zip -r ../utilities.alfredworkflow *

install: package
	mkdir -p "$(WORKFLOW_DIR)"
	cp utilities.alfredworkflow "$(WORKFLOW_DIR)"

clean:
	rm -f workflow/utilities
	rm -f utilities.alfredworkflow

all: clean build package install 