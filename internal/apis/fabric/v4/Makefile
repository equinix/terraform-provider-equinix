.PHONY: all gen patch fetch

CURRENT_UID := $(shell id -u)
CURRENT_GID := $(shell id -g)

# https://github.com/OpenAPITools/openapi-generator-cli
SPEC_URL:=https://api.equinix.com/metal/v1/api-docs
SPEC_FETCHED_FILE:=equinix-metal.fetched.json
SPEC_PATCHED_FILE:=equinix-metal.patched.json
IMAGE=openapitools/openapi-generator-cli
GIT_ORG=t0mk
GIT_REPO=gometal
PACKAGE_MAJOR=v1

SWAGGER=docker run --rm -u ${CURRENT_UID}:${CURRENT_GID} -v $(CURDIR):/local ${IMAGE}
GOLANGCI_LINT=golangci-lint

all: pull fetch fix-tags patch clean gen mod test stage

pull:
	docker pull ${IMAGE}

fetch:
	curl -o ${SPEC_FETCHED_FILE} ${SPEC_URL}

fix-tags:
	- jq '. | select(((.paths[][].tags| type=="array"), length) > 1).paths[][].tags |= [.[0]]' ${SPEC_FETCHED_FILE} | diff -d -U6 ${SPEC_FETCHED_FILE} - >  patches/01-tag-from-last-in-path.patch

patch:
	# patch is idempotent, always starting with the fetched
	# fetched file to create the patched file.
	ARGS="-o ${SPEC_PATCHED_FILE} ${SPEC_FETCHED_FILE}"; \
	for diff in $(shell find patches -name \*.patch | sort -n); do \
		patch --no-backup-if-mismatch -N -t $$ARGS $$diff; \
		ARGS=${SPEC_PATCHED_FILE}; \
	done
	find ${SPEC_PATCHED_FILE} -empty -exec cp ${SPEC_FETCHED_FILE} ${SPEC_PATCHED_FILE} \;

clean:
	rm -rf v1/

gen:
	${SWAGGER} generate -g go \
		--package-name ${PACKAGE_MAJOR} \
		--model-package types \
		--api-package models \
		--git-user-id ${GIT_ORG} \
		--git-repo-id ${GIT_REPO} \
		-o /local/${PACKAGE_MAJOR} \
		-i /local/${SPEC_PATCHED_FILE}

validate:
	${SWAGGER} validate \
		--recommend \
		-i /local/${SPEC_PATCHED_FILE}

mod:
	cd v1 && go mod tidy

test:
	cd v1 && go test -v ./...

# https://github.com/OpenAPITools/openapi-generator/issues/741#issuecomment-569791780
remove-dupe-requests: ## Removes duplicate Request structs from the generated code
	@for struct in $$(grep -h 'type .\{1,\} struct' $(PACKAGE_MAJOR)/*.go | grep Request  | sort | uniq -c | grep -v '^      1' | awk '{print $$3}'); do \
	  for f in $$(/bin/ls $(PACKAGE_MAJOR)/*.go); do \
	    if grep -qF "type $${struct} struct" "$${f}"; then \
	      if eval "test -z \$${$${struct}}"; then \
	        echo "skipping first appearance of $${struct} in file $${f}"; \
	        eval "export $${struct}=1"; \
	      else \
	        echo "removing dupe $${struct} from file $${f}"; \
	        tr '\n' '\r' <"$${f}" | sed 's~// '"$${struct}"'.\{1,\}type '"$${struct}"' struct {[^}]\{1,\}}~~' | tr '\r' '\n' >"$${f}.tmp"; \
	        mv -f "$${f}.tmp" "$${f}"; \
	      fi; \
	    fi \
	  done \
	done
lint:
	@$(GOLANGCI_LINT) run -v --no-config --fast=false --fix --disable-all --enable goimports $(PACKAGE_MAJOR)

stage:
	git add --intent-to-add v1
