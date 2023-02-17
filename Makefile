DACE_BIN := bin/daikin-ac-exporter
DACI_BIN := bin/daikin-ac-info

GO ?= go

VERSION := $(shell cat VERSION)
USE_VENDOR =
LOCAL_LDFLAGS = -buildmode=pie -ldflags "-X=main.Version=$(VERSION)"

.PHONY: all api build vendor
all: dep build

dep: ## Get the dependencies
	@$(GO) get -v -d ./...

update: ## Get and update the dependencies
	@$(GO) get -v -d -u ./...

tidy: ## Clean up dependencies
	@$(GO) mod tidy

vendor: dep ## Create vendor directory
	@$(GO) mod vendor

build: ## Build the binary files
	$(GO) build -v -o $(DACE_BIN) $(USE_VENDOR) $(LOCAL_LDFLAGS) ./cmd/daikin-ac-exporter
	$(GO) build -v -o $(DACI_BIN) $(USE_VENDOR) $(LOCAL_LDFLAGS) ./cmd/daikin-ac-info

clean: ## Remove previous builds
	@rm -f $(DACE_BIN) $(DACI_BIN)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: release
release: ## create release package from git
	git clone https://github.com/thkukuk/daikin-gomod
	mv daikin-gomod daikin-gomod-$(VERSION)
	sed -i -e 's|USE_VENDOR =|USE_VENDOR = -mod vendor|g' daikin-gomod-$(VERSION)/Makefile
	make -C daikin-gomod-$(VERSION) vendor
	cp VERSION daikin-gomod-$(VERSION)
	tar --exclude .git -cJf daikin-gomod-$(VERSION).tar.xz daikin-gomod-$(VERSION)
	rm -rf daikin-gomod-$(VERSION)
