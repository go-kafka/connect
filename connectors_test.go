package connect_test

import (
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/go-kafka/connect"
)

var _ = Describe("Connectors", func() {
	var client *Client
	var server *ghttp.Server

	jsonAcceptHeader := http.Header{"Accept": []string{"application/json"}}

	fileSourceConfig := ConnectorConfig{
		"connector.class": "FileStreamSource",
		"file":            "/tmp/test.txt",
		"name":            "local-file-source",
		"tasks.max":       "1",
		"topic":           "go-kafka-connect-test",
	}

	BeforeEach(func() {
		server = ghttp.NewServer()
		url, _ := url.Parse(server.URL())
		client = NewClient(nil)
		client.BaseURL = url
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("CreateConnector", func() {
		It("creates a new instance given a valid Connector", func() {
			connector := Connector{Name: "test"}
			err := CreateConnector(connector)
			Expect(err).NotTo(HaveOccurred())
		})

		PIt("returns an error given an invalid Connector", func() {
		})
	})

	Describe("ListConnectors", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/connectors"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWith(http.StatusOK, `["test", "phony-hdfs-sink"]`),
				),
			)
		})

		It("returns list of connector names", func() {
			names, _, err := client.ListConnectors()
			Expect(err).NotTo(HaveOccurred())
			Expect(names).To(Equal([]string{"test", "phony-hdfs-sink"}))
		})
	})

	Describe("GetConnector", func() {
		var resultConnector *Connector
		var statusCode int

		BeforeEach(func() {
			resultConnector = &Connector{
				Name:   "local-file-source",
				Config: fileSourceConfig,
				Tasks:  []TaskID{TaskID{"local-file-source", 0}},
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/connectors/local-file-source"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &resultConnector),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
			})

			It("returns a connector", func() {
				connector, _, err := client.GetConnector("local-file-source")
				Expect(err).NotTo(HaveOccurred())
				Expect(connector).To(Equal(resultConnector))
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns an error response", func() {
				connector, resp, err := client.GetConnector("local-file-source")
				Expect(err).To(HaveOccurred())
				Expect(*connector).To(BeZero())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("GetConnectorConfig", func() {
		var statusCode int

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/connectors/local-file-source/config"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &fileSourceConfig),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
			})

			It("returns configuration for connector", func() {
				config, _, err := client.GetConnectorConfig("local-file-source")
				Expect(err).NotTo(HaveOccurred())
				Expect(config).To(Equal(fileSourceConfig))
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns an error response", func() {
				config, resp, err := client.GetConnectorConfig("local-file-source")
				Expect(err).To(HaveOccurred())
				Expect(config).To(BeEmpty())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("GetConnectorTasks", func() {
		var resultTasks []Task
		var statusCode int

		BeforeEach(func() {
			resultTasks = []Task{
				Task{
					ID: TaskID{"local-file-source", 0},
					Config: map[string]string{
						"file":       "/tmp/test.txt",
						"task.class": "org.apache.kafka.connect.file.FileStreamSourceTask",
						"topic":      "go-kafka-connect-test",
					},
				},
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/connectors/local-file-source/tasks"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &resultTasks),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
			})

			It("returns tasks for an extant connector", func() {
				tasks, _, err := client.GetConnectorTasks("local-file-source")
				Expect(err).NotTo(HaveOccurred())
				Expect(tasks).To(Equal(resultTasks))
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns an error response", func() {
				tasks, resp, err := client.GetConnectorTasks("local-file-source")
				Expect(err).To(HaveOccurred())
				Expect(tasks).To(BeEmpty())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("GetConnectorStatus", func() {
		var resultStatus *ConnectorStatus
		var statusCode int

		BeforeEach(func() {
			resultStatus = &ConnectorStatus{
				Name: "local-file-source",
				Connector: ConnectorState{
					State:    "RUNNING",
					WorkerID: "127.0.0.1:8083",
				},
				Tasks: []TaskState{
					TaskState{
						ID:       0,
						State:    "RUNNING",
						WorkerID: "127.0.0.1:8083",
					},
				},
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/connectors/local-file-source/status"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &resultStatus),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
			})

			It("returns connector status", func() {
				status, _, err := client.GetConnectorStatus("local-file-source")
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(resultStatus))
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns an error response", func() {
				status, resp, err := client.GetConnectorStatus("local-file-source")
				Expect(err).To(HaveOccurred())
				Expect(*status).To(BeZero())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("UpdateConnectorConfig", func() {
		It("returns updated connector when successful", func() {
			connector, err := UpdateConnectorConfig("test", ConnectorConfig{})
			Expect(err).NotTo(HaveOccurred())
			Expect(connector).NotTo(BeNil()) // TODO: assert properties of the connector
		})
	})

	Describe("DeleteConnector", func() {
		It("deletes a connector instance by name", func() {
			deleted := DeleteConnector("test")
			Expect(deleted).To(BeTrue())
		})

		PIt("returns false given an nonexistent connector name", func() {
			deleted := DeleteConnector("invalid")
			Expect(deleted).To(BeFalse())
		})
	})
})
