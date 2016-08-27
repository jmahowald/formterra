SOURCEDIR= .
SOURCES := $(shell find $(SOURCEDIR) -type f -name '*.go')
# ASSETS := $(shell find $())
BINARY=formterra

ASSETS=cmd/assets.go
ASSET_SOURCES := $(wildcard cmd/templates/*)

# H/T https://ariejan.net/2015/10/03/a-makefile-for-golang-cli-tools/
# VERSION=1.0.0
# BUILD_TIME=`date +%FT%T%z`

LDFLAGS=
#LDFLAGS=-ldflags "-X github.com/ariejan/roll/core.Version=${VERSION} -X github.com/ariejan/roll/core.BuildTime=${BUILD_TIME}"

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
	go generate -x ./cmd

build:
	go build

# TODO make this look into asssest directory
generate:
