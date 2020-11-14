TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=metal
GOPATH?=$(HOME)/go

default: build

build: fmtcheck
	go install

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test $(TEST) -v $(TESTARGS) -timeout=30s -parallel=10

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout=120m -parallel=10

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

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile docs-lint docs-lint-fix tfproviderlint tfproviderlint-fix tfproviderdocs-check

