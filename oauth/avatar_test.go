package oauth_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wl/oauth"
)

var _ = Describe("client - Avatar operations", func() {
	var (
		userID uint
		size   int

		fallback bool
	)

	BeforeEach(func() {
		userID = uint(1234)
		size = 0
		fallback = true
	})

	Describe("getting avatar", func() {
		Context("with size=0 and fallback=true", func() {
			It("performs GET requests to /avatar with correct headers", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/avatar", "user_id=1234"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
					),
				)

				client.AvatarURL(userID, size, fallback)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("with nonzero size and fallback=true", func() {
			BeforeEach(func() {
				size = 128
			})

			It("performs GET requests to /avatar with correct headers and params", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/avatar", "user_id=1234&size=128"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
					),
				)

				client.AvatarURL(userID, size, fallback)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("with size=0 and fallback=true", func() {
			BeforeEach(func() {
				fallback = false
			})

			It("performs GET requests to /avatar with correct headers", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/avatar", "user_id=1234&fallback=false"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
					),
				)

				client.AvatarURL(userID, size, fallback)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedAvatarURL := "https://my-avatar-url"

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusFound, expectedAvatarURL, http.Header{"Location": []string{expectedAvatarURL}}),
					),
				)

				avatarURL, err := client.AvatarURL(userID, size, fallback)
				Expect(err).NotTo(HaveOccurred())

				Expect(avatarURL).To(Equal(expectedAvatarURL))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.AvatarURL(userID, size, fallback)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.AvatarURL(userID, size, fallback)

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

				_, err := client.AvatarURL(userID, size, fallback)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)

				_, err := client.AvatarURL(userID, size, fallback)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, "invalid json response"),
					),
				)

				_, err := client.AvatarURL(userID, size, fallback)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when fallback is false and no fallback avatar is found", func() {
			BeforeEach(func() {
				fallback = false

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNoContent, ""),
					),
				)
			})

			It("returns empty string without error", func() {
				avatarURL, err := client.AvatarURL(userID, size, fallback)

				Expect(err).NotTo(HaveOccurred())
				Expect(avatarURL).To(Equal(""))
			})
		})
	})
})
