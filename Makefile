MODULE_NAME=aferocopy

VENDOR_DIR = vendor

GOLANGCI_LINT_VERSION ?= v2.13.0

GO ?= go
GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint-$(GOLANGCI_LINT_VERSION)

GITHUB_OUTPUT ?= /dev/null

# Other config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

ifeq ($(V),1)
  Q = @set -x;
else
  Q = @
endif

.PHONY: $(VENDOR_DIR)
$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor
	@$(GO) mod tidy

.PHONY: bump-deps
bump-deps:
	$(Q)$(GO) get -u ./...

.PHONY: tidy
tidy:
	$(Q)$(GO) mod tidy

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@printf -- "$(OK_COLOR)==> lint$(NO_COLOR)\n"
	@$(GOLANGCI_LINT) run -c .golangci.yaml

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	@printf -- "$(OK_COLOR)==> unit test$(NO_COLOR)\n"
	$(Q)$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./... -tags testcoverage

#.PHONY: test-integration
#test-integration:
#	@printf -- "$(OK_COLOR)==> integration test$(NO_COLOR)\n"
#	$(Q)$(GO) test ./features/... -gcflags=-l -coverprofile=features.coverprofile -coverpkg ./... -race --godog

.PHONY: $(GITHUB_OUTPUT)
$(GITHUB_OUTPUT):
	@echo "MODULE_NAME=$(MODULE_NAME)" >> "$@"
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)" >> "$@"

$(GOLANGCI_LINT):
	@printf -- "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)\n"
	$(Q)curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b /tmp "$(GOLANGCI_LINT_VERSION)"
	$(Q)$(call install-dep,/tmp/golangci-lint,$(GOLANGCI_LINT))

define install-dep
	if [ "$(1)" != "$(2)" ]; then \
		mkdir -p $$(dirname $(2)) || true; \
		mv $(1) $(2); \
	fi

	chmod +x "$(2)"
endef
