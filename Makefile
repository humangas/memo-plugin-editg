PLUGIN_DIR=$(shell grep "pluginsdir.*" ~/.config/memo/config.toml | grep -o "\".*\"" | sed -e 's/"//g')

.DEFAULT_GOAL := help

.PHONY: all help setup deps install

all:

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "target:"
	@echo " - install:    install memo-plugin-editg"
	@echo " - deps:       dep ensure"
	@echo ""

setup:
	go get -u github.com/golang/dep/cmd/dep

deps: setup
	dep ensure

install: deps
	GOOS=darwin go build -o editg *.go
	mv editg $(PLUGIN_DIR)
