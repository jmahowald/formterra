SOURCEDIR= .
SOURCES := $(shell find $(SOURCEDIR) -type f -name '*.go')
# ASSETS := $(shell find $())
BINARY=formterra
LINUXBINARY=formterralinux
GONAME=gihub.com/jmahowald/formterra
ASSETS=tfproject/assets.go
ASSET_SOURCES := $(shell find tfproject/assets/* -type f -print)

DOC_SOURCES := $(shell find cmd -type f -name '*.go')
DOCKER_TAG ?= $(BINARY)

# H/T https://ariejan.net/2015/10/03/a-makefile-for-golang-cli-tools/
# VERSION=1.0.0
# BUILD_TIME=`date +%FT%T%z`
VERSION ?= 0.1
BUILD_TIME=$(shell date +%FT%T%z)

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

vendor:
	glide install



# This section is used to create a relatively minimal docker image 
# Thx - https://developer.atlassian.com/blog/2015/07/osx-static-golang-binaries-with-docker/


$(LINUXBINARY): vendor
	$(MAKE) buildlinux

buildgo: 
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o $(LINUXBINARY) ./go/src/github.com/jmahowald/formterra

buildlinux: vendor
	docker build -t build-$(BINARY) -f ./Dockerfile.build .
	docker run -t build-$(BINARY) /bin/true
	docker cp `docker ps -q -n=1`:/$(LINUXBINARY) .
	chmod 755 ./$(LINUXBINARY)

builddocker:  $(LINUXBINARY)
	docker build --rm=true --tag=$(DOCKER_TAG) -f Dockerfile.static .
