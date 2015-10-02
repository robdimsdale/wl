package oauth_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wl"
	"github.com/robdimsdale/wl/oauth"
)

var _ = Describe("client - File operations", func() {
	Describe("getting file preview for fileID", func() {
		var (
			fileID   uint
			platform string
			size     string
		)

		BeforeEach(func() {
			fileID = 1234
			platform = ""
			size = ""
		})

		It("performs GET requests with correct headers and query params to /previews", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/previews", "file_id=1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.FilePreview(fileID, platform, size)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when size is provided", func() {
			BeforeEach(func() {
				size = "some-size"
			})

			It("includes size in query params", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/previews", "file_id=1234&size=some-size"),
					),
				)

				client.FilePreview(fileID, platform, size)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when platform is provided", func() {
			BeforeEach(func() {
				platform = "some-platform"
			})

			It("includes platform in query params", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/previews", "file_id=1234&platform=some-platform"),
					),
				)

				client.FilePreview(fileID, platform, size)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFilePreview := wl.FilePreview{URL: "some-url"}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFilePreview)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFilePreview)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				filePreview, err := client.FilePreview(fileID, platform, size)
				Expect(err).NotTo(HaveOccurred())

				Expect(filePreview).To(Equal(expectedFilePreview))
			})
		})

		Context("when fileID == 0", func() {
			BeforeEach(func() {
				fileID = 0
			})

			It("returns an error", func() {
				_, err := client.FilePreview(fileID, platform, size)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.FilePreview(fileID, platform, size)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.FilePreview(fileID, platform, size)

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

				_, err := client.FilePreview(fileID, platform, size)

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

				_, err := client.FilePreview(fileID, platform, size)

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

				_, err := client.FilePreview(fileID, platform, size)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
