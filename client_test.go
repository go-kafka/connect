package connect_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/go-kafka/connect"
)

var _ = Describe("NewClient", func() {
	It("uses the default host URL", func() {
		client := NewClient()
		Expect(client.Host()).To(Equal(DefaultHostURL))
	})

	It("uses a given host URL", func() {
		host := "http://example.com"
		client := NewClient(host)
		Expect(client.Host()).To(Equal(host))
	})

	Context("given an invalid host URL", func() {
		initClient := func() {
			NewClient("&*%$fdasj")
		}

		It("panics", func() {
			Expect(initClient).To(Panic())
		})
	})

	Context("given multiple host arguments", func() {
		initClient := func() {
			NewClient("one", "another")
		}

		It("panics", func() {
			Expect(initClient).To(Panic())
		})
	})
})
