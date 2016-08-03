# This makefile probably requires GNU make >= 3.81

GO := go

build:
	$(GO) build

install: build
	$(GO) install $$(glide novendor)

test: build
	ginkgo -v $$(glide novendor)

# TODO: tests and coverage for CLI
coverage: build
	ginkgo --cover $$(glide novendor)
	$(GO) tool cover --func connect.coverprofile

browse-coverage: coverage
	$(GO) tool cover --html connect.coverprofile

# golint only takes one package or else it thinks multiple arguments are
# directories (which it also doesn't support), so `glide novendor` won't work :-/
lint:
	golint --set_exit_status . && \
		golint --set_exit_status ./cmd/...

zen:
	ginkgo watch -notify $$(glide novendor)

get-devtools:
	@echo Getting golint...
	$(GO) get -u -v github.com/golang/lint/golint
	@echo Getting the Ginkgo test runner...
	$(GO) get -u -v github.com/onsi/ginkgo/ginkgo

clean:
	$(RM) *.coverprofile
	$(RM) -r man

# In case you forget -s -v when using `glide get`.
clean-vendor:
	glide install --strip-vcs --strip-vendor

distclean: clean
	$(GO) clean -i github.com/go-kafka/connect...

.PHONY: build install test coverage browse-coverage lint zen get-devtools
.PHONY: clean clean-vendor distclean
