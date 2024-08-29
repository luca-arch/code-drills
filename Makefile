SHELL:=/bin/bash
GCI_LINT:=v1.60.3


help: ### Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help


lint-assets: ### Run eslint and prettier
	cd ./frontend-app && npm run lint-fix;
.PHONY: lint-assets


lint-go: ### Run go fmt and golangci-lint
	go fmt ./...;
	echo -e '\033[1mgo fmt finished!\033[0m';
	docker run --rm -v $(PWD):/mnt -w /mnt golangci/golangci-lint:$(GCI_LINT) golangci-lint run --fix -v
	echo -e '\033[1mgolangci-lint passed!\033[0m';
	docker run --rm -v $(PWD):/data cytopia/goimports:latest -w .
	echo -e '\033[1mgoimports passed!\033[0m';
.PHONY: lint-go


tests-assets: ### Run unit test (frontend app)
	cd ./frontend-app && npm run test;
.PHONY: tests-assets


tests-go: ### Run unit test (backend app)
	docker run --rm -ti \
		-w /mnt \
		-v $(PWD):/mnt golang:1.23.0-alpine3.20 \
		go test ./...
.PHONY: tests-go