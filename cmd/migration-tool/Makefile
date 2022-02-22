
BINARY =equinix-migration-tool
GOCMD  =go
TEST   ?=$$(go list ./... |grep -v 'vendor')

default: clean build test

all: default

test:
	echo $(TEST) | \
		xargs -t ${GOCMD} test -v -timeout=10m

clean:
	${GOCMD} clean
	rm -f ${BINARY}

build:
	${GOCMD} build -o ${BINARY}

.PHONY: build clean release
