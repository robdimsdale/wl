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
	Describe("getting files for list", func() {
		var (
			listID uint
		)

		BeforeEach(func() {
			listID = 1234
		})

		It("performs GET requests with correct headers to /files", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/files", "list_id=1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.FilesForListID(listID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFiles := []wl.File{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFiles)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFiles)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)
				files, err := client.FilesForListID(listID)
				Expect(err).NotTo(HaveOccurred())

				Expect(files).To(Equal(expectedFiles))
			})
		})

		Context("when ListID == 0", func() {
			BeforeEach(func() {
				listID = 0
			})

			It("returns an error", func() {
				_, err := client.FilesForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.FilesForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.FilesForListID(listID)

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

				_, err := client.FilesForListID(listID)

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

				_, err := client.FilesForListID(listID)

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

				_, err := client.FilesForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting files for task", func() {
		var (
			taskID uint
		)

		BeforeEach(func() {
			taskID = 1234
		})

		It("performs GET requests with correct headers to /files", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/files", "task_id=1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.FilesForTaskID(taskID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFiles := []wl.File{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFiles)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFiles)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)
				files, err := client.FilesForTaskID(taskID)
				Expect(err).NotTo(HaveOccurred())

				Expect(files).To(Equal(expectedFiles))
			})
		})

		Context("when TaskID == 0", func() {
			BeforeEach(func() {
				taskID = 0
			})

			It("returns an error", func() {
				_, err := client.FilesForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.FilesForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.FilesForTaskID(taskID)

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

				_, err := client.FilesForTaskID(taskID)

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

				_, err := client.FilesForTaskID(taskID)

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

				_, err := client.FilesForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("creating a new file", func() {
		var (
			taskID   uint
			uploadID uint
		)

		BeforeEach(func() {
			taskID = 1234
			uploadID = 2345
		})

		It("performs GET requests with correct headers to /files", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/files"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSON(`{"task_id":1234,"upload_id":2345}`),
				),
			)

			client.CreateFile(uploadID, taskID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFile := wl.File{ID: 1234}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFile)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFile)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusCreated, expectedBody),
					),
				)

				file, err := client.CreateFile(uploadID, taskID)
				Expect(err).NotTo(HaveOccurred())

				Expect(file).To(Equal(expectedFile))
			})
		})

		Context("when fileID == 0", func() {
			BeforeEach(func() {
				taskID = 0
			})

			It("returns an error", func() {
				_, err := client.CreateFile(uploadID, taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.CreateFile(uploadID, taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.CreateFile(uploadID, taskID)

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

				_, err := client.CreateFile(uploadID, taskID)

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

				_, err := client.CreateFile(uploadID, taskID)

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

				_, err := client.CreateFile(uploadID, taskID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting file by ID", func() {
		var (
			fileID uint
		)

		BeforeEach(func() {
			fileID = 1234
		})

		It("performs GET requests with correct headers to /files/:id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/files/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.File(fileID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedFile := wl.File{ID: 2345}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedFile)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedFile)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				file, err := client.File(fileID)
				Expect(err).NotTo(HaveOccurred())

				Expect(file).To(Equal(expectedFile))
			})
		})

		Context("when fileID == 0", func() {
			BeforeEach(func() {
				fileID = 0
			})

			It("returns an error", func() {
				_, err := client.File(fileID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.File(fileID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.File(fileID)

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

				_, err := client.File(fileID)

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

				_, err := client.File(fileID)

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

				_, err := client.File(fileID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("deleting a file", func() {
		var (
			file wl.File
		)

		BeforeEach(func() {
			file = wl.File{ID: 1234, Revision: 23}
		})

		It("performs DELETE requests with correct headers to /files", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/files/1234", "revision=23"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.DestroyFile(file)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNoContent, nil),
					),
				)

				err := client.DestroyFile(file)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				err := client.DestroyFile(file)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				err := client.DestroyFile(file)

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

				err := client.DestroyFile(file)

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

				err := client.DestroyFile(file)

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

				err := client.DestroyFile(file)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
