package oauth_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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

		Context("with valid size and fallback=true", func() {
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

		DescribeTable("size validation",
			func(size int, valid bool) {
				if size <= 0 {
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
				} else if valid {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/avatar", fmt.Sprintf("user_id=1234&size=%d", size)),
							ghttp.VerifyHeader(http.Header{
								"X-Access-Token": []string{dummyAccessToken},
								"X-Client-ID":    []string{dummyClientID},
							}),
						),
					)

					client.AvatarURL(userID, size, fallback)

					Expect(server.ReceivedRequests()).Should(HaveLen(1))
				} else {
					_, err := client.AvatarURL(userID, size, fallback)
					Expect(err).Should(HaveOccurred())
				}
			},
			Entry("negative", -1, false),
			Entry("0", 0, false),
			Entry("25", 25, true),
			Entry("28", 28, true),
			Entry("30", 30, true),
			Entry("32", 32, true),
			Entry("50", 50, true),
			Entry("54", 54, true),
			Entry("56", 56, true),
			Entry("60", 60, true),
			Entry("64", 64, true),
			Entry("108", 108, true),
			Entry("128", 128, true),
			Entry("135", 135, true),
			Entry("256", 256, true),
			Entry("270", 270, true),
			Entry("512", 512, true),
			Entry("513", 513, false),
		)

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
