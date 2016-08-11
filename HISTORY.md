Release Notes
=============

The CLI tool and library are versioned independently and both adhere to
[Semantic Versioning] policy. Compatibility with Kafka Connect API versions
will be noted in the release history below.

Status
------

The library is tentatively at 1.0 status pending a trial release, the API is
expected to be stable without breaking changes for the Kafka 0.10.0.x series.
The CLI's output, exit codes, etc. are provisional and subject to change until
a 1.0 release—please consult the change log before installing new versions if
you rely on these for scripting.

kafka-connect CLI
-----------------

### v0.9.0 - 11 August, 2016 ###

**Library version: v0.9.0**

Initial release. Covers nearly all API functionality (minus connector plugins),
but output may not yet be stable.

go-kafka/connect Library
------------------------

### v0.9.0 - 11 August, 2016 ###

Initial release. Supports the Kafka Connect REST API as of Kafka v0.10.0.0.


[Semantic Versioning]: http://semver.org/
<!-- vim:set tw=79: -->
