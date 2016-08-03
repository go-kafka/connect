# This makefile probably requires GNU make >= 3.81

GO := go

build:
	$(GO) build

install: build
	$(GO) install $$(glide novendor)

test: build
	ginkgo -v $$(glide novendor)

zen:
	ginkgo watch -notify $$(glide novendor)

get-devtools:
	@echo Getting the Ginkgo test runner...
	$(GO) get -v github.com/onsi/ginkgo/ginkgo

clean:
	$(RM) *.coverprofile
	$(RM) -r man

# In case you forget -s -v when using `glide get`.
clean-vendor:
	glide install --strip-vcs --strip-vendor

distclean: clean
	$(GO) clean -i github.com/go-kafka/connect...

.PHONY: build install test get-devtools clean clean-vendor distclean
