package connect

import (
	"errors"
	"fmt"
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

// CreateConnector creates a new connector instance. If successful, conn is
// updated with the connector's state returned by the API, including Tasks.
//
// Passing an object that already contains Tasks produces an error.
//
// See: http://docs.confluent.io/current/connect/userguide.html#post--connectors
func (c *Client) CreateConnector(conn *Connector) (*http.Response, error) {
	if len(conn.Tasks) != 0 {
		return nil, errors.New("Cannot create Connector with existing Tasks")
	}
	path := "connectors"
	response, err := c.doRequest("POST", path, conn, conn)
	return response, err
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
// the given name, returning the new state of the Connector.
//
// If the connector does not exist, it will be created, and the returned HTTP
// response will indicate a 201 Created status.
//
// See: http://docs.confluent.io/current/connect/userguide.html#put--connectors-(string-name)-config
func (c *Client) UpdateConnectorConfig(name string, config ConnectorConfig) (*Connector, *http.Response, error) {
	path := fmt.Sprintf("connectors/%v/config", name)
	connector := new(Connector)
	response, err := c.doRequest("PUT", path, config, connector)
	return connector, response, err
}

// DeleteConnector deletes a connector with the given name, halting all tasks
// and deleting its configuration.
//
// See: http://docs.confluent.io/current/connect/userguide.html#delete--connectors-(string-name)-
func (c *Client) DeleteConnector(name string) (*http.Response, error) {
	return c.delete("connectors/" + name)
}

// PauseConnector pauses a connector and its tasks, which stops message
// processing until the connector is resumed. Tasks will transition to PAUSED
// state asynchronously.
//
// See: http://docs.confluent.io/current/connect/userguide.html#put--connectors-(string-name)-pause
func (c *Client) PauseConnector(name string) (*http.Response, error) {
	path := fmt.Sprintf("connectors/%v/pause", name)
	return c.doRequest("PUT", path, nil, nil)
}

// ResumeConnector resumes a paused connector. Tasks will transition to RUNNING
// state asynchronously.
//
// See: http://docs.confluent.io/current/connect/userguide.html#put--connectors-(string-name)-resume
func (c *Client) ResumeConnector(name string) (*http.Response, error) {
	path := fmt.Sprintf("connectors/%v/resume", name)
	return c.doRequest("PUT", path, nil, nil)
}

// RestartConnector restarts a connector and its tasks.
//
// See http://docs.confluent.io/current/connect/userguide.html#post--connectors-(string-name)-restart
func (c *Client) RestartConnector(name string) (*http.Response, error) {
	path := fmt.Sprintf("connectors/%v/restart", name)
	return c.doRequest("POST", path, nil, nil)
}
