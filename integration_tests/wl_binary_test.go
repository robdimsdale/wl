package wundergo_integration_test

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

		It("displays version with '-v'", func() {
			args = append(args, fmt.Sprintf("-v"))

			command := exec.Command(wlBinPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gbytes.Say("dev"))
			Eventually(session).Should(gexec.Exit(0))
		})

		It("displays version with '--version'", func() {
			args = append(args, fmt.Sprintf("--version"))

			command := exec.Command(wlBinPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gbytes.Say("dev"))
			Eventually(session).Should(gexec.Exit(0))
		})
	})

	Context("when provided with valid credentials", func() {
		BeforeEach(func() {
			args = append(args, fmt.Sprintf("-accessToken=%s", wlAccessToken))
			args = append(args, fmt.Sprintf("-clientID=%s", wlClientID))
		})

		It("exits with failure code if no arguments are provided", func() {
			command := exec.Command(wlBinPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(2))
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
	})
})
