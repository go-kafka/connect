# This makefile probably requires GNU make >= 3.81

GO ?= go
# OS X, use sha256sum or gsha256sum elsewhere
SHASUM := shasum -a 256
VERSION := $(shell git describe)

packages := ./...
fordist := find * -type d -exec

# Should use -mod=readonly in CI
build:
	$(GO) build $(packages)

install:
	$(GO) install $(packages)

# Use `test all` now for CI? https://research.swtch.com/vgo-cmd
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

check: lint errcheck

lint:
	golint --set_exit_status ./...

errcheck:
	errcheck --asserts --exclude=err-excludes.txt $(packages)

zen:
	ginkgo watch -notify $(packages)

get-devtools:
	@echo Getting golint...
	$(GO) install golang.org/x/lint/golint
	@echo Getting errcheck...
	$(GO) install github.com/kisielk/errcheck

get-reltools:
	@echo Getting gox...
	$(GO) install github.com/mitchellh/gox

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

distclean: clean
	$(RM) -r dist/
	$(GO) clean -i $(packages)

.PHONY: build install test spec coverage browse-coverage
.PHONY: check lint errcheck zen get-devtools
.PHONY: dist get-reltools man release
