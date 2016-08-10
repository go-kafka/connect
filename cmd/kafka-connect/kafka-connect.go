package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/go-kafka/connect"
)

const hostenv = "KAFKA_CONNECT_CLI_HOST"

// ValidationError indicates that command arguments break an expected invariant.
type ValidationError struct {
	Message string

	// Flag suggesting that a handler should display contextual usage help in
	// addition to the short error message (e.g. kingpin.FatalUsageContext).
	SuggestUsage bool
}

func (e ValidationError) Error() string {
	return e.Message
}

// A connectorAction is a function that performs an imperative action on a
// Connector by name.
type connectorAction func(name string) (*http.Response, error)

var (
	pipedinput bool

	host     *url.URL
	connName string

	// For matching which execution we dispatch without proliferating strings
	listCmd, createCmd, updateCmd, deleteCmd *kingpin.CmdClause
	showCmd, configCmd, tasksCmd, statusCmd  *kingpin.CmdClause
	pauseCmd, resumeCmd, restartCmd          *kingpin.CmdClause

	newConnectorFilePath, connectorConfigPath string
)

func init() {
	pipedinput = !isatty(os.Stdin)
}

// BuildApp constructs the kafka-connect command line interface.
func BuildApp() *kingpin.Application {
	app := kingpin.New("kafka-connect", "Command line utility for managing Kafka Connect.").
		Version("kafka-connect CLI " + connect.VERSION).
		Author("Ches Martin").
		UsageWriter(os.Stdout)

	app.HelpFlag.Short('h')

	app.Flag("host", "Host address for the Kafka Connect REST API instance.").
		Short('H').
		Default(connect.DefaultHostURL).
		Envar(hostenv).
		URLVar(&host)

	// The modular style of Kingpin setup might cut down on the non-local vars,
	// but it feels pretty heavy and less declarative, so I'm undecided...
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
	addConnectorNameArg := func(cmdName, hint string, required bool) {
		command := app.GetCommand(cmdName)
		desc := fmt.Sprintf("Name of the connector to %v.", hint)
		if required {
			command.Arg("name", desc).Required().StringVar(&connName)
		} else {
			command.Arg("name", desc).StringVar(&connName)
		}
	}

	addConnectorNameArg("create", "create", false)
	hintedByName := []string{"update", "delete", "show", "pause", "resume", "restart"}
	for _, name := range hintedByName {
		addConnectorNameArg(name, name, true)
	}
	for _, name := range []string{"config", "tasks", "status"} {
		addConnectorNameArg(name, "look up", true)
	}

	createCmd.Flag("from-file", "A JSON file matching API request format, including connector name.").
		Short('f').
		PlaceHolder("FILE").
		ExistingFileVar(&newConnectorFilePath)
	createCmd.Flag("config", "A JSON file containing connector config.").
		Short('c').
		PlaceHolder("FILE").
		ExistingFileVar(&connectorConfigPath)

	updateCmd.Flag("config", "A JSON file containing connector config.").
		Short('c').
		PlaceHolder("FILE").
		ExistingFileVar(&connectorConfigPath)

	// Re-initialize global state for in-process tests, yeah kinda gross
	connName, newConnectorFilePath, connectorConfigPath = "", "", ""
	host = nil

	return app
}

// ValidateArgs parses and validates CLI arguments, isolated from execution
// logic. Returns user's parsed subcommand if invocation is valid, else error.
func ValidateArgs(app *kingpin.Application, argv []string) (subcommand string, err error) {
	// TODO: Employ kingpin Validate, but it's currently difficult to use:
	// https://github.com/alecthomas/kingpin/issues/125
	if subcommand, err = app.Parse(argv); err != nil {
		return
	}

	if !(*host).IsAbs() {
		msg := fmt.Sprintf("host %v is not a valid absolute URL", (*host).String())
		if os.Getenv(hostenv) != "" {
			msg += fmt.Sprintf(" (set by %v)", hostenv)
		}
		err = ValidationError{msg, false}
		return
	}

	switch subcommand {
	case createCmd.FullCommand():
		if pipedinput && (newConnectorFilePath != "" || connectorConfigPath != "") {
			err = ValidationError{"--from-file and --config cannot be used with input from stdin", false}
			return
		}

		if connName == "" {
			if connectorConfigPath != "" {
				err = ValidationError{"--config requires a connector name", true}
				return
			}
			if newConnectorFilePath == "" && !pipedinput {
				err = ValidationError{"either a connector name or --from-file is required", true}
				return
			}
		} else {
			if connectorConfigPath == "" && !pipedinput {
				err = ValidationError{"--config is required with a connector name", true}
				return
			}
			// Kingpin v3 might give us first-class mutual exclusivity support with
			// nice usage output: https://github.com/alecthomas/kingpin/issues/103
			if newConnectorFilePath != "" {
				err = ValidationError{"--from-file and --config are mutually exclusive", true}
				return
			}
		}
	case updateCmd.FullCommand():
		if pipedinput && connectorConfigPath != "" {
			err = ValidationError{"--config cannot be used with input from stdin", false}
			return
		}
		if connectorConfigPath == "" && !pipedinput {
			err = ValidationError{"configuration input is required, try --config or pipe to stdin", true}
			return
		}
	}

	return
}

