## Makefile

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get -u github.com/golang/dep/cmd/dep
	@$(MAKE) dep

.PHONY: dep
dep: ## run dep ensure and prune
	dep ensure --vendor-only

.PHONY: fmt
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: lint
lint: ## Run all the linters
	gometalinter --vendor --disable-all \
	--enable=deadcode \
	--enable=ineffassign \
	--enable=gosimple \
	--enable=staticcheck \
	--enable=gofmt \
	--enable=goimports \
	--enable=misspell \
	--enable=errcheck \
	--enable=vet \
	--enable=vetshadow \
	--deadline=10m \
	./...

.PHONY: build
build: kv client server


.PHONY: kv
kv:
	CGO_ENABLED=0 GOOS=darwin go build -o ./bin/kv ./kv/main.go

.PHONY: client
client:
	CGO_ENABLED=0 GOOS=darwin go build -o ./bin/client ./echo_client/main.go

.PHONY: server
server:
	CGO_ENABLED=0 GOOS=darwin go build -o ./bin/server ./echo_server/main.go

.PHONY:clean
clean: ## Remove temporary files
	go clean

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build

