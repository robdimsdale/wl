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

var _ = Describe("Client - User operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("getting user", func() {

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /user", func() {
			expectedUrl := fmt.Sprintf("%s/user", apiURL)

			fakeJSONHelper.UnmarshalReturns(&wundergo.User{}, nil)
			client.User()

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.User()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.User()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.User()

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
				_, err := client.User()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.User()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUser := &wundergo.User{
				Name: "testy",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedUser, nil)
			})

			It("returns the unmarshalled user without error", func() {
				user, err := client.User()

				Expect(err).To(BeNil())
				Expect(user).To(Equal(expectedUser))
			})
		})
	})

	Describe("updating user", func() {
		user := wundergo.User{
			Name:     "username",
			Revision: 12,
		}

		BeforeEach(func() {
			fakeHTTPHelper.PutReturns(dummyResponse, nil)
		})

		It("performs PUT requests with new username to /user", func() {
			expectedUrl := fmt.Sprintf("%s/user", apiURL)

			fakeJSONHelper.UnmarshalReturns(&wundergo.User{}, nil)
			expectedBody := []byte(fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name))

			client.UpdateUser(user)

			Expect(fakeHTTPHelper.PutCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PutArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when httpHelper.Put returns an error", func() {
			expectedError := errors.New("httpHelper PUT error")

			BeforeEach(func() {
				fakeHTTPHelper.PutReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateUser(user)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PutReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.UpdateUser(user)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when reading body returns an error", func() {
			expectedError := errors.New("read error")
			BeforeEach(func() {
				dummyResponse.Body = erroringReadCloser{
					readError: expectedError,
				}
				fakeHTTPHelper.PutReturns(dummyResponse, nil)
			})

			It("forwards the error", func() {
				_, err := client.UpdateUser(user)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateUser(user)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUser := &wundergo.User{
				Name: "Testy",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedUser, nil)
			})

			It("returns the unmarshalled user without error", func() {
				user, err := client.UpdateUser(user)

				Expect(err).To(BeNil())
				Expect(user).To(Equal(expectedUser))
			})
		})
	})

	Describe("getting users", func() {

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID is not provided", func() {
			expectedUrl := fmt.Sprintf("%s/users", apiURL)

			It("performs GET requests to /users", func() {
				fakeJSONHelper.UnmarshalReturns(&[]wundergo.User{}, nil)
				client.Users()

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})
		})

		Context("when ListID == 0", func() {
			listID := uint(0)
			expectedUrl := fmt.Sprintf("%s/users", apiURL)

			It("performs GET requests to /users", func() {
				fakeJSONHelper.UnmarshalReturns(&[]wundergo.User{}, nil)
				client.UsersForListID(listID)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})
		})

		Context("when listID > 0", func() {
			listID := uint(12345)
			expectedUrl := fmt.Sprintf("%s/users?list_id=%d", apiURL, listID)

			It("performs GET requests to /users with list_id param", func() {
				fakeJSONHelper.UnmarshalReturns(&[]wundergo.User{}, nil)
				client.UsersForListID(listID)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})

			Context("when httpHelper.Get returns an error", func() {
				expectedError := errors.New("httpHelper GET error")

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(nil, expectedError)
				})

				It("forwards the error", func() {
					_, err := client.UsersForListID(listID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when response status code is unexpected", func() {
				BeforeEach(func() {
					dummyResponse.StatusCode = http.StatusBadRequest
				})

				It("returns an error", func() {
					_, err := client.UsersForListID(listID)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when response body is nil", func() {
				BeforeEach(func() {
					dummyResponse.Body = nil
					fakeHTTPHelper.GetReturns(dummyResponse, nil)
				})

				It("returns an error", func() {
					_, err := client.UsersForListID(listID)

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
					_, err := client.UsersForListID(listID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				expectedError := errors.New("jsonHelper error")

				BeforeEach(func() {
					fakeJSONHelper.UnmarshalReturns(nil, expectedError)
				})

				It("forwards the error", func() {
					_, err := client.UsersForListID(listID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when valid response is received", func() {
				expectedUsers := &[]wundergo.User{
					wundergo.User{
						Name: "Testy",
					},
				}

				BeforeEach(func() {
					fakeJSONHelper.UnmarshalReturns(expectedUsers, nil)
				})

				It("returns the unmarshalled array of users without error", func() {
					users, err := client.UsersForListID(listID)

					Expect(err).To(BeNil())
					Expect(users).To(Equal(expectedUsers))
				})
			})
		})
	})
})
