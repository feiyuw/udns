# disable all default rules
.SUFFIXES:
MAKEFLAGS+=-r

.PHONY: all clean build test
.DEFAULT: all

# version & build time
VERSION:=$(shell git describe --dirty --tags)
ifeq (,$(VERSION))
VERSION:="master"
endif
TARGET:="udns"

all: clean build test

clean:
	@echo cleaning...
	@rm -rf ./build/$(TARGET)
	@rm -rf ./build/*.rpm

test:
	@echo unit testing...
	go test udns/...

build:
	@echo building...
	go build -o ../build/$(TARGET) -ldflags "-X main.Version=$(VERSION)"

rpm:
	@echo generate rpm...
	cd build; fpm -s dir -t rpm --prefix /opt/udns -n udns -v $(VERSION) --after-upgrade ./start.sh --after-install ./start.sh --after-remove ./remove.sh .
