package connect

import (
	"fmt"
	"log"
	"net/http"
)

// A Connector represents a Kafka Connect connector instance.
//
// See: http://docs.confluent.io/current/connect/userguide.html#connectors-tasks-and-workers
type Connector struct {
	Name   string          `json:"name"`
	Config ConnectorConfig `json:"config,omitempty"`
	Tasks  []TaskID        `json:"tasks,omitempty"`
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
type Task struct {
	ID     TaskID            `json:"id"`
	Config map[string]string `json:"config"`
}

// A TaskID has two components, a numerical ID and a connector name by which the
// ID is scoped.
type TaskID struct {
	ConnectorName string `json:"connector"`
	ID            int    `json:"task"`
}

// TODO: Probably need to URL-encode connector names

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

// ListConnectors retrieves a list of active connector names.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors
func (c *Client) ListConnectors() ([]string, *http.Response, error) {
	path := "connectors"
	var names []string
	response, err := c.get(path, &names)
	return names, response, err
}

// GetConnector retrieves information about a connector with the given name.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors-(string-name)
func (c *Client) GetConnector(name string) (*Connector, *http.Response, error) {
	path := "connectors/" + name
	connector := new(Connector)
	response, err := c.get(path, connector)
	return connector, response, err
}

// GetConnectorConfig retrieves configuration for a connector with the given
// name.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors-(string-name)-config
func (c *Client) GetConnectorConfig(name string) (ConnectorConfig, *http.Response, error) {
	path := fmt.Sprintf("connectors/%v/config", name)
	config := make(ConnectorConfig)
	response, err := c.get(path, &config)
	return config, response, err
}

// GetConnectorTasks retrieves a list of tasks currently running for a connector
// with the given name.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors-(string-name)-tasks
func (c *Client) GetConnectorTasks(name string) ([]Task, *http.Response, error) {
	path := fmt.Sprintf("connectors/%v/tasks", name)
	var tasks []Task
	response, err := c.get(path, &tasks)
	return tasks, response, err
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

func (c *Client) get(path string, v interface{}) (*http.Response, error) {
	request, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.Do(request, v)
	if err != nil {
		return response, err
	}

	return response, err
}
