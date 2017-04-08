package oauth_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wl/oauth"
)

var _ = Describe("client - Curl operations", func() {
	Describe("curling a URL", func() {
		var (
			method string
			url    string
			body   []byte
		)

		BeforeEach(func() {
			method = "PATCH"
			url = "/path/to/some/url"
			body = []byte("some body")
		})

		It("performs requests with correct headers to provided url with body", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(method, url),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyBody(body),
				),
			)

			client.Curl(method, url, body)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("strips preceding slashes", func() {
			url = fmt.Sprintf("//%s", url)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(method, url),
				),
			)

			client.Curl(method, url, body)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedResponse := []byte("expected-response")

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedResponse),
					),
				)
				returned, err := client.Curl(method, url, body)
				Expect(err).NotTo(HaveOccurred())

				Expect(returned.StatusCode).To(Equal(http.StatusOK))

				returnedBody, err := ioutil.ReadAll(returned.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(returnedBody).To(Equal(expectedResponse))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Curl(method, url, body)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Curl(method, url, body)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
