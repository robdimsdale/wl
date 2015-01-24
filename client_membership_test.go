package wundergo_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Client - Membership operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("Getting memberships for list", func() {
		listID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.MembershipsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /memberships?list_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/memberships?list_id=%d", apiUrl, listID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Membership{}, nil)
			client.MembershipsForListID(listID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.MembershipsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.MembershipsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.MembershipsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when reading body returns an error", func() {
			expectedError := errors.New("read error")
			BeforeEach(func() {
				dummyResponse.Body = erroringReadCloser{
					readError: expectedError,
				}
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("forwards the error", func() {
				_, err := client.MembershipsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.MembershipsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedMemberships := &[]wundergo.Membership{
				wundergo.Membership{
					UserID: 1234,
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedMemberships, nil)
			})

			It("returns the unmarshalled notes without error", func() {
				memberships, err := client.MembershipsForListID(listID)

				Expect(err).To(BeNil())
				Expect(memberships).To(Equal(expectedMemberships))
			})
		})
	})

	Describe("Getting memberships", func() {

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /memberships", func() {
			expectedUrl := fmt.Sprintf("%s/memberships", apiUrl)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Membership{}, nil)
			client.Memberships()

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Memberships()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.Memberships()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.Memberships()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when reading body returns an error", func() {
			expectedError := errors.New("read error")
			BeforeEach(func() {
				dummyResponse.Body = erroringReadCloser{
					readError: expectedError,
				}
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("forwards the error", func() {
				_, err := client.Memberships()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Memberships()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedMemberships := &[]wundergo.Membership{
				wundergo.Membership{
					UserID: 1234,
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedMemberships, nil)
			})

			It("returns the unmarshalled memberships without error", func() {
				memberships, err := client.Memberships()

				Expect(err).To(BeNil())
				Expect(memberships).To(Equal(expectedMemberships))
			})
		})
	})

	Describe("getting membership by ID", func() {
		membershipID := uint(1)
		expectedUrl := fmt.Sprintf("%s/memberships/%d", apiUrl, membershipID)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /memberships/:id", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.Membership{}, nil)
			client.Membership(membershipID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Membership(membershipID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.Membership(membershipID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.Membership(membershipID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when reading body returns an error", func() {
			expectedError := errors.New("read error")
			BeforeEach(func() {
				dummyResponse.Body = erroringReadCloser{
					readError: expectedError,
				}
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("forwards the error", func() {
				_, err := client.Membership(membershipID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Membership(membershipID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedMembership := &wundergo.Membership{
				UserID: 1234,
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedMembership, nil)
			})

			It("returns the unmarshalled note without error", func() {
				note, err := client.Membership(membershipID)

				Expect(err).To(BeNil())
				Expect(note).To(Equal(expectedMembership))
			})
		})
	})

	Describe("adding member to a list", func() {
		listID := uint(1)
		muted := true

		Describe("using userID", func() {

			userID := uint(2)

			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusCreated
				fakeHTTPHelper.PostReturns(dummyResponse, nil)
			})

			It("performs POST requests to /memberships?user_id=:userID&list_id=:listID&muted=:muted", func() {
				expectedUrl := fmt.Sprintf("%s/memberships?user_id=%d&list_id=%d&muted=%t", apiUrl, userID, listID, muted)

				fakeJSONHelper.UnmarshalReturns(&wundergo.Membership{}, nil)
				client.AddMemberToListViaUserID(userID, listID, muted)

				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))
			})

			Context("when userID == 0", func() {
				userID := uint(0)

				It("returns an error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when listID == 0", func() {
				listID := uint(0)

				It("returns an error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when httpHelper.Post returns an error", func() {
				expectedError := errors.New("httpHelper POST error")

				BeforeEach(func() {
					fakeHTTPHelper.PostReturns(nil, expectedError)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when response status code is unexpected", func() {
				BeforeEach(func() {
					dummyResponse.StatusCode = http.StatusBadRequest
				})

				It("returns an error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when response body is nil", func() {
				BeforeEach(func() {
					dummyResponse.Body = nil
					fakeHTTPHelper.PostReturns(dummyResponse, nil)
				})

				It("returns an error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when reading body returns an error", func() {
				expectedError := errors.New("read error")
				BeforeEach(func() {
					dummyResponse.Body = erroringReadCloser{
						readError: expectedError,
					}
					fakeHTTPHelper.PostReturns(dummyResponse, nil)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				expectedError := errors.New("jsonHelper error")

				BeforeEach(func() {
					fakeJSONHelper.UnmarshalReturns(nil, expectedError)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when valid response is received", func() {
				expectedMembership := &wundergo.Membership{
					UserID: 1234,
				}

				BeforeEach(func() {
					fakeJSONHelper.UnmarshalReturns(expectedMembership, nil)
				})

				It("returns the unmarshalled note without error", func() {
					note, err := client.AddMemberToListViaUserID(userID, listID, muted)

					Expect(err).To(BeNil())
					Expect(note).To(Equal(expectedMembership))
				})
			})
		})

		Describe("using email address", func() {

			emailAddress := "my-email-address"

			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusCreated
				fakeHTTPHelper.PostReturns(dummyResponse, nil)
			})

			It("performs POST requests to /memberships?email=:emailAddress&list_id=:listID&muted=:muted", func() {
				expectedUrl := fmt.Sprintf("%s/memberships?email=%s&list_id=%d&muted=%t", apiUrl, emailAddress, listID, muted)

				fakeJSONHelper.UnmarshalReturns(&wundergo.Membership{}, nil)
				client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))
			})

			Context("when emailAddress is empty", func() {
				userID := ""

				It("returns an error", func() {
					_, err := client.AddMemberToListViaEmailAddress(userID, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when listID == 0", func() {
				listID := uint(0)

				It("returns an error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when httpHelper.Post returns an error", func() {
				expectedError := errors.New("httpHelper POST error")

				BeforeEach(func() {
					fakeHTTPHelper.PostReturns(nil, expectedError)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when response status code is unexpected", func() {
				BeforeEach(func() {
					dummyResponse.StatusCode = http.StatusBadRequest
				})

				It("returns an error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when response body is nil", func() {
				BeforeEach(func() {
					dummyResponse.Body = nil
					fakeHTTPHelper.PostReturns(dummyResponse, nil)
				})

				It("returns an error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when reading body returns an error", func() {
				expectedError := errors.New("read error")
				BeforeEach(func() {
					dummyResponse.Body = erroringReadCloser{
						readError: expectedError,
					}
					fakeHTTPHelper.PostReturns(dummyResponse, nil)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				expectedError := errors.New("jsonHelper error")

				BeforeEach(func() {
					fakeJSONHelper.UnmarshalReturns(nil, expectedError)
				})

				It("forwards the error", func() {
					_, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when valid response is received", func() {
				expectedMembership := &wundergo.Membership{
					UserID: 1234,
				}

				BeforeEach(func() {
					fakeJSONHelper.UnmarshalReturns(expectedMembership, nil)
				})

				It("returns the unmarshalled note without error", func() {
					note, err := client.AddMemberToListViaEmailAddress(emailAddress, listID, muted)

					Expect(err).To(BeNil())
					Expect(note).To(Equal(expectedMembership))
				})
			})
		})
	})
})
