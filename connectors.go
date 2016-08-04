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

// ConnectorStatus reflects the status of a Connector and state of its Tasks.
//
// Having connector name and a "connector" object at top level is a little
// awkward and produces stuttering, but it's their design, not ours.
type ConnectorStatus struct {
	Name      string         `json:"name"`
	Connector ConnectorState `json:"connector"`
	Tasks     []TaskState    `json:"tasks"`
}

// ConnectorState reflects the running state of a Connector and the worker where
// it is running.
type ConnectorState struct {
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
}

// TaskState reflects the running state of a Task and the worker where it is
// running.
type TaskState struct {
	ID       int    `json:"id"`
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
	Trace    string `json:"trace,omitempty"`
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

// GetConnectorStatus gets current status of the connector, including whether it
// is running, failed or paused, which worker it is assigned to, error
// information if it has failed, and the state of all its tasks.
//
// See: http://docs.confluent.io/current/connect/userguide.html#get--connectors-(string-name)-status
func (c *Client) GetConnectorStatus(name string) (*ConnectorStatus, *http.Response, error) {
	path := fmt.Sprintf("connectors/%v/status", name)
	status := new(ConnectorStatus)
	response, err := c.get(path, status)
	return status, response, err
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
// See: http://docs.confluent.io/current/connect/userguide.html#delete--connectors-(string-name)-
func (c *Client) DeleteConnector(name string) (*http.Response, error) {
	return c.delete(fmt.Sprintf("connectors/%v", name))
}

func (c *Client) get(path string, v interface{}) (*http.Response, error) {
	return c.doRequest("GET", path, nil, v)
}

func (c *Client) delete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil, nil)
}

func (c *Client) doRequest(method, path string, body, v interface{}) (*http.Response, error) {
	request, err := c.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	return c.Do(request, v)
}
