TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
VERSION?=0.6.5
TF_PLUGINS_DIR?=$(HOME)/.terraform.d/plugins/$$(go env GOOS)_$$(go env GOARCH)
TARGETS=darwin_amd64 darwin_arm64 linux_amd64 windows_amd64
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=circleci
SWEEP_RUN=circleci_project
SWEEP_REGION=us

default: build

build: fmtcheck
	go build -o terraform-provider-circleci_v$(VERSION) main.go

install: build
	mkdir -p $(TF_PLUGINS_DIR)
	mv terraform-provider-circleci_v$(VERSION) $(TF_PLUGINS_DIR)

targets: $(TARGETS)

$(TARGETS):
	GOOS=$(firstword $(subst _, ,$@)) GOARCH=$(lastword $(subst _, ,$@)) CGO_ENABLED=0 go build -o "dist/terraform-provider-circleci_$$(git describe --tags)_$@"
	zip -j dist/terraform-provider-circleci_$$(git describe --tags)_$@.zip dist/terraform-provider-circleci_$$(git describe --tags)_$@

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

sweep:
	go test $(TEST) -v -sweep=$(SWEEP_REGION) -sweep-run=$(SWEEP_RUN)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile website website-test
