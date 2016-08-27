SOURCEDIR= .
SOURCES := $(shell find $(SOURCEDIR) -type f -name '*.go')
# ASSETS := $(shell find $())
BINARY=formterra

ASSETS=assets.go
ASSET_SOURCES := $(wildcard tfproject/assets/*)

# H/T https://ariejan.net/2015/10/03/a-makefile-for-golang-cli-tools/
# VERSION=1.0.0
# BUILD_TIME=`date +%FT%T%z`
VERSION ?= 0.1
BUILD_TIME=$(shell date +%FT%T%z)

LDFLAGS=-ldflags "-X github.com/jmahowald/formterra/core.Version=${VERSION} -X github.com/jmahowald/formterra/core.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(ASSETS) $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} main.go


.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

$(ASSETS): $(ASSET_SOURCES)
	go generate -x ./tfproject

build:
	go build


# TODO make this look into asssest directory
