package wl_integration_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("basic curl functionality", func() {
	It("can get the lists", func() {
		By("Getting lists")
		Eventually(func() error {
			resp, err := client.Curl("GET", "/lists", nil, nil)
			if err != nil {
				return err
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("Not the expected 200 yet")
			}

			return nil
		}).Should(Succeed())
	})
})
