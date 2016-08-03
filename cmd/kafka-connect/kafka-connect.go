package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/go-kafka/connect"
)

var (
	//-------------------------------------------------------------------------
	// CLI Commands, Args, Flags with Kingpin
	//-------------------------------------------------------------------------

	app = kingpin.New("kafka-connect", "Command line utility for managing Kafka Connect.").
		Version(connect.VERSION).
		Author("Ches Martin")

	// debug = app.Flag("debug", "Enable debug mode.").Envar("KAFKA_CONNECT_CLI_DEBUG").Bool()

	host = app.Flag("host", "Host address for the Kafka Connect REST API instance.").
		Envar("KAFKA_CONNECT_CLI_HOST").Short('H').Default(connect.DefaultHostURL).String()

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

	statusCmd  = app.Command("status", "Gets current status of a connector.")
	statusName = statusCmd.Arg("name", "Name of the connector to look up.").Required().String()

	// TODO: New stuff
	// pause
	// resume
	// restart
	// plugins
)

func main() {
	// Localize use of os.Exit because it doesn't run deferreds
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	subcommand := kingpin.MustParse(app.Parse(os.Args[1:]))

	client := connect.NewClient(nil)
	var apiResult interface{}
	var err error
	var output string

	// Dispatch subcommands
	switch subcommand {
	case listCmd.FullCommand():
		apiResult, _, err = client.ListConnectors()

	case createCmd.FullCommand():
		connect.CreateConnector(connect.Connector{Name: *createName})

	case updateCmd.FullCommand():
		connect.UpdateConnectorConfig(*updateName, connect.ConnectorConfig{})

	case deleteCmd.FullCommand():
		connect.DeleteConnector(*deleteName)

	case showCmd.FullCommand():
		apiResult, _, err = client.GetConnector(*showName)

	case configCmd.FullCommand():
		apiResult, _, err = client.GetConnectorConfig(*configName)

	case tasksCmd.FullCommand():
		apiResult, _, err = client.GetConnectorTasks(*tasksName)

	case statusCmd.FullCommand():
		apiResult, _, err = client.GetConnectorStatus(*statusName)
	}

	if err != nil {
		return err
	}

	if output, err = formatPrettyJSON(apiResult); err != nil {
		return err
	}

	fmt.Println(output)
	return nil
}

// TODO: Some kind of formatter abstraction
func formatPrettyJSON(v interface{}) (string, error) {
	pretty, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(pretty), nil
}
