# This makefile probably requires GNU make >= 3.81

build:
	go build

install: build
	go install $$(glide novendor)

test:
	go test $$(glide novendor)

# In case you forget -s -v when using `glide get`.
clean-vendor:
	glide install --strip-vcs --strip-vendor

.PHONY: build install test clean-vendor
