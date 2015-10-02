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

var _ = Describe("client - User operations", func() {
	Describe("getting user", func() {
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

			client.User()

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedUser := wl.User{Name: "some user"}
				expectedBody, err := json.Marshal(expectedUser)
				Expect(err).NotTo(HaveOccurred())

				// Unmarshal to ensure exact object is returned
				// avoids odd behavior with the time fields
				err = json.Unmarshal(expectedBody, &expectedUser)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				user, err := client.User()
				Expect(err).NotTo(HaveOccurred())

				Expect(user).To(Equal(expectedUser))
			})

		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.User()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.User()

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

				_, err := client.User()

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

				_, err := client.User()

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

				_, err := client.User()

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("updating user", func() {
		user := wl.User{
			Name:     "username",
			Revision: 12,
		}

		It("performs PUT requests with correct headers to /user", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/user"),
					ghttp.VerifyContentType("application/json"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					// TODO: Add body check here
				),
			)

			client.UpdateUser(wl.User{})

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedUser := wl.User{Name: "some user"}

				// Unmarshal to ensure exact object is returned
				// avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedUser)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedUser)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				user, err := client.UpdateUser(expectedUser)
				Expect(err).NotTo(HaveOccurred())

				Expect(user).To(Equal(expectedUser))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.UpdateUser(user)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.UpdateUser(user)

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

				_, err := client.UpdateUser(user)

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

				_, err := client.UpdateUser(user)

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

				_, err := client.UpdateUser(user)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting users", func() {
		Context("when ListID is not provided", func() {
			It("performs GET requests to /users with correct headers", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/users"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
					),
				)

				client.Users()

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("performs GET requests to /users with correct headers", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/users"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
					),
				)

				client.UsersForListID(listID)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when listID > 0", func() {
			listID := uint(12345)

			It("performs GET requests to /users with correct headers", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/users", "list_id=12345"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
					),
				)

				client.UsersForListID(listID)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedUsers := []wl.User{{Name: "some user"}}

				// Unmarshal to ensure exact object is returned
				// avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedUsers)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedUsers)
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				users, err := client.Users()
				Expect(err).NotTo(HaveOccurred())

				Expect(users).To(Equal(expectedUsers))
			})

		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Users()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Users()

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

				_, err := client.Users()

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

				_, err := client.Users()

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

				_, err := client.Users()

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
