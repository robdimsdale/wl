package wundergo_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("basic root functionality", func() {
	It("gets root correctly", func() {
		root, err := client.Root()
		Expect(err).NotTo(HaveOccurred())

		Expect(root.ID).To(BeNumerically(">", 0))
	})
})
