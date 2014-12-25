GO = /usr/bin/go
BUILD_DIR = build
PROJECT = github.com/opbk/openbook-landing

.PHONY: all build clear

all: build

dependencies = github.com/gorilla/mux \
github.com/gorilla/sessions \
github.com/gorilla/schema \
github.com/cihub/seelog \

dependencies_paths := $(addprefix $(GOPATH)/src/,$(dependencies))
$(dependencies_paths):
	for i in $(dependencies); do $(GO) get -d $$i; done

dependencies: $(dependencies_paths)

build:
	rm -rf build/
	$(GO) build -o $(BUILD_DIR)/landing $(PROJECT)
	cp -r resources build/
	touch build/users.csv

clear:
	rm -rf build/