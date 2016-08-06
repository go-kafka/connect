package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("kafka-connect CLI", func() {
	var command *exec.Cmd

	Context("with no arguments", func() {
		BeforeEach(func() {
			command = exec.Command(pathToCLI)
		})

		It("executes successfully", func() {
			session, err := Start(command, nil, nil) // Don't need to see this output with -v
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(Exit(0))
		})

		It("outputs usage help", func() {
			session, _ := Start(command, nil, nil)
			Eventually(session).Should(Say("usage: kafka-connect"))
		})
	})
})
