Kafka Connect CLI
=================

A fast, portable, self-documenting CLI tool to inspect and manage [Kafka
Connect] connectors via the [REST API]. Because you don't want to be fumbling
through runbooks of `curl` commands when something's going wrong, or ever
really.

This project also contains a Go library for the Kafka Connect API usable by
other Go-based tools or applications. See below for details.

Usage
-----

The tool is self-documenting: run `kafka-connect help` or `kafka-connect help
<subcommand>` when you need a reference.

    $ kafka-connect list
    $ kafka-connect create
    $ kafka-connect update
    $ kafka-connect delete
    $ kafka-connect show
    $ kafka-connect config
    $ kafka-connect tasks

Installation
------------

Binary releases [are available on GitHub][releases] with checksums.

Fetch the appropriate version for your platform and place it somewhere on your
`PATH`:

```sh
$ curl -L -o kafka-connect https://github.com/go-kafka/connect/releases/download/0.9/kafka-connect-0.9-linux-amd64
$ sha256sum kafka-connect  # Verify checksum for your arch with releases page
$ mv kafka-connect /usr/local/bin/
```

TODO: explain versioning scheme; publish checksums

*Cross-compiled binaries are possibly untested—please report any issues. If you
would like a binary build for a platform that is not currently published, I'm
happy to make one available as long as Go can cross-compile it without
problem—please open an issue.*

To build your own version from source, see the below *Building and Development*
section.

### Command Completion ###

Shell completion is built in for bash and zsh, just add the following to your
shell profile initialization (`~/.bash_profile` or the like):

```sh
which kafka-connect >/dev/null && eval "$(kafka-connect --completion-script-bash)"
```

Predictably, use `--completion-script-zsh` for zsh.

### Manual Page ###

If you'd like a `man` page, you can generate one and place it on your
`MANPATH`:

```sh
$ kafka-connect --help-man > /usr/local/share/man/man1/kafka-connect.1
```

Configuration
-------------

- API host address: `--host / -H`, default `http://localhost:8083/`, supports
  environment variable `KAFKA_CONNECT_CLI_HOST`

Others? `CLASSPATH` that the shell scripts support for connector plugins.

Output format? JSON for scripting, tabular, Java properties?

Building and Development
------------------------

This project is implemented in Go and uses [Glide] to achieve reproducible
builds. You don't need Glide unless you want to make changes to dependencies,
though, since the dependency sources are checked into source control.

Once you [have a working Go toolchain][write go], it is simple to build like
any standard Go project:

```sh
$ go get -d github.com/go-kafka/connect/...
$ cd $GOPATH/github.com/go-kafka/connect
$ go build    # or
$ go install  # or
$ go test     # etc.
```

If you are building with Go 1.5, be sure to `export GO15VENDOREXPERIMENT=1` to
resolve the dependencies in `vendor/`. This is default in Go 1.6.

Cross-compiling is again standard Go procedure: set `GOOS` and `GOARCH`. For
example if you wanted to build a CLI tool binary for Linux on ARM:

```sh
$ env GOOS=linux GOARCH=arm go build ./cmd/...
$ file ./kafka-connect
kafka-connect: ELF 32-bit LSB executable, ARM, version 1 (SYSV), statically linked, not stripped
```

#### Testing with Ginkgo ####

This project uses the [Ginkgo] BDD testing library. You can run the tests
normally with `go test`, but Ginkgo also has its own CLI tool [with additional
options][ginkgo cli] for running tests, watching for file changes, generating
stub test files, etc. This tool is not included in `vendor/`, in part because
[Glide doesn't support that yet][glide execs], so if you'd like to use it
you'll need to `go install github.com/onsi/ginkgo/ginkgo` to install it on your
regular `GOPATH`. For convenience there is also a `make get-devtools` task.

### Using the Go Library ###

[![GoDoc][godoc-badge]][godoc]

To use the Go library, simply use `go get` and import it in your code as usual:

```sh
$ go get -u github.com/go-kafka/connect
```

The library has no dependencies beyond the standard library. Dependencies in
this repository's `vendor/` are for the CLI tool (the `cmd` sub-package).

TODO: use gopkg.in?

Alternatives
------------

- <https://github.com/datamountaineer/kafka-connect-tools>

  When I wanted a tool like this, I found this one already available. But it's
  written in Scala—I <3 Scala, but JVM start-up time is sluggish for CLI tools,
  and it's much easier to distribute self-contained native binaries to
  management hosts that don't require a JVM installed. Hence, I wrote this Go
  variant. Kudos to the kafka-connect-tools authors for inspiration.

Contributing
------------

Please see [the Contributing Guide](CONTRIBUTING.md).


[Kafka Connect]: http://docs.confluent.io/current/connect/intro.html
[REST API]: http://docs.confluent.io/current/connect/userguide.html#rest-interface
[releases]: https://github.com/go-kafka/connect/releases
[Glide]: https://glide.sh/
[write go]: https://golang.org/doc/install
[Ginkgo]: https://onsi.github.io/ginkgo/
[ginkgo cli]: https://onsi.github.io/ginkgo/#the-ginkgo-cli
[glide execs]: https://github.com/Masterminds/glide/pull/331

[godoc-badge]: http://img.shields.io/badge/godoc-reference-blue.svg?style=flat
[godoc]: https://godoc.org/github.com/go-kafka/connect

<!-- vim:set expandtab shiftwidth=2 textwidth=79: -->
