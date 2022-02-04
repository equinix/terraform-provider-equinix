GOCMD				=go
TEST				?=$$(go list ./... |grep -v 'vendor')
INSTALL_DIR			=~/.terraform.d/plugins
BINARY				=terraform-provider-equinix
SWEEP				?=all #Flag required to define the regions that the sweeper is to be ran in
SWEEP_DIR			?=./equinix
SWEEP_ARGS			=-timeout 60m
ACCTEST_TIMEOUT		?= 180m
ACCTEST_PARALLELISM ?= 8
ACCTEST_COUNT       ?= 1
GOFMT_FILES         ?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME            =equinix

ifneq ($(origin TESTS_REGEXP), undefined)
	TESTARGS = -run='$(TESTS_REGEXP)'
endif

default: clean build test

all: default
	
test:
	echo $(TEST) | \
		xargs -t ${GOCMD} test -v $(TESTARGS) -timeout=10m

testacc:
	TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 ${GOCMD} test $(TEST) -v -count $(ACCTEST_COUNT) -parallel $(ACCTEST_PARALLELISM) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT)

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep=$(SWEEP) $(SWEEP_ARGS)

build:
	${GOCMD} build -o ${BINARY}

install: test build
	@if [ -d ${INSTALL_DIR} ]; then \
		echo "==> [INFO] installing in ${INSTALL_DIR} directory"; \
		cp ${BINARY} ${INSTALL_DIR}; \
	else \
		echo "==> [ERROR] installation plugin directory ${INSTALL_DIR} does not exist"; \
	fi

clean:
	${GOCMD} clean
	rm -f ${BINARY}


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


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

docs-lint:
	@echo "==> Checking docs against linters..."
	@misspell -error -source=text docs/ || (echo; \
		echo "Unexpected misspelling found in docs files."; \
		echo "To automatically fix the misspelling, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli docs/ || (echo; \
		echo "Unexpected issues found in docs Markdown files."; \
		echo "To apply any automatic fixes, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)

docs-lint-fix:
	@echo "==> Applying automatic docs linter fixes..."
	@misspell -w -source=text docs/
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli --fix docs/

tfproviderlint:
	@echo "==> Checking provider code against bflad/tfproviderlint..."
	@docker run -v $(PWD):/src bflad/tfproviderlint ./... || (echo; \
		echo "Unexpected issues found in code with bflad/tfproviderlint."; \
		echo "To apply automated fixes for check that support them, run 'make tfproviderlint-fix'."; \
		exit 1)

tfproviderlint-fix:
	@echo "==> Applying fixes with bflad/tfproviderlint..."
	@docker run -v $(PWD):/src bflad/tfproviderlint -fix ./...

tfproviderdocs-check:
	@echo "==> Check provider docs with bflad/tfproviderdocs..."
	@docker run -v $(PWD):/src bflad/tfproviderdocs check -provider-name=metal || (echo; \
		echo "Unexpected issues found in code with bflad/tfproviderdocs."; \
		exit 1)

.PHONY: test testacc build install clean  fmt fmtcheck errcheck test-compile docs-lint docs-lint-fix tfproviderlint tfproviderlint-fix tfproviderdocs-check
