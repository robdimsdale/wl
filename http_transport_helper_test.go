package wundergo_test

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Default HTTP Transport Helper", func() {
	var transportHelper *wundergo.DefaultHTTPTransportHelper

	BeforeEach(func() {
		transportHelper = &wundergo.DefaultHTTPTransportHelper{}
	})

	Describe("NewRequest", func() {
		It("Returns a new http.Request without error", func() {
			method := "aMethod"
			urlStr := "aUrl"
			body := strings.NewReader("someBody")

			expectedRequest, err := http.NewRequest(method, urlStr, body)
			Expect(err).NotTo(HaveOccurred())

			request, err := transportHelper.NewRequest(method, urlStr, body)

			Expect(err).NotTo(HaveOccurred())
			Expect(request).To(Equal(expectedRequest))
		})
	})
})
