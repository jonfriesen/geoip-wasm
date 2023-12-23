# Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOWASM=GOOS=js GOARCH=wasm $(GOBUILD)

# Build directories
WASMDIR=cmd/wasm
PUBLICDIR=./public

# WASM file name
WASMFILE=ip_location.wasm

# Get Go root directory
GOROOT:=$(shell go env GOROOT)

# GeoLite URL
GEOLITE_DB_URL := $(shell curl -s https://api.github.com/repos/P3TERX/GeoLite.mmdb/releases/latest | jq -r '.assets[] | select(.name == "GeoLite2-City.mmdb") | .browser_download_url')

all: download-geolite build-wasm copy-wasm-exec

build-wasm:
	$(GOWASM) -o $(PUBLICDIR)/$(WASMFILE) $(WASMDIR)/main.go

copy-wasm-exec:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(PUBLICDIR)/

clean:
	$(GOCLEAN)
	rm -f $(PUBLICDIR)/$(WASMFILE)
	rm -f $(PUBLICDIR)/wasm_exec.js
	rm -f ./assets/db/GeoLite2-City.mmdb

download-geolite:
	@echo "Downloading GeoLite2-City.mmdb..."
	@curl -L $(GEOLITE_DB_URL) -o ./assets/db/GeoLite2-City.mmdb

run-web-example:
	$(GORUN) ./cmd/web-example	
