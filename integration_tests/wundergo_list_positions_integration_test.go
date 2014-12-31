package wundergo_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Basic list position functionality", func() {
	It("gets list positions", func() {
		var err error
		var listPositions *[]wundergo.Position
		Eventually(func() error {
			listPositions, err = client.ListPositions()
			return err
		}).Should(Succeed())
		Expect(listPositions).ToNot(BeNil())
	})
})
