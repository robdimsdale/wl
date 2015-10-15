package wl_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic root functionality", func() {
	It("gets root correctly", func() {
		var err error
		var root wl.Root
		Eventually(func() error {
			root, err = client.Root()
			return err
		}).Should(Succeed())

		Expect(root.ID).To(BeNumerically(">", 0))
	})
})
