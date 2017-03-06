SOURCEDIR= .
SOURCES := $(shell find $(SOURCEDIR) -type f -name '*.go')
# ASSETS := $(shell find $())
BINARY=formterra
LINUXBINARY=formterralinux
GONAME=github.com/backpackhealth/formterra
ASSETS=tfproject/assets.go
ASSET_SOURCES := $(shell find tfproject/assets/* -type f -print)
#ASSET_SOURCES := $(shell find .$(PROJ_GO_SRC)/tfproject/assets/* -type f -print)

# This allows us to test our packages without testing vendors.
# H/T https://coderwall.com/p/urusna/test-all-go-packages-in-a-project-without-testing-vendor
BASE=$(shell echo $PWD | sed "s|$GOPATH/src/||")
GO_TEST_PACKAGES=$(shell go list $(PROJ_GO_SRC)/... | grep -v vendor | sed "s|$(BASE)/|./|" )


DOC_SOURCES := $(shell find cmd -type f -name '*.go')
DOCKER_TAG ?= $(BINARY)

export PROJ_GO_SRC ?=go/src/github.com/backpackhealth/formterra/
# H/T https://ariejan.net/2015/10/03/a-makefile-for-golang-cli-tools/
# VERSION=1.0.0
# BUILD_TIME=`date +%FT%T%z`
VERSION ?= 0.2
BUILD_TIME=$(shell date +%FT%T%z)


#VERSION_FLAGS=-X $(GONAME)/core.Version=${VERSION} -X $(GONAME)/core.BuildTime=${BUILD_TIME}
LDFLAGS=-ldflags "-X $(GONAME)/core.Version=${VERSION} -X $(GONAME)/core.BuildTime=${BUILD_TIME}"
.DEFAULT_GOAL: $(BINARY)

.PHONY: install clean

$(BINARY): $(ASSETS) $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} main.go

$(ASSETS): $(ASSET_SOURCES)
	echo "assets are $(ASSET_SOURCES)"
	go generate -x ./tfproject

install:
	go install ${LDFLAGS} .


clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${LINUXBINARY} ] ; then rm ${LINUXBINARY} ; fi

docs: $(DOC_SOURCES)
	cd build; go run build.go ../docs

test: $(BINARY)

	cd tfproject; go test $(LDFLAGS) .




# This section is used to create a relatively minimal docker image 
# Thx - https://developer.atlassian.com/blog/2015/07/osx-static-golang-binaries-with-docker/


$(LINUXBINARY): 
	$(MAKE) buildlinux

buildtest:
	echo go test ${LDFLAGS} $(GO_TEST_PACKAGES)

# TODO figure out how to get build information when building linux in docker
buildgo: 
	CGO_ENABLED=0 GOOS=linux go build -ldflags  "-s" -a -installsuffix cgo -o $(LINUXBINARY) ./$(PROJ_GO_SRC)

buildlinux: 
	docker build -t build-$(BINARY) -f ./Dockerfile.build .
	docker run -t build-$(BINARY) /bin/true
	docker cp `docker ps -q -n=1`:/$(LINUXBINARY) .
	chmod 755 ./$(LINUXBINARY)

builddocker:  $(LINUXBINARY)
	docker build --rm=true --tag=$(DOCKER_TAG) -f Dockerfile.static .
