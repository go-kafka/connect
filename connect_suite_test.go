package connect_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConnect(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kafka Connect Suite")
}
