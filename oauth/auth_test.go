package oauth_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wl/oauth"
)

var _ = Describe("client - Auth operations", func() {
	Describe("Authed()", func() {
		It("performs GET requests to /user with correct headers", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/user"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			_, _ = client.Authed()

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the response status code is 200", func() {
			It("returns successfully", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)

				authed, err := client.Authed()
				Expect(err).NotTo(HaveOccurred())

				Expect(authed).To(BeTrue())
			})
		})

		Context("when response status code is 401", func() {
			It("returns false without error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusUnauthorized, nil),
					),
				)

				authed, err := client.Authed()
				Expect(err).NotTo(HaveOccurred())

				Expect(authed).To(BeFalse())
			})
		})

		Context("when response status code is 403", func() {
			It("returns false without error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusForbidden, nil),
					),
				)

				authed, err := client.Authed()
				Expect(err).NotTo(HaveOccurred())

				Expect(authed).To(BeFalse())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Authed()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Authed()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response status code is unexpected", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNotFound, nil),
					),
				)

				_, err := client.Authed()

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
