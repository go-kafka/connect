/*
kafka-connect is a command line utility for managing Kafka Connect.

With it, you can inspect the status of connector instances running in a Kafka
cluster, start new connectors or update the configuration of existing ones, or
invoke lifecycle operations like pausing or restarting connectors.

For usage information, please see:

	kafka-connect help
	kafka-connect help <subcommand>

The tool is mostly self-documenting in this manner, but note that you can also
pass data for creating and updating connectors via standard input. These two
invocations are equivalent, but one may be more convenient for a particular
scripting task:

	kafka-connect create --from-file connector.json
	cat connector.json | kafka-connect create

Updating works similarly:

	kafka-connect update connector-name --config config.json
	cat config.json | kafka-connect update connector-name

In these examples, connector.json represents a JSON structure accepted by the
Connect REST API for creating connectors, and config.json is only the config
object of such a structure. The latter can be obtained for an existing connector
with:

	kafka-connect config connector-name

And the former is the output of the show command minus active tasks—using the jq
tool:

	kafka-connect show connector-name | jq 'del(.tasks)'

If you have configurations, you can also create new connector instances by
specifying names for them on the command line:

	kafka-connect create new-connector --config config.json
	cat config.json | kafka-connect create new-connector

API Host

By default kafka-connect will attempt to make requests to a Kafka Connect API
instance running on localhost and the default API port of 8083. This can be
changed by giving a full URL with the --host (or -H) flag, or with the
environment variable KAFKA_CONNECT_CLI_HOST.

Putting it all together, you might migrate connectors from one cluster to
another:

	kafka-connect show connector-name | jq 'del(.tasks)' | \
		kafka-connect -H http://newcluster:8083 create
	kafka-connect delete connector-name

For complete details of the data structures, see the REST API documentation:
http://docs.confluent.io/latest/connect/userguide.html#connect-userguide-rest.
*/
package main
