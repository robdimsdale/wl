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
				expectedFolders := []wl.Folder{{ID: 2345}}

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
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Folders()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
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
				expectedFolder := wl.Folder{ID: 2345}

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

		Context("when listIDs is empty", func() {
			BeforeEach(func() {
				listIDs = []uint{}
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
				client = oauth.NewClient("", "", "", testLogger)
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
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
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
				expectedFolder := wl.Folder{ID: 2345}

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
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Folder(folderID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
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

	Describe("updating a folder", func() {
		var folder wl.Folder

		BeforeEach(func() {
			folder = wl.Folder{ID: 1234}
		})

		It("performs PUT requests with correct headers to /folders/1234", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "/folders/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.UpdateFolder(folder)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFolder := wl.Folder{ID: 2345}

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

				folder, err := client.UpdateFolder(folder)
				Expect(err).NotTo(HaveOccurred())

				Expect(folder).To(Equal(expectedFolder))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.UpdateFolder(folder)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.UpdateFolder(folder)

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

				_, err := client.UpdateFolder(folder)

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

				_, err := client.UpdateFolder(folder)

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

				_, err := client.UpdateFolder(folder)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("deleting a folder", func() {
		var folder wl.Folder

		BeforeEach(func() {
			folder = wl.Folder{ID: 1234, Revision: 23}
		})

		It("performs DELETE requests with correct headers to /folders/1234", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/folders/1234", "revision=23"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.DeleteFolder(folder)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusNoContent, nil),
					),
				)

				err := client.DeleteFolder(folder)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				err := client.DeleteFolder(folder)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				err := client.DeleteFolder(folder)

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

				err := client.DeleteFolder(folder)

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

				err := client.DeleteFolder(folder)

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

				err := client.DeleteFolder(folder)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting folder revisions", func() {
		It("performs GET requests with correct headers to /folder_revisions", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/folder_revisions"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.FolderRevisions()

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFolderRevisions := []wl.FolderRevision{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFolderRevisions)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFolderRevisions)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusOK, expectedFolderRevisions),
					),
				)

				folderRevisions, err := client.FolderRevisions()
				Expect(err).NotTo(HaveOccurred())

				Expect(folderRevisions).To(Equal(expectedFolderRevisions))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.FolderRevisions()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.FolderRevisions()

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

				_, err := client.FolderRevisions()

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

				_, err := client.FolderRevisions()

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

				_, err := client.FolderRevisions()

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
