package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func runWundergo(args ...string) *gexec.Session {
	command := exec.Command(wundergoBinPath, args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gbytes.Say("Wundergo"))
	return session
}

var _ = Describe("Main", func() {
	Describe("Basic functionality", func() {
		It("executes succesfully", func() {
			runWundergo()
		})
	})
})
