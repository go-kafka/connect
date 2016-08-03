package connect_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/go-kafka/connect"
)

var _ = Describe("Connectors", func() {
	Describe("CreateConnector", func() {
		It("creates a new instance given a valid Connector", func() {
			connector := Connector{Name: "test"}
			err := CreateConnector(connector)
			Expect(err).NotTo(HaveOccurred())
		})

		PIt("returns an error given an invalid Connector", func() {
		})
	})

	Describe("GetConnectors", func() {
		It("returns list of connector names", func() {
			names := GetConnectors()
			Expect(names).To(HaveLen(0))
		})
	})

	Describe("GetConnector", func() {
		PIt("returns a connector instance by name", func() {
			connector, err := GetConnector("test")
			Expect(err).NotTo(HaveOccurred())
			Expect(connector.Name).To(Equal("test"))
		})
	})

	Describe("GetConnectorConfig", func() {
		It("returns configuration for an extant connector", func() {
			config, err := GetConnectorConfig("test")
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil()) // TODO: assert properties of the config
		})

		PIt("returns error given a nonexistent connector", func() {
			_, err := GetConnectorConfig("invalid")
			Expect(err).To(HaveOccurred()) // TODO: assert properties of error
		})
	})

	Describe("GetConnectorTasks", func() {
		It("returns tasks for an extant connector", func() {
			tasks, err := GetConnectorTasks("test")
			Expect(err).NotTo(HaveOccurred())
			Expect(tasks).To(HaveLen(0)) // TODO: assert task data
		})

		PIt("returns error given a nonexistent connector", func() {
			_, err := GetConnectorTasks("invalid")
			Expect(err).To(HaveOccurred()) // TODO: assert properties of error
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
