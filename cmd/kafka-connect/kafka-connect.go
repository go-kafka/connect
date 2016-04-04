package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/go-kafka/connect"
)

// VERSION is a programmatically-available declaration of the kafka-connect
// CLI's version. This tracks the associated Kafka minor release version,
// initially.
const VERSION = "0.9-pre"

var (
	//-------------------------------------------------------------------------
	// CLI Commands, Args, Flags with Kingpin
	//-------------------------------------------------------------------------

	app = kingpin.New("kafka-connect", "Command line utility for managing Kafka Connect.").
		Version(VERSION).
		Author("Ches Martin")

	// debug = app.Flag("debug", "Enable debug mode.").Envar("KAFKA_CONNECT_CLI_DEBUG").Bool()

	host = app.Flag("host", "Host address for the Kafka Connect REST API instance.").
		Envar("KAFKA_CONNECT_CLI_HOST").Short('H').Default("http://localhost:8083").String()

	listCmd = app.Command("list", "Lists active connectors. Aliased as 'ls'.").Alias("ls")

	// TODO: Take a properties file or JSON for config?
	createCmd  = app.Command("create", "Creates a new connector instance.")
	createName = createCmd.Arg("name", "Name of the connector to create.").Required().String()

	// TODO: Take a properties file or JSON for config?
	updateCmd  = app.Command("update", "Updates a connector.")
	updateName = updateCmd.Arg("name", "Name of the connector to update.").Required().String()

	deleteCmd  = app.Command("delete", "Deletes a connector. Aliased as 'rm'.").Alias("rm")
	deleteName = deleteCmd.Arg("name", "Name of the connector to delete.").Required().String()

	showCmd  = app.Command("show", "Shows information about a connector and its tasks.")
	showName = showCmd.Arg("name", "Name of the connector to look up.").Required().String()

	configCmd  = app.Command("config", "Displays configuration of a connector.")
	configName = configCmd.Arg("name", "Name of the connector to look up.").Required().String()

	tasksCmd  = app.Command("tasks", "Displays tasks currently running for a connector.")
	tasksName = tasksCmd.Arg("name", "Name of the connector to look up.").Required().String()
)

func main() {
	subcommand := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Dispatch subcommands
	switch subcommand {
	case listCmd.FullCommand():
		connect.GetConnectors()

	case createCmd.FullCommand():
		connect.CreateConnector(connect.Connector{Name: *createName})

	case updateCmd.FullCommand():
		connect.UpdateConnectorConfig(*updateName, connect.ConnectorConfig{})

	case deleteCmd.FullCommand():
		connect.DeleteConnector(*deleteName)

	case showCmd.FullCommand():
		connect.GetConnector(*showName)

	case configCmd.FullCommand():
		connect.GetConnectorConfig(*configName)

	case tasksCmd.FullCommand():
		connect.GetConnectorTasks(*tasksName)
	}
}
