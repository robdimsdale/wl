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

var _ = Describe("client - Folder operations", func() {
	Describe("getting folders", func() {
		It("performs GET requests with correct headers to /folders", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/folders"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.Folders()

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFolders := []wundergo.Folder{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFolders)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFolders)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusOK, expectedFolders),
					),
				)

				folders, err := client.Folders()
				Expect(err).NotTo(HaveOccurred())

				Expect(folders).To(Equal(expectedFolders))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", logger)
			})

			It("forwards the error", func() {
				_, err := client.Folders()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", logger)
			})

			It("forwards the error", func() {
				_, err := client.Folders()

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

				_, err := client.Folders()

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

				_, err := client.Folders()

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

				_, err := client.Folders()

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("creating a new folder", func() {
		var (
			title   string
			listIDs []uint
		)

		BeforeEach(func() {
			title = "folder title"
			listIDs = []uint{1234, 5678}
		})

		It("performs POST requests with correct headers to /folders", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/folders"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSON(`{"list_ids":[1234,5678],"title":"folder title"}`),
				),
			)

			client.CreateFolder(
				title,
				listIDs,
			)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFolder := wundergo.Folder{ID: 2345}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFolder)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFolder)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusCreated, expectedFolder),
					),
				)

				folder, err := client.CreateFolder(
					title,
					listIDs,
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(folder).To(Equal(expectedFolder))
			})
		})

		Context("when title is empty", func() {
			BeforeEach(func() {
				title = ""
			})

			It("returns an error", func() {
				_, err := client.CreateFolder(
					title,
					listIDs,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when listIDs is nil", func() {
			BeforeEach(func() {
				listIDs = nil
			})

			It("returns an error", func() {
				_, err := client.CreateFolder(
					title,
					listIDs,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", logger)
			})

			It("forwards the error", func() {
				_, err := client.CreateFolder(
					title,
					listIDs,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", logger)
			})

			It("forwards the error", func() {
				_, err := client.CreateFolder(
					title,
					listIDs,
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

				_, err := client.CreateFolder(
					title,
					listIDs,
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

				_, err := client.CreateFolder(
					title,
					listIDs,
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

				_, err := client.CreateFolder(
					title,
					listIDs,
				)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting folder by ID", func() {
		var folderID uint

		BeforeEach(func() {
			folderID = 1234
		})

		It("performs GET requests with correct headers to /folders", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/folders/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.Folder(folderID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFolder := wundergo.Folder{ID: 2345}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFolder)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFolder)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusOK, expectedFolder),
					),
				)

				folder, err := client.Folder(folderID)
				Expect(err).NotTo(HaveOccurred())

				Expect(folder).To(Equal(expectedFolder))
			})
		})

		Context("when FolderID == 0", func() {
			BeforeEach(func() {
				folderID = 0
			})

			It("returns an error", func() {
				_, err := client.Folder(folderID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", logger)
			})

			It("forwards the error", func() {
				_, err := client.Folder(folderID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", logger)
			})

			It("forwards the error", func() {
				_, err := client.Folder(folderID)

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

				_, err := client.Folder(folderID)

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

				_, err := client.Folder(folderID)

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

				_, err := client.Folder(folderID)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
