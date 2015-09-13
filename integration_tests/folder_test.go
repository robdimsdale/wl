package wundergo_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("basic webhook functionality", func() {

	It("lists folders", func() {
		var err error

		By("listing folders")
		_, err = client.Folders()
		Expect(err).NotTo(HaveOccurred())
	})
})
