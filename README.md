Kafka Connect CLI
=================

[![Release][release-badge]][latest release]
[![Build Status][travis-badge]][build status]
[![Coverage Status][coverage-badge]][coverage status]
[![Go Report Card][go-report-badge]][go report card]
[![GoDoc][godoc-badge]][godoc]

A fast, portable, self-documenting CLI tool to inspect and manage [Kafka
Connect] connectors via the [REST API]. Because you don't want to be fumbling
through runbooks of `curl` commands when something's going wrong, or ever
really.

This project also contains a Go library for the Kafka Connect API usable by
other Go tools or applications. See [Using the Go
Library](#using-the-go-library) for details.

Usage
-----

The tool is self-documenting: run `kafka-connect help` or `kafka-connect help
<subcommand>` when you need a reference. A summary of functionality:

    $ kafka-connect
    usage: kafka-connect [<flags>] <command> [<args> ...]

    Command line utility for managing Kafka Connect.

    Flags:
      -h, --help     Show context-sensitive help (also try --help-long and --help-man).
          --version  Show application version.
      -H, --host=http://localhost:8083/
                    Host address for the Kafka Connect REST API instance.

    Commands:
      help [<command>...]
        Show help.

      list
        Lists active connectors. Aliased as 'ls'.

      create [<flags>] [<name>]
        Creates a new connector instance.

      update <name>
        Updates a connector.

      delete <name>
        Deletes a connector. Aliased as 'rm'.

      show <name>
        Shows information about a connector and its tasks.

      config <name>
        Displays configuration of a connector.

      tasks <name>
        Displays tasks currently running for a connector.

      status <name>
        Gets current status of a connector.

      pause <name>
        Pause a connector and its tasks.

      resume <name>
        Resume a paused connector.

      restart <name>
        Restart a connector and its tasks.

      version
        Shows kafka-connect version information.

For examples, see [the Godoc page for the command][cmd doc].

The process exits with a zero status when operations are successful and
non-zero in the case of errors.

[cmd doc]: https://godoc.org/github.com/go-kafka/connect/cmd/kafka-connect

### Manual Page ###

If you'd like a `man` page, you can generate one and place it on your
`MANPATH`:

```sh
$ kafka-connect --help-man > /usr/local/share/man/man1/kafka-connect.1
```

### Options ###

Expanded details for select parameters:

- `--host / -H`: API host address, default `http://localhost:8083/`. Can be set
  with environment variable `KAFKA_CONNECT_CLI_HOST`. Note that you can target
  any host in a Kafka Connect cluster.

Installation
------------

Binary releases [are available on GitHub][releases], signed and with checksums.

Fetch the appropriate version for your platform and place it somewhere on your
`PATH`. The YOLO way:

```sh
$ curl -L https://github.com/go-kafka/connect/releases/download/v0.9.0/kafka-connect-v0.9.0-linux-amd64.zip
$ unzip kafka-connect-v0.9.0-linux-amd64.zip
$ mv linux-amd64/kafka-connect /usr/local/bin/
```

The prudent way:

```sh
$ curl -L https://github.com/go-kafka/connect/releases/download/v0.9.0/kafka-connect-v0.9.0-linux-amd64.zip
$ curl -L https://github.com/go-kafka/connect/releases/download/v0.9.0/kafka-connect-v0.9.0-linux-amd64.zip.sha256sum
# Verify integrity of the archive file, on OS X try shasum --check
$ sha256sum --check kafka-connect-v0.9.0-linux-amd64.zip.sha256sum
$ unzip kafka-connect-v0.9.0-linux-amd64.zip
$ mv linux-amd64/kafka-connect /usr/local/bin/
```

Or best of all, the careful way:

```sh
$ curl -L https://github.com/go-kafka/connect/releases/download/v0.9.0/kafka-connect-v0.9.0-linux-amd64.zip
$ unzip kafka-connect-v0.9.0-linux-amd64.zip
# Verify signature of the binary:
$ gpg --verify linux-amd64/kafka-connect{.asc,}
$ mv linux-amd64/kafka-connect /usr/local/bin/
```

You can find my GPG key distributed on keyservers with ID `8638EE95`. The
fingerprint is:

    23D6 18B5 3AB8 209F F172  C070 6E5C D3ED 8638 EE95

For a more detailed primer on GPG signatures and key authenticity, check out
[the Apache Software Foundation's doc](http://www.apache.org/info/verification.html).

*Cross-compiled binaries are possibly untested—please report any issues. If you
would like a binary build for a platform that is not currently published, I'm
happy to make one available as long as Go can cross-compile it without
problem—please open an issue.*

To build your own version from source, see the below [Building and
Development](#building-and-development) section.

### Command Completion ###

Shell completion is built in for bash and zsh, just add the following to your
shell profile initialization (`~/.bash_profile` or the like):

```sh
which kafka-connect >/dev/null && eval "$(kafka-connect --completion-script-bash)"
```

Predictably, use `--completion-script-zsh` for zsh.

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
normally with `go test` or `make test`. If you wish to use additional features
of [the Ginkgo CLI tool][ginkgo cli] like `watch` mode or generating stub test
files, etc. you'll need to install it on your regular `GOPATH` using:

    $ go install github.com/onsi/ginkgo/ginkgo

This tool is not included in `vendor/`, in part because [Glide doesn't support
that yet][glide execs], and also because it's optional.

### Using the Go Library ###

[![GoDoc][godoc-badge]][godoc]

To use the Go library, simply use `go get` and import it in your code as usual:

```sh
$ go get -u github.com/go-kafka/connect
```

The library has no dependencies beyond the standard library. Dependencies in
this repository's `vendor/` are for the CLI tool (the `cmd` sub-package, not
installed unless you append `/...` to the `go get` command above).

See the API documentation linked above for examples.

Versions
--------

For information about versioning policy and compatibility status please see
[the release notes](HISTORY.md).

Alternatives
------------

<https://github.com/datamountaineer/kafka-connect-tools>

When I wanted a tool like this, I found this one. It's written in Scala—I <3
Scala, but JVM start-up time is sluggish for CLI tools, and it's much easier to
distribute self-contained native binaries to management hosts that don't
require a JVM installed.

Similar things can be said of Kafka's packaged management scripts, which are
less ergonomic. Hence, I wrote this Go variant.

Kudos to the kafka-connect-tools authors for inspiration.

Contributing
------------

Please see [the Contributing Guide](CONTRIBUTING.md)!

License
-------

The library and CLI tool are made available under the terms of the MIT license,
see the [LICENSE](LICENSE) file for full details.


[Kafka Connect]: http://docs.confluent.io/current/connect/intro.html
[REST API]: http://docs.confluent.io/current/connect/userguide.html#rest-interface
[releases]: https://github.com/go-kafka/connect/releases
[Glide]: https://glide.sh/
[write go]: https://golang.org/doc/install
[Ginkgo]: https://onsi.github.io/ginkgo/
[ginkgo cli]: https://onsi.github.io/ginkgo/#the-ginkgo-cli
[glide execs]: https://github.com/Masterminds/glide/pull/331

[release-badge]: https://img.shields.io/github/release/go-kafka/connect.svg?maxAge=2592000
[latest release]: https://github.com/go-kafka/connect/releases/latest
[travis-badge]:https://travis-ci.org/go-kafka/connect.svg?branch=master
[build status]: https://travis-ci.org/go-kafka/connect
[coverage-badge]: https://codecov.io/gh/go-kafka/connect/branch/master/graph/badge.svg
[coverage status]: https://codecov.io/gh/go-kafka/connect
[go-report-badge]: https://goreportcard.com/badge/github.com/go-kafka/connect
[go report card]: https://goreportcard.com/report/github.com/go-kafka/connect
[godoc-badge]: http://img.shields.io/badge/godoc-reference-blue.svg?style=flat
[godoc]: https://godoc.org/github.com/go-kafka/connect

<!-- vim:set expandtab shiftwidth=2 textwidth=79: -->
