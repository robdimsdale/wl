package oauth_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wl"
	"github.com/robdimsdale/wl/oauth"
)

var _ = Describe("client - Membership operations", func() {
	Describe("getting memberships for list", func() {
		var listID uint

		BeforeEach(func() {
			listID = 1234
		})

		It("performs GET requests with correct headers to /memberships", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/memberships", "list_id=1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.MembershipsForListID(listID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedMemberships := []wl.Membership{{ID: 2345}}

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusOK, expectedMemberships),
					),
				)

				memberships, err := client.MembershipsForListID(listID)
				Expect(err).NotTo(HaveOccurred())

				Expect(memberships).To(Equal(expectedMemberships))
			})
		})

		Context("when ListID == 0", func() {
			BeforeEach(func() {
				listID = 0
			})

			It("returns an error", func() {
				_, err := client.MembershipsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.MembershipsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.MembershipsForListID(listID)

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

				_, err := client.MembershipsForListID(listID)

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

				_, err := client.MembershipsForListID(listID)

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

				_, err := client.MembershipsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting memberships", func() {
		It("performs GET requests with correct headers to /memberships", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/memberships"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.Memberships()

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedMemberships := []wl.Membership{{ID: 2345}}

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusOK, expectedMemberships),
					),
				)

				memberships, err := client.Memberships()
				Expect(err).NotTo(HaveOccurred())

				Expect(memberships).To(Equal(expectedMemberships))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Memberships()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Memberships()

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

				_, err := client.Memberships()

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

				_, err := client.Memberships()

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

				_, err := client.Memberships()

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting membership by ID", func() {
		var id uint

		BeforeEach(func() {
			id = 1234
		})

		It("performs GET requests with correct headers to /memberships/:id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/memberships/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.Membership(id)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedMembership := wl.Membership{ID: id}

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusOK, expectedMembership),
					),
				)

				membership, err := client.Membership(id)
				Expect(err).NotTo(HaveOccurred())

				Expect(membership).To(Equal(expectedMembership))
			})
		})

		Context("when id == 0", func() {
			BeforeEach(func() {
				id = 0
			})

			It("returns an error", func() {
				_, err := client.Membership(id)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Membership(id)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Membership(id)

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

				_, err := client.Membership(id)

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

				_, err := client.Membership(id)

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

				_, err := client.Membership(id)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("adding member to a list", func() {
		var (
			muted  bool
			listID uint
		)

		BeforeEach(func() {
			listID = 2345
			muted = true
		})

		Describe("using userID", func() {
			var userID uint

			BeforeEach(func() {
				userID = 1234
			})

			It("performs POST requests with correct headers to /memberships/:userID", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/memberships"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
						ghttp.VerifyJSON(`{"user_id":1234,"list_id":2345,"muted":true}`),
					),
				)

				client.AddMemberToListViaUserID(userID, listID, muted)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})

			Context("when the request is valid", func() {
				It("returns successfully", func() {
					expectedMembership := wl.Membership{ID: 3456}

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.RespondWithJSONEncoded(http.StatusCreated, expectedMembership),
						),
					)

					membership, err := client.AddMemberToListViaUserID(userID, listID, muted)
					Expect(err).NotTo(HaveOccurred())

					Expect(membership).To(Equal(expectedMembership))
				})
			})

			Context("when userID == 0", func() {
				BeforeEach(func() {
					userID = 0
				})

				It("returns an error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when listID == 0", func() {
				BeforeEach(func() {
					listID = 0
				})

				It("returns an error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when creating request fails with error", func() {
				BeforeEach(func() {
					client = oauth.NewClient("", "", "", testLogger)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when executing request fails with error", func() {
				BeforeEach(func() {
					client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

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

					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

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

					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

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

					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("using email address", func() {
			var emailAddress string

			BeforeEach(func() {
				emailAddress = "my email address"
			})

			It("performs POST requests with correct headers to /memberships/:id", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/memberships"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
						ghttp.VerifyJSON(`{"email":"my email address","list_id":2345,"muted":true}`),
					),
				)

				client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})

			Context("when the request is valid", func() {
				It("returns successfully", func() {
					expectedMembership := wl.Membership{ID: 3456}

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.RespondWithJSONEncoded(http.StatusCreated, expectedMembership),
						),
					)

					membership, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)
					Expect(err).NotTo(HaveOccurred())

					Expect(membership).To(Equal(expectedMembership))
				})
			})

			Context("when emailAddress is empty", func() {
				BeforeEach(func() {
					emailAddress = ""
				})

				It("returns an error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when creating request fails with error", func() {
				BeforeEach(func() {
					client = oauth.NewClient("", "", "", testLogger)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when executing request fails with error", func() {
				BeforeEach(func() {
					client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

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

					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

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

					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

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

					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("marking member as accepted", func() {
		var (
			membership wl.Membership
		)

		BeforeEach(func() {
			membership = wl.Membership{ID: 1234}
		})

		It("performs POST requests with correct headers to /memberships/:userID", func() {
			expectedMembership := wl.Membership{
				ID:    membership.ID,
				State: "accepted",
			}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "/memberships/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSONRepresenting(expectedMembership),
				),
			)

			client.AcceptMember(membership)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedMembership := wl.Membership{
					ID: membership.ID,
				}

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusOK, &expectedMembership),
					),
				)

				actualMembership, err := client.AcceptMember(membership)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualMembership).To(Equal(expectedMembership))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.AcceptMember(membership)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.AcceptMember(membership)

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

				_, err := client.AcceptMember(membership)

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

				_, err := client.AcceptMember(membership)

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

				_, err := client.AcceptMember(membership)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
