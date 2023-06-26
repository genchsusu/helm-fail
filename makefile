HELM_PLUGIN_NAME := fail
VERSION := $(shell cat plugin.yaml | grep "version" | cut -d '"' -f 2)
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: simulate build tag

simulate:
	helm upgrade -i simulate ./chart --wait

build:
	export CGO_ENABLED=0 && \
	go build -o bin/${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) .

tag:
	@scripts/tag.sh