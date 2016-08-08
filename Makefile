# This makefile probably requires GNU make >= 3.81

GO := go

packages := . ./cmd/...

build:
	$(GO) build $(packages)

install:
	$(GO) install $(packages)

# Re: ginkgo, https://github.com/onsi/ginkgo/issues/278
test:
	$(GO) test $(packages)

spec:
	ginkgo -r -v

# TODO: coverage for CLI? https://github.com/onsi/ginkgo/issues/89
coverage:
	ginkgo --cover $(packages) --covermode count
	$(GO) tool cover --func connect.coverprofile

browse-coverage: coverage
	$(GO) tool cover --html connect.coverprofile

# golint only takes one package or else it thinks multiple arguments are
# directories (which it also doesn't support), ./... includes vendor :-/
lint:
	$(foreach pkg, $(packages), golint --set_exit_status $(pkg);)

# TODO: add to CI after dropping 1.5 support
# https://github.com/kisielk/errcheck/issues/75
errcheck:
	errcheck --asserts --ignore 'io:Close' $(packages)

zen:
	ginkgo watch -notify $(packages)

get-devtools:
	@echo Getting golint...
	$(GO) get -u github.com/golang/lint/golint
	@echo Getting errcheck...
	$(GO) get -u github.com/kisielk/errcheck

clean:
	$(RM) *.coverprofile
	$(RM) -r man

# In case you forget -s -v when using `glide get`.
clean-vendor:
	glide install --strip-vcs --strip-vendor

distclean: clean
	$(GO) clean -i github.com/go-kafka/connect...

.PHONY: build install test spec coverage browse-coverage
.PHONY: lint errcheck zen get-devtools
.PHONY: clean clean-vendor distclean
