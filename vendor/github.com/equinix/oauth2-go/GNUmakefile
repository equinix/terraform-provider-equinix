GOCMD=go
TEST?=$$(go list ./... |grep -v 'vendor')

default: clean build test

all: default

test:
	${GOCMD} test -v

build:
	${GOCMD} build

clean:
	${GOCMD} clean

.PHONY: test build clean