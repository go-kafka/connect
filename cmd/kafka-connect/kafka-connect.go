package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/go-kafka/connect"
)

// A connectorAction is a function that performs an imperative action on a
// Connector by name.
type connectorAction func(name string) (*http.Response, error)

var (
	//-------------------------------------------------------------------------
	// CLI Commands, Args, Flags with Kingpin
	//-------------------------------------------------------------------------

	app = kingpin.New("kafka-connect", "Command line utility for managing Kafka Connect.").
		Version("kafka-connect CLI " + connect.VERSION).
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

	pauseCmd  = app.Command("pause", "Pause a connector and its tasks.")
	pauseName = pauseCmd.Arg("name", "Name of the connector to pause.").Required().String()

	resumeCmd  = app.Command("resume", "Resume a paused connector.")
	resumeName = resumeCmd.Arg("name", "Name of the connector to resume.").Required().String()

	restartCmd  = app.Command("restart", "Restart a connector and its tasks.")
	restartName = restartCmd.Arg("name", "Name of the connector to restart.").Required().String()

	// TODO: New stuff
	// plugin subcommand: list (default), validate
)

func main() {
	app.HelpFlag.Short('h')

	argv := os.Args[1:]
	subcommand, err := app.Parse(argv)

	if err != nil {
		context, _ := app.ParseContext(argv)
		app.FatalUsageContext(context, err.Error())
	}

	if !(*host).IsAbs() {
		app.Fatalf("host %v is not a valid absolute URL", *host)
	}

	// Localize use of os.Exit because it doesn't run deferreds
	app.FatalIfError(run(subcommand), "")
}

func run(subcommand string) error {
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
		// TODO: verify error output of 409 Conflict
		return affectConnector(*deleteName, client.DeleteConnector, "Deleted")

	case showCmd.FullCommand():
		apiResult, _, err = client.GetConnector(*showName)

	case configCmd.FullCommand():
		apiResult, _, err = client.GetConnectorConfig(*configName)

	case tasksCmd.FullCommand():
		apiResult, _, err = client.GetConnectorTasks(*tasksName)

	case statusCmd.FullCommand():
		apiResult, _, err = client.GetConnectorStatus(*statusName)

	case pauseCmd.FullCommand():
		return affectConnector(*pauseName, client.PauseConnector, "Paused")

	case resumeCmd.FullCommand():
		return affectConnector(*resumeName, client.ResumeConnector, "Resumed")

	case restartCmd.FullCommand():
		// TODO: verify error output of 409 Conflict
		return affectConnector(*restartName, client.RestartConnector, "Restarted")
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

func affectConnector(name string, action connectorAction, desc string) error {
	_, err := action(name)
	if err == nil {
		fmt.Printf("%v connector %v.\n", desc, name)
	}

	return err
}

// TODO: Some kind of formatter abstraction
func formatPrettyJSON(v interface{}) (string, error) {
	pretty, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(pretty), nil
}
