# This makefile probably requires GNU make >= 3.81

GO ?= go
# OS X, use sha256sum or gsha256sum elsewhere
SHASUM := shasum -a 256
VERSION := $(shell git describe)

packages := . ./cmd/...
fordist := find * -type d -exec

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
	$(GO) test --covermode count --coverprofile connect.coverprofile .
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

get-reltools:
	@echo Getting gox...
	$(GO) get -u github.com/mitchellh/gox

dist: test
	@echo Cross-compiling binaries...
	gox -verbose \
		-ldflags "-s -w" \
		-os="darwin linux windows" \
		-arch="amd64 386" \
		-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" ./cmd/...

release: dist
	@echo Preparing distributions...
	@cd dist && \
		$(fordist) sh -c 'gpg --detach-sign --armor {}/kafka-connect*' \; && \
		$(fordist) cp ../LICENSE {} \; && \
		$(fordist) cp ../README.md {} \; && \
		$(fordist) cp ../HISTORY.md {} \; && \
		$(fordist) tar -zcf kafka-connect-${VERSION}-{}.tar.gz {} \; && \
		$(fordist) zip -r kafka-connect-${VERSION}-{}.zip {} \; && \
		echo Computing checksums... && \
		find . \( -name '*.tar.gz' -or -name '*.zip' \) -exec \
			sh -c '$(SHASUM) {} > {}.sha256sum' \; && \
	cd ..
	@echo Done

man: install
	mkdir -p man
	kafka-connect --help-man > man/kafka-connect.1
	nroff -man man/kafka-connect.1
	@echo
	@echo -----------------------------------------
	@echo Man page generated at man/kafka-connect.1
	@echo -----------------------------------------

clean:
	$(RM) *.coverprofile
	$(RM) -r man

# In case you forget -s -v when using `glide get`.
clean-vendor:
	glide install --strip-vcs --strip-vendor

distclean: clean
	$(RM) -r dist/
	$(GO) clean -i $(packages)

.PHONY: build install test spec coverage browse-coverage
.PHONY: lint errcheck zen get-devtools
.PHONY: dist get-reltools man release
