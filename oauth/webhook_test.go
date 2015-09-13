package oauth_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/oauth"
)

var _ = Describe("client - Webhook operations", func() {

	Describe("creating a new webhook", func() {
		var (
			url           string
			listID        uint
			processorType string
			configuration string
		)

		BeforeEach(func() {
			url = "some-url"
			listID = 1234
			processorType = "generic"
			configuration = "some configuration"
		})

		It("performs POST requests with correct headers to /webhooks", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/webhooks"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSON(`{"list_id":1234,"url":"some-url","processor_type":"generic","configuration":"some configuration"}`),
				),
			)

			client.CreateWebhook(
				listID,
				url,
				processorType,
				configuration,
			)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedWebhook := wundergo.Webhook{ID: 2345}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedWebhook)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedWebhook)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusCreated, expectedWebhook),
					),
				)

				note, err := client.CreateWebhook(
					listID,
					url,
					processorType,
					configuration,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(note).To(Equal(expectedWebhook))
			})
		})

		Context("when NoteID == 0", func() {
			BeforeEach(func() {
				listID = 0
			})

			It("returns an error", func() {
				_, err := client.CreateWebhook(
					listID,
					url,
					processorType,
					configuration,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			client := oauth.NewClient("", "", "", logger)

			It("forwards the error", func() {
				_, err := client.CreateWebhook(
					listID,
					url,
					processorType,
					configuration,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			client := oauth.NewClient("", "", "http://not-a-real-url.com", logger)

			It("forwards the error", func() {
				_, err := client.CreateWebhook(
					listID,
					url,
					processorType,
					configuration,
				)

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

				_, err := client.CreateWebhook(
					listID,
					url,
					processorType,
					configuration,
				)

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

				_, err := client.CreateWebhook(
					listID,
					url,
					processorType,
					configuration,
				)

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

				_, err := client.CreateWebhook(
					listID,
					url,
					processorType,
					configuration,
				)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
