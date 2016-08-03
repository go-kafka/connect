package connect

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// A Connector represents a Kafka Connect connector instance.
//
// See: http://docs.confluent.io/current/connect/userguide.html#connectors-tasks-and-workers
//
// TODO: maybe omit Tasks field here, see POST /connectors endpoint
type Connector struct {
	Name   string          `json:"name"`
	Config ConnectorConfig `json:"config"`
	Tasks  []Task          `json:"tasks"`
}

// ConnectorConfig is a key-value mapping of configuration for connectors, where
// keys are in the form of Java properties.
//
// See: http://docs.confluent.io/current/connect/userguide.html#configuring-connectors
type ConnectorConfig map[string]string

// A Task is a unit of work dispatched by a Connector to parallelize the work of
// a data copy job.
//
// See: http://docs.confluent.io/current/connect/userguide.html#connectors-tasks-and-workers
//
// TODO: see /connectors/<name>/tasks endpoint -- might want to just encode this
// as an anon []struct member of Connector (verify the documented output, the
// text is ambiguous)
type Task struct {
	ConnectorName string `json:"connector"`
	ID            int    `json:"task"`
}

// TODO: Should we always return (success, error) pairs in case of HTTP errors?

// CreateConnector creates a new connector instance. It returns an error if
// creation is unsuccessful.
//
// See: http://docs.confluent.io/current/connect/userguide.html#post--connectors
//
// TODO: return the full connector info so that tasks are listed?
func CreateConnector(conn Connector) error {
	log.Println("Called CreateConnector")
	return nil
}

// GetConnectors retrieves a list of active connector names.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors
func GetConnectors() []string {
	log.Println("Called GetConnectors")
	return make([]string, 0)
}

// GetConnector retrieves information about a connector with the given name.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors-(string-name)
func GetConnector(name string) (Connector, error) {
	log.Println("Called GetConnector")
	return Connector{}, nil
}

// GetConnectorConfig retrieves configuration for a connector with the given
// name.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors-(string-name)-config
func GetConnectorConfig(name string) (ConnectorConfig, error) {
	log.Println("Called GetConnectorConfig")
	return ConnectorConfig{}, nil
}

// GetConnectorTasks retrieves a list of tasks currently running for a connector
// with the given name.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors-(string-name)-tasks
//
// TODO: Erm, the example output for this endpoint is totally inconsistent with
// what the docs describe as the output format. Need to test. And it says
// *request* format, when there should not be a request body...
func GetConnectorTasks(name string) ([]Task, error) {
	log.Println("Called GetConnectorTasks")
	return make([]Task, 0), nil
}

// UpdateConnectorConfig updates configuration for an existing connector with
// the given name. If the connector does not exist, it will be created.
//
// See: http://docs.confluent.io/current/connect/userguide.html#put--connectors-(string-name)-config
func UpdateConnectorConfig(name string, config ConnectorConfig) (Connector, error) {
	log.Println("Called UpdateConnectorConfig")
	return Connector{}, nil
}

// DeleteConnector deletes a connector with the given name, halting all tasks
// and deleting its configuration. Returns whether deletion was successful or
// not.
//
// TODO: return error instead, forwarding HTTP error?
//
// See: http://docs.confluent.io/current/connect/userguide.html#delete--connectors-(string-name)-
func DeleteConnector(name string) bool {
	log.Println("Called DeleteConnector")
	return true
}

// TODO: either return decoded response entity, or accept one to mutate
func sendRequest(url string) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		// TODO: parse APIError, implement Error() and return using the message
		return errors.New("kafka connect: HTTP error " + response.Status + " on " + url)
	}

	decoder := json.NewDecoder(response.Body)
	// TODO: decode to appropriate entity for request
	connector := Connector{}
	err = decoder.Decode(&connector)
	if err != nil {
		return err
	}

	return nil
}
