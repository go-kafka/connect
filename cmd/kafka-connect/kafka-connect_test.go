package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"gopkg.in/alecthomas/kingpin.v2"

	. "github.com/go-kafka/connect/cmd/kafka-connect"
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

var _ = Describe("Argument Validation", func() {
	var app *kingpin.Application
	var argv []string
	var subcommand string
	var err error

	JustBeforeEach(func() {
		app = BuildApp()
		subcommand, err = ValidateArgs(app, argv)
	})

	Describe("for global flags", func() {
		Context("with a --host that is not an absolute URL", func() {
			BeforeEach(func() { argv = []string{"--host", "asdfjk", "list"} })

			It("fails", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("host asdfjk is not a valid absolute URL"))
				Expect(subcommand).To(Equal("list"))
			})

			It("does not suggest usage", func() {
				verr, _ := err.(ValidationError)
				Expect(verr.SuggestUsage).To(BeFalse())
			})
		})
	})

	Describe("for a nonexistent command", func() {
		BeforeEach(func() { argv = []string{"asdfjk"} })

		It("fails", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected command but got \"asdfjk\""))
		})
	})

	Describe("for create", func() {
		var existingFilepath string

		BeforeEach(func() { argv = []string{"create"} })

		Context("without a connector name", func() {
			Context("without --from-file", func() {
				It("fails", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("either a connector name or --from-file is required"))
				})

				It("suggests usage", func() {
					verr, _ := err.(ValidationError)
					Expect(verr.SuggestUsage).To(BeTrue())
				})
			})

			Context("with --config", func() {
				BeforeEach(func() {
					tmpfile, _ := ioutil.TempFile("", "connector-config")
					existingFilepath = tmpfile.Name()
					argv = append(argv, "--config", existingFilepath)
				})

				AfterEach(func() {
					os.Remove(existingFilepath)
				})

				It("fails", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("--config requires a connector name"))
				})

				It("suggests usage", func() {
					verr, _ := err.(ValidationError)
					Expect(verr.SuggestUsage).To(BeTrue())
				})
			})
		})

		Context("with a connector name", func() {
			BeforeEach(func() { argv = append(argv, "a-name") })

			Context("without --config", func() {
				It("fails", func() {
					Expect(err).To(HaveOccurred(), "argv: %v", argv)
					Expect(err.Error()).To(Equal("--config is required with a connector name"))
				})

				It("suggests usage", func() {
					verr, _ := err.(ValidationError)
					Expect(verr.SuggestUsage).To(BeTrue())
				})
			})

			Context("with --config", func() {
				BeforeEach(func() {
					tmpfile, _ := ioutil.TempFile("", "connector-config")
					existingFilepath = tmpfile.Name()
					argv = append(argv, "--config", existingFilepath)
				})

				AfterEach(func() {
					os.Remove(existingFilepath)
				})

				Context("and --from-file", func() {
					BeforeEach(func() {
						argv = append(argv, "--from-file", existingFilepath)
					})

					It("fails", func() {
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("--from-file and --config are mutually exclusive"))
					})

					It("suggests usage", func() {
						verr, _ := err.(ValidationError)
						Expect(verr.SuggestUsage).To(BeTrue())
					})
				})
			})
		})
	})

	Describe("for update", func() {
		BeforeEach(func() { argv = []string{"update"} })

		Context("without a connector name", func() {
			It("fails", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("required argument 'name' not provided"))
			})
		})

		Context("with a connector name", func() {
			BeforeEach(func() { argv = append(argv, "a-name") })

			Context("without --config", func() {
				It("fails", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("configuration input is required"))
				})

				It("suggests usage", func() {
					verr, _ := err.(ValidationError)
					Expect(verr.SuggestUsage).To(BeTrue())
				})
			})
		})
	})
})
