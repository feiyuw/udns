# disable all default rules
.SUFFIXES:
MAKEFLAGS+=-r

.PHONY: all clean build test
.DEFAULT: all

# GO shit
ROOT_DIR:=$(realpath $(CURDIR))
export GOPATH:=$(ROOT_DIR)

# version & build time
VERSION:=$(shell git describe --dirty --tags)
ifeq (,$(VERSION))
VERSION:="master"
endif
TARGET:="udns"

all: clean build

clean:
	@echo cleaning...
	@rm -rf $(TARGET)

test:
	@echo unit testing...
	go test udns/...

build:
	@echo building...
	cd src/udns; go build -o ../../$(TARGET) -ldflags "-X main.Version=$(VERSION)"
