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

var _ = Describe("client - TaskComment operations", func() {
	Describe("getting task comments for list", func() {
		var (
			listID uint
		)

		BeforeEach(func() {
			listID = 1234
		})

		It("performs GET requests with correct headers to /task_comments", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/task_comments", "list_id=1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.TaskCommentsForListID(listID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedTaskComments := []wl.TaskComment{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedTaskComments)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedTaskComments)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)
				taskComments, err := client.TaskCommentsForListID(listID)
				Expect(err).NotTo(HaveOccurred())

				Expect(taskComments).To(Equal(expectedTaskComments))
			})
		})

		Context("when ListID == 0", func() {
			BeforeEach(func() {
				listID = 0
			})

			It("returns an error", func() {
				_, err := client.TaskCommentsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.TaskCommentsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.TaskCommentsForListID(listID)

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

				_, err := client.TaskCommentsForListID(listID)

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

				_, err := client.TaskCommentsForListID(listID)

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

				_, err := client.TaskCommentsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting task comments for task", func() {
		var (
			taskID uint
		)

		BeforeEach(func() {
			taskID = 1234
		})

		It("performs GET requests with correct headers to /task_comments", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/task_comments", "task_id=1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.TaskCommentsForTaskID(taskID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedTaskComments := []wl.TaskComment{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedTaskComments)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedTaskComments)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)
				taskComments, err := client.TaskCommentsForTaskID(taskID)
				Expect(err).NotTo(HaveOccurred())

				Expect(taskComments).To(Equal(expectedTaskComments))
			})
		})

		Context("when TaskID == 0", func() {
			BeforeEach(func() {
				taskID = 0
			})

			It("returns an error", func() {
				_, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.TaskCommentsForTaskID(taskID)

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

				_, err := client.TaskCommentsForTaskID(taskID)

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

				_, err := client.TaskCommentsForTaskID(taskID)

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

				_, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("creating a new task comment", func() {
		var (
			taskID uint
			text   string
		)

		BeforeEach(func() {
			taskID = 1234
			text = "some text"
		})

		It("performs GET requests with correct headers to /task_comments", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/task_comments"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSON(`{"task_id":1234,"text":"some text"}`),
				),
			)

			client.CreateTaskComment(text, taskID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedTaskComment := wl.TaskComment{ID: 1234}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedTaskComment)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedTaskComment)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusCreated, expectedBody),
					),
				)

				taskComment, err := client.CreateTaskComment(text, taskID)
				Expect(err).NotTo(HaveOccurred())

				Expect(taskComment).To(Equal(expectedTaskComment))
			})
		})

		Context("when taskCommentID == 0", func() {
			BeforeEach(func() {
				taskID = 0
			})

			It("returns an error", func() {
				_, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.CreateTaskComment(text, taskID)

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

				_, err := client.CreateTaskComment(text, taskID)

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

				_, err := client.CreateTaskComment(text, taskID)

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

				_, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting task comment by ID", func() {
		var (
			taskCommentID uint
		)

		BeforeEach(func() {
			taskCommentID = 1234
		})

		It("performs GET requests with correct headers to /task_comments/:id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/task_comments/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.TaskComment(taskCommentID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedTaskComment := wl.TaskComment{ID: 2345}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedTaskComment)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedTaskComment)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				taskComment, err := client.TaskComment(taskCommentID)
				Expect(err).NotTo(HaveOccurred())

				Expect(taskComment).To(Equal(expectedTaskComment))
			})
		})

		Context("when taskCommentID == 0", func() {
			BeforeEach(func() {
				taskCommentID = 0
			})

			It("returns an error", func() {
				_, err := client.TaskComment(taskCommentID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.TaskComment(taskCommentID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.TaskComment(taskCommentID)

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

				_, err := client.TaskComment(taskCommentID)

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

				_, err := client.TaskComment(taskCommentID)

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

				_, err := client.TaskComment(taskCommentID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("deleting a task comment", func() {
		var (
			taskComment wl.TaskComment
		)

		BeforeEach(func() {
			taskComment = wl.TaskComment{ID: 1234, Revision: 23}
		})

		It("performs DELETE requests with correct headers to /task_comments", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/task_comments/1234", "revision=23"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.DeleteTaskComment(taskComment)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNoContent, nil),
					),
				)

				err := client.DeleteTaskComment(taskComment)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				err := client.DeleteTaskComment(taskComment)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				err := client.DeleteTaskComment(taskComment)

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

				err := client.DeleteTaskComment(taskComment)

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

				err := client.DeleteTaskComment(taskComment)

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

				err := client.DeleteTaskComment(taskComment)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
