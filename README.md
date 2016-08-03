Kafka Connect CLI
=================

[![GoDoc][godoc-badge]][godoc]

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

    $ kafka-connect
    usage: kafka-connect [<flags>] <command> [<args> ...]

    Command line utility for managing Kafka Connect.

    Flags:
          --help     Show context-sensitive help (also try --help-long and --help-man).
          --version  Show application version.
      -H, --host="http://localhost:8083/"
                    Host address for the Kafka Connect REST API instance.

    Commands:
      help [<command>...]
        Show help.

      list
        Lists active connectors. Aliased as 'ls'.

      create <name>
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

The process exits with a zero status when operations are successful and
non-zero in the case of errors.

### Manual Page ###

If you'd like a `man` page, you can generate one and place it on your
`MANPATH`:

```sh
$ kafka-connect --help-man > /usr/local/share/man/man1/kafka-connect.1
```

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

Options
-------

Expanded details for select parameters, including supported environment
variables:

- `--host / -H`: API host address, default `http://localhost:8083/`. Supports
  environment variable `KAFKA_CONNECT_CLI_HOST`. Note that you can target any
  host in a Kafka Connect cluster.

Others? `CLASSPATH` that the shell scripts support for connector plugins.

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
this repository's `vendor/` are for the CLI tool (the `cmd` sub-package, not
installed unless you append `/...` to the `go get` command above).

See the API documentation linked above for examples.

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

TODO
----

### Features ###

- [ ] Dynamic shell completion of connector names.
- [ ] Other output/input support: ASCII tables, Java properties.
- [ ] TLS/SSL?
- [ ] Logging?

### Enhancements ###

- [ ] Do something useful in known error conditions, like 409 for restart
  during rebalance.
- [ ] Output compact JSON when stdout is not a TTY, with an option to force.
  Mimic jq's options.
- [ ] More efficient byte stream de/encoding than unmarshaling JSON and then
  marshaling again to print it?
- [ ] Consider testing the CLI with Gomega's gexec features.

### Meta ###

- [ ] Decide and document versioning scheme. Might be best to version CLI
  separately.
- [] Use gopkg.in for the library's sake?
- [ ] Publish checksums/sigs for releases, document `gpg --verify` steps.
- [ ] Write the package documentation in `version.go` or `doc.go`.
- [ ] Add some examples for library usage.
- [ ] Drop vendored dependency sources when Glide merges gps solver.
- [ ] Follow [dropping protobuf test dependency from
  Gomega](https://github.com/onsi/gomega/issues/123)


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
