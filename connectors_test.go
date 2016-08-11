package connect_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/go-kafka/connect"
)

var (
	client *Client
	server *ghttp.Server

	jsonAcceptHeader  = http.Header{"Accept": []string{"application/json"}}
	jsonContentHeader = http.Header{"Content-Type": []string{"application/json"}}
)

var _ = Describe("Connector CRUD", func() {
	BeforeEach(func() {
		server = ghttp.NewServer()
		client = NewClient(server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	fileSourceConfig := ConnectorConfig{
		"connector.class": "FileStreamSource",
		"file":            "/tmp/test.txt",
		"tasks.max":       "1",
		"topic":           "go-kafka-connect-test",
	}

	Describe("CreateConnector", func() {
		var connector, resultConnector Connector
		var statusCode int

		BeforeEach(func() {
			connector = Connector{
				Name:   "local-file-source",
				Config: fileSourceConfig,
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/connectors"),
					ghttp.VerifyHeader(jsonContentHeader),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.VerifyJSONRepresenting(connector),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &resultConnector),
				),
			)
		})

		Context("when a valid Connector is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusCreated
				resultConnector = connector
				resultConnector.Tasks = []TaskID{{"local-file-source", 0}}
				Expect(connector).NotTo(Equal(resultConnector))
			})

			It("updates reference with state of newly-created instance", func() {
				_, err := client.CreateConnector(&connector)
				Expect(err).NotTo(HaveOccurred())
				Expect(connector).To(Equal(resultConnector))
			})
		})

		// The API ought to return a 422 but it currently returns 500 instead
		// (and the response is text/html despite Accept).
		// TODO: report this upstream as a bug
		Context("when an invalid Connector is given", func() {
			var origConnector Connector

			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
				origConnector = connector
			})

			// TODO: if 422 is returned in the future (see above), assert on
			// APIError value.
			It("returns an error", func() {
				resp, err := client.CreateConnector(&connector)
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})

			It("does not mutate Connector reference", func() {
				_, err := client.CreateConnector(&connector)
				Expect(err).To(HaveOccurred())
				Expect(connector).To(Equal(origConnector))
			})
		})

		Context("when a Connector with extant Tasks is given", func() {
			BeforeEach(func() {
				connector.Tasks = []TaskID{{"local-file-source", 0}}
			})

			It("returns an error", func() {
				_, err := client.CreateConnector(&connector)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Cannot create Connector with existing Tasks"))
			})
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
		var resultConnector interface{}
		var statusCode int

		BeforeEach(func() {
			resultConnector = Connector{
				Name:   "local-file-source",
				Config: fileSourceConfig,
				Tasks:  []TaskID{{"local-file-source", 0}},
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
				Expect(*connector).To(Equal(resultConnector.(Connector)))
			})
		})

		Context("when nonexisting connector name is given", func() {
			apiError := APIError{
				Code:    404,
				Message: "Connector local-file-source not found",
			}

			BeforeEach(func() {
				statusCode = http.StatusNotFound
				resultConnector = apiError
			})

			It("returns an error response", func() {
				connector, resp, err := client.GetConnector("local-file-source")
				Expect(err).To(MatchError(err.(APIError)))
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
				{
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

			It("returns tasks", func() {
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
					{
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
		var statusCode int

		resultConnector := Connector{
			Name:   "local-file-source",
			Config: fileSourceConfig,
			Tasks:  []TaskID{{"local-file-source", 0}},
		}

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/connectors/local-file-source/config"),
					ghttp.VerifyHeader(jsonContentHeader),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.VerifyJSONRepresenting(fileSourceConfig),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &resultConnector),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
			})

			It("returns updated connector", func() {
				connector, resp, err := client.UpdateConnectorConfig("local-file-source", fileSourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(connector.Config["file"]).To(Equal("/tmp/test.txt"))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusCreated
			})

			It("returns newly created connector with a Created response", func() {
				connector, resp, err := client.UpdateConnectorConfig("local-file-source", fileSourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(*connector).To(Equal(resultConnector))
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			})
		})
	})

	Describe("DeleteConnector", func() {
		var statusCode int

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/connectors/local-file-source"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWithPtr(&statusCode, nil),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusNoContent
			})

			It("deletes connector", func() {
				resp, err := client.DeleteConnector("local-file-source")
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			Context("when rebalance is in process", func() {
				BeforeEach(func() {
					statusCode = http.StatusConflict
				})

				It("returns error with a conflict response", func() {
					resp, err := client.DeleteConnector("local-file-source")
					Expect(err).To(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusConflict))
				})
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns error with a not found response", func() {
				resp, err := client.DeleteConnector("local-file-source")
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})

var _ = Describe("Connector Lifecycle", func() {
	BeforeEach(func() {
		server = ghttp.NewServer()
		client = NewClient(server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("PauseConnector", func() {
		var statusCode int

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/connectors/local-file-source/pause"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWithPtr(&statusCode, nil),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusAccepted
			})

			It("pauses connector", func() {
				resp, err := client.PauseConnector("local-file-source")
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns error with a not found response", func() {
				resp, err := client.PauseConnector("local-file-source")
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("ResumeConnector", func() {
		var statusCode int

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/connectors/local-file-source/resume"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWithPtr(&statusCode, nil),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusAccepted
			})

			It("resumes connector", func() {
				resp, err := client.ResumeConnector("local-file-source")
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns error with a not found response", func() {
				resp, err := client.ResumeConnector("local-file-source")
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("RestartConnector", func() {
		var statusCode int

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/connectors/local-file-source/restart"),
					ghttp.VerifyHeader(jsonAcceptHeader),
					ghttp.RespondWithPtr(&statusCode, nil),
				),
			)
		})

		Context("when existing connector name is given", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
			})

			It("restarts connector", func() {
				resp, err := client.RestartConnector("local-file-source")
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Context("when rebalance is in process", func() {
				BeforeEach(func() {
					statusCode = http.StatusConflict
				})

				It("returns error with a conflict response", func() {
					resp, err := client.RestartConnector("local-file-source")
					Expect(err).To(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusConflict))
				})
			})
		})

		Context("when nonexisting connector name is given", func() {
			BeforeEach(func() {
				// The API actually throws a 500 on POST to nonexistent
				statusCode = http.StatusInternalServerError
			})

			It("returns error with a server error response", func() {
				resp, err := client.RestartConnector("local-file-source")
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
