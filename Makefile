.PHONY: default build clean
APPS        := bot
BLDDIR      ?= bin
VERSION     ?= $(shell cat VERSION)
IMPORT_BASE := github.com/opencars/bot
LDFLAGS     := -ldflags "-X $(IMPORT_BASE)/pkg/version.Version=$(VERSION)"

default: clean build

build: $(APPS)

$(BLDDIR)/%:
	go build $(LDFLAGS) -o $@ ./cmd/$*

$(APPS): %: $(BLDDIR)/%

lint:
	@revive -formatter stylish -config=revive.toml ./...

clean:
	@mkdir -p $(BLDDIR)
	@for app in $(APPS) ; do \
		rm -f $(BLDDIR)/$$app ; \
	done
