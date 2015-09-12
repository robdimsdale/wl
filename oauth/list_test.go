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

var _ = Describe("client - List operations", func() {
	BeforeEach(func() {
		oauth.NewLogger = func() wundergo.Logger {
			return &fakeLogger
		}

		client = oauth.NewClient(
			dummyAccessToken,
			dummyClientID,
			apiURL,
		)
	})

	Describe("getting lists", func() {
		It("performs GET requests with correct headers to /lists", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/lists"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.Lists()

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedLists := &[]wundergo.List{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedLists)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, expectedLists)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				lists, err := client.Lists()
				Expect(err).NotTo(HaveOccurred())

				Expect(*lists).To(Equal(*expectedLists))
			})
		})

		Context("when creating request fails with error", func() {
			client := oauth.NewClient("", "", "")

			It("forwards the error", func() {
				_, err := client.Lists()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			client := oauth.NewClient("", "", "http://not-a-real-url.com")

			It("forwards the error", func() {
				_, err := client.Lists()

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

				_, err := client.Lists()

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

				_, err := client.Lists()

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

				_, err := client.Lists()

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting list by ID", func() {
		var listID uint

		BeforeEach(func() {
			listID = 1234
		})

		It("performs GET requests with correct headers to /list/:id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/lists/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.List(listID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedList := &wundergo.List{
					ID: listID,
				}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedList)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, expectedList)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				list, err := client.List(listID)
				Expect(err).NotTo(HaveOccurred())

				Expect(*list).To(Equal(*expectedList))
			})
		})

		Context("when creating request fails with error", func() {
			client := oauth.NewClient("", "", "")

			It("forwards the error", func() {
				_, err := client.List(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			client := oauth.NewClient("", "", "http://not-a-real-url.com")

			It("forwards the error", func() {
				_, err := client.List(listID)

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

				_, err := client.List(listID)

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

				_, err := client.List(listID)

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

				_, err := client.List(listID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("creating a new list", func() {
		var title string

		BeforeEach(func() {
			title = "a list"
		})

		It("performs POST requests with correct headers to /lists", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/lists"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSON(`{"title":"a list"}`),
				),
			)

			client.CreateList(title)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedList := &wundergo.List{ID: 2345}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedList)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, expectedList)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusCreated, expectedBody),
					),
				)

				actualList, err := client.CreateList(title)
				Expect(err).NotTo(HaveOccurred())

				Expect(*actualList).To(Equal(*expectedList))
			})
		})

		Context("when the title is empty", func() {
			BeforeEach(func() {
				title = ""
			})

			It("returns an error", func() {
				_, err := client.CreateList(title)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			client := oauth.NewClient("", "", "")

			It("returns an error", func() {
				_, err := client.CreateList(title)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			client := oauth.NewClient("", "", "http://not-a-real-url.com")

			It("forwards the error", func() {
				_, err := client.CreateList(title)

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

				_, err := client.CreateList(title)

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

				_, err := client.CreateList(title)

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

				_, err := client.CreateList(title)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("updating a list", func() {
		var list *wundergo.List

		BeforeEach(func() {
			list = &wundergo.List{
				ID: 1234,
			}
		})

		It("performs PATCH requests with correct headers to /lists/:id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "/lists/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSONRepresenting(list),
				),
			)

			client.UpdateList(*list)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedList := &wundergo.List{ID: list.ID}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedList)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, expectedList)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				actualList, err := client.UpdateList(*list)
				Expect(err).NotTo(HaveOccurred())

				Expect(*actualList).To(Equal(*expectedList))
			})
		})

		Context("when creating request fails with error", func() {
			client := oauth.NewClient("", "", "")

			It("forwards the error", func() {
				_, err := client.UpdateList(*list)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			client := oauth.NewClient("", "", "http://not-a-real-url.com")

			It("forwards the error", func() {
				_, err := client.UpdateList(*list)

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

				_, err := client.UpdateList(*list)

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

				_, err := client.UpdateList(*list)

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

				_, err := client.UpdateList(*list)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("deleting a list", func() {
		var list *wundergo.List

		BeforeEach(func() {
			list = &wundergo.List{
				ID: 1234,
			}
		})

		It("performs DELETE requests with correct headers to /lists/:id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/lists/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.DeleteList(*list)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNoContent, nil),
					),
				)

				err := client.DeleteList(*list)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			client := oauth.NewClient("", "", "")

			It("forwards the error", func() {
				err := client.DeleteList(*list)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			client := oauth.NewClient("", "", "http://not-a-real-url.com")

			It("forwards the error", func() {
				err := client.DeleteList(*list)

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

				err := client.DeleteList(*list)

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

				err := client.DeleteList(*list)

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

				err := client.DeleteList(*list)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
