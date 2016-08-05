package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/go-kafka/connect"
)

// A connectorAction is a function that performs an imperative action on a
// Connector by name.
type connectorAction func(name string) (*http.Response, error)

var (
	host     *url.URL
	connName string

	// For matching which execution we dispatch without proliferating strings
	listCmd, createCmd, updateCmd, deleteCmd *kingpin.CmdClause
	showCmd, configCmd, tasksCmd, statusCmd  *kingpin.CmdClause
	pauseCmd, resumeCmd, restartCmd          *kingpin.CmdClause
)

// The modular style of Kingpin setup might cut down on the non-local vars, but
// it feels pretty heavy and less declarative, so I'm undecided...
func buildApp() *kingpin.Application {
	app := kingpin.New("kafka-connect", "Command line utility for managing Kafka Connect.").
		Version("kafka-connect CLI " + connect.VERSION).
		Author("Ches Martin").
		UsageWriter(os.Stdout)

	app.HelpFlag.Short('h')

	app.Flag("host", "Host address for the Kafka Connect REST API instance.").
		Short('H').
		Default(connect.DefaultHostURL).
		Envar("KAFKA_CONNECT_CLI_HOST").
		URLVar(&host)

	listCmd = app.Command("list", "Lists active connectors. Aliased as 'ls'.").Alias("ls")
	createCmd = app.Command("create", "Creates a new connector instance.")
	updateCmd = app.Command("update", "Updates a connector.")
	deleteCmd = app.Command("delete", "Deletes a connector. Aliased as 'rm'.").Alias("rm")
	showCmd = app.Command("show", "Shows information about a connector and its tasks.")
	configCmd = app.Command("config", "Displays configuration of a connector.")
	tasksCmd = app.Command("tasks", "Displays tasks currently running for a connector.")
	statusCmd = app.Command("status", "Gets current status of a connector.")
	pauseCmd = app.Command("pause", "Pause a connector and its tasks.")
	resumeCmd = app.Command("resume", "Resume a paused connector.")
	restartCmd = app.Command("restart", "Restart a connector and its tasks.")

	// TODO: New stuff
	// plugin subcommand: list (default), validate

	// Most commands need a connector name, reduce the boilerplate.
	hintedByName := []string{"create", "update", "delete", "show", "pause", "resume", "restart"}
	for _, name := range hintedByName {
		addConnectorNameArg(app, name, name)
	}
	for _, name := range []string{"config", "tasks", "status"} {
		addConnectorNameArg(app, name, "look up")
	}

	return app
}

func addConnectorNameArg(app *kingpin.Application, cmdName, hint string) {
	command := app.GetCommand(cmdName)
	desc := fmt.Sprintf("Name of the connector to %v.", hint)
	command.Arg("name", desc).Required().StringVar(&connName)
}

func main() {
	app := buildApp()
	argv := os.Args[1:]
	subcommand, err := app.Parse(argv)

	if err != nil {
		context, _ := app.ParseContext(argv)
		app.FatalUsageContext(context, err.Error())
	}

	// TODO: use kingpin Validate, but it's currently difficult to use:
	// https://github.com/alecthomas/kingpin/issues/125
	if !(*host).IsAbs() {
		app.Fatalf("host %v is not a valid absolute URL", (*host).String())
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
		connect.CreateConnector(connect.Connector{Name: connName})

	case updateCmd.FullCommand():
		connect.UpdateConnectorConfig(connName, connect.ConnectorConfig{})

	case deleteCmd.FullCommand():
		// TODO: verify error output of 409 Conflict
		return affectConnector(connName, client.DeleteConnector, "Deleted")

	case showCmd.FullCommand():
		apiResult, _, err = client.GetConnector(connName)

	case configCmd.FullCommand():
		apiResult, _, err = client.GetConnectorConfig(connName)

	case tasksCmd.FullCommand():
		apiResult, _, err = client.GetConnectorTasks(connName)

	case statusCmd.FullCommand():
		apiResult, _, err = client.GetConnectorStatus(connName)

	case pauseCmd.FullCommand():
		return affectConnector(connName, client.PauseConnector, "Paused")

	case resumeCmd.FullCommand():
		return affectConnector(connName, client.ResumeConnector, "Resumed")

	case restartCmd.FullCommand():
		// TODO: verify error output of 409 Conflict
		return affectConnector(connName, client.RestartConnector, "Restarted")
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