func main() {
	app := BuildApp()
	argv := os.Args[1:]
	subcommand, err := ValidateArgs(app, argv)

	if err != nil {
		if verr, ok := err.(ValidationError); ok && !verr.SuggestUsage {
			app.Fatalf(verr.Error())
		}
		context, _ := app.ParseContext(argv)
		app.FatalUsageContext(context, err.Error())
	}

	// Localize use of os.Exit because it doesn't run deferreds
	app.FatalIfError(run(subcommand), "")
}

func run(subcommand string) error {
	client := connect.NewClient(nil)
	if host != nil {
		client.BaseURL = host
	}

	// Dispatch subcommands
	switch subcommand {
	case listCmd.FullCommand():
		return maybePrintAPIResult(client.ListConnectors())

	case createCmd.FullCommand():
		// TODO: verify/improve error output of 409 Conflict
		return createConnector(connName, client)

	case updateCmd.FullCommand():
		config, err := decodeConnectorConfig()
		if err != nil {
			return err
		}
		return maybePrintAPIResult(client.UpdateConnectorConfig(connName, config))

	case deleteCmd.FullCommand():
		// TODO: verify error output of 409 Conflict
		return affectConnector(connName, client.DeleteConnector, "Deleted")

	case showCmd.FullCommand():
		return maybePrintAPIResult(client.GetConnector(connName))

	case configCmd.FullCommand():
		return maybePrintAPIResult(client.GetConnectorConfig(connName))

	case tasksCmd.FullCommand():
		return maybePrintAPIResult(client.GetConnectorTasks(connName))

	case statusCmd.FullCommand():
		return maybePrintAPIResult(client.GetConnectorStatus(connName))

	case pauseCmd.FullCommand():
		return affectConnector(connName, client.PauseConnector, "Paused")

	case resumeCmd.FullCommand():
		return affectConnector(connName, client.ResumeConnector, "Resumed")

	case restartCmd.FullCommand():
		// TODO: verify error output of 409 Conflict
		return affectConnector(connName, client.RestartConnector, "Restarted")

	default: // won't reach here, arg parsing handles unknown commands
		return fmt.Errorf("Command `%v` is missing implementation!", subcommand)
	}
}

func maybePrintAPIResult(data interface{}, resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if output, err := formatPrettyJSON(data); err == nil {
		fmt.Println(output)
	}

	return err
}

func affectConnector(name string, action connectorAction, desc string) error {
	_, err := action(name)
	if err == nil {
		fmt.Printf("%v connector %v.\n", desc, name)
	}

	return err
}

func createConnector(name string, client *connect.Client) (err error) {
	var connector connect.Connector

	// Are we creating from --from-file, --config, or stdin?
	getDataSource := func() string {
		if pipedinput {
			return os.Stdin.Name()
		}
		if newConnectorFilePath != "" {
			return newConnectorFilePath
		}
		return connectorConfigPath
	}

	// TODO: more error message context
	// TODO: should really buffer stdin just in case...
	unmarshalInto := func(target interface{}, filepath string) error {
		if contents, err := ioutil.ReadFile(filepath); err == nil {
			return json.Unmarshal(contents, target)
		}
		return err
	}

	readFrom := getDataSource()
	if readFrom == newConnectorFilePath || name == "" {
		// kafka-connect create --from-file
		// cat connector.json | kafka-connect create
		err = unmarshalInto(&connector, readFrom)
	} else {
		// kafka-connect create my-conn-name --config
		// cat config.json | kafka-connect create my-conn-name
		var config connect.ConnectorConfig
		err = unmarshalInto(&config, readFrom)
		connector = connect.Connector{Name: name, Config: config}
	}

	if err != nil {
		return
	}

	// The API dubiously allows creating connectors with blank names... That's
	// probably a mistake, let's try to avoid it. It also sometimes returns name
	// as an attribute of config, so this might be present in roundtrip
	// scripting.
	// TODO: could apply this workaround in the library, but I'd rather get
	// the behavior acknowledged as a bug and not do that.
	if connector.Name == "" && connector.Config["name"] != "" {
		connector.Name = connector.Config["name"]
	}

	if _, err = client.CreateConnector(&connector); err == nil {
		if output, err := formatPrettyJSON(connector); err == nil {
			fmt.Println(output)
		}
	}

	return
}

func decodeConnectorConfig() (config connect.ConnectorConfig, err error) {
	var source string

	if pipedinput {
		source = os.Stdin.Name()
	} else {
		source = connectorConfigPath
	}

	contents, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &config)
	return config, err
}

// TODO: Some kind of formatter abstraction
func formatPrettyJSON(v interface{}) (string, error) {
	pretty, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(pretty), nil
}

// TODO: This probably doesn't work on Windows.
// https://github.com/mattn/go-isatty
func isatty(file *os.File) bool {
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
