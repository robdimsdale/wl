package wl_integration_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func runMainWithArgs(args ...string) *gexec.Session {
	command := exec.Command(wlBinPath, args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}

var _ = Describe("wl binary", func() {
	var (
		args []string
	)

	BeforeEach(func() {
		args = []string{}
	})

	Describe("Displaying version", func() {
		It("displays version with 'version'", func() {
			args = append(args, fmt.Sprintf("version"))

			command := exec.Command(wlBinPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gbytes.Say("dev"))
			Eventually(session).Should(gexec.Exit(0))
		})
	})

	Context("when provided with valid credentials", func() {
		BeforeEach(func() {
			args = append(args, fmt.Sprintf("--accessToken=%s", wlAccessToken))
			args = append(args, fmt.Sprintf("--clientID=%s", wlClientID))
		})

		It("renders usage if no arguments are provided", func() {
			command := exec.Command(wlBinPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gbytes.Say("Usage:"))
			Eventually(session).Should(gexec.Exit(0))
		})

		Context("with output rendered as json", func() {
			BeforeEach(func() {
				args = append(args, "-j")
			})

			Describe("list operations", func() {
				It("gets all lists", func() {
					args = append(args, fmt.Sprintf("lists"))

					command := exec.Command(wlBinPath, args...)
					session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					Eventually(session).Should(gbytes.Say(`"list_type":"inbox"`))
					Eventually(session).Should(gexec.Exit(0))
				})
			})

			Describe("folder operations", func() {
				It("gets all folders", func() {
					args = append(args, fmt.Sprintf("folders"))

					command := exec.Command(wlBinPath, args...)
					session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())

					// We have no idea how many folders we should expect
					// it could be many, or an empty list
					// Either way we check for '[' and ']' which will always be rendered
					// These characters are escaped to prevent them being interpreted as regex
					Eventually(session).Should(gbytes.Say("\\["))
					Eventually(session).Should(gbytes.Say("\\]"))
					Eventually(session).Should(gexec.Exit(0))
				})
			})
		})
	})
})
