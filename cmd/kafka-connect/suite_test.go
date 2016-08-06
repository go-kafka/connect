package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var pathToCLI string

func TestKafkaConnectCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-kafka/connect CLI Suite")
}

var _ = BeforeSuite(func() {
	var err error

	// Build the executable in a sandbox
	pathToCLI, err = gexec.Build("github.com/go-kafka/connect/cmd/kafka-connect")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
