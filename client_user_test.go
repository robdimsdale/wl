package wundergo_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Client - User operations", func() {

	BeforeEach(func() {
		initializeFakes()
		initializeClient()
	})

	Describe("getting user", func() {
		It("performs GET requests to /user", func() {
			expectedUrl := fmt.Sprintf("%s/user", apiUrl)

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

			It("returns an empty user", func() {
				user, _ := client.User()

				Expect(user).To(Equal(wundergo.User{}))
			})

			It("forwards the error", func() {
				_, err := client.User()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns([]byte("invalid json response"), nil)
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("returns an empty user", func() {
				user, _ := client.User()

				Expect(user).To(Equal(wundergo.User{}))
			})

			It("forwards the error", func() {
				_, err := client.User()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUser := wundergo.User{Name: "testy"}

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns([]byte(""), nil)
				fakeJSONHelper.UnmarshalReturns(&expectedUser, nil)
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
		It("performs PUT requests with new username to /user", func() {
			expectedUrl := fmt.Sprintf("%s/user", apiUrl)

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

			It("returns an empty user", func() {
				user, _ := client.UpdateUser(user)

				Expect(user).To(Equal(wundergo.User{}))
			})

			It("forwards the error", func() {
				_, err := client.UpdateUser(user)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeHTTPHelper.PutReturns([]byte("invalid json response"), nil)
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("returns an empty user", func() {
				user, _ := client.UpdateUser(user)

				Expect(user).To(Equal(wundergo.User{}))
			})

			It("forwards the error", func() {
				_, err := client.UpdateUser(user)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUser := wundergo.User{}

			BeforeEach(func() {
				fakeHTTPHelper.PutReturns([]byte(""), nil)
				fakeJSONHelper.UnmarshalReturns(&expectedUser, nil)
			})

			It("returns the unmarshalled user without error", func() {
				user, err := client.UpdateUser(user)

				Expect(err).To(BeNil())
				Expect(user).To(Equal(expectedUser))
			})
		})
	})

	Describe("getting users", func() {
		Context("when ListID is not provided", func() {
			expectedUrl := fmt.Sprintf("%s/users", apiUrl)

			It("performs GET requests to /users", func() {
				fakeJSONHelper.UnmarshalReturns(&[]wundergo.User{}, nil)
				client.Users()

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})

			Context("when httpHelper.Get returns an error", func() {
				expectedError := errors.New("httpHelper GET error")

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(nil, expectedError)
				})

				It("returns an empty array of users", func() {
					user, _ := client.Users()

					Expect(user).To(Equal([]wundergo.User{}))
				})

				It("forwards the error", func() {
					_, err := client.Users()

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				expectedError := errors.New("jsonHelper error")

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns([]byte("invalid json response"), nil)
					fakeJSONHelper.UnmarshalReturns(nil, expectedError)
				})

				It("returns an empty array of users", func() {
					user, _ := client.Users()

					Expect(user).To(Equal([]wundergo.User{}))
				})

				It("forwards the error", func() {
					_, err := client.Users()

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when valid response is received", func() {
				expectedUsers := []wundergo.User{}

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns([]byte(""), nil)
					fakeJSONHelper.UnmarshalReturns(&expectedUsers, nil)
				})

				It("returns the unmarshalled array of users without error", func() {
					users, err := client.Users()

					Expect(err).To(BeNil())
					Expect(users).To(Equal(expectedUsers))
				})
			})
		})

		Context("when ListID == 0", func() {
			ListID := uint(0)
			expectedUrl := fmt.Sprintf("%s/users", apiUrl)

			It("performs GET requests to /users", func() {
				fakeJSONHelper.UnmarshalReturns(&[]wundergo.User{}, nil)
				client.UsersForListID(ListID)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})

			Context("when httpHelper.Get returns an error", func() {
				expectedError := errors.New("httpHelper GET error")
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(nil, expectedError)
				})

				It("returns an empty array of users", func() {
					user, _ := client.UsersForListID(ListID)

					Expect(user).To(Equal([]wundergo.User{}))
				})

				It("forwards the error", func() {
					_, err := client.UsersForListID(ListID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				expectedError := errors.New("jsonHelper error")

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns([]byte("invalid json response"), nil)
					fakeJSONHelper.UnmarshalReturns(nil, expectedError)
				})

				It("returns an empty array of users", func() {
					user, _ := client.UsersForListID(ListID)

					Expect(user).To(Equal([]wundergo.User{}))
				})

				It("forwards the error", func() {
					_, err := client.UsersForListID(ListID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when valid response is received", func() {
				expectedUsers := []wundergo.User{}

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns([]byte(""), nil)
					fakeJSONHelper.UnmarshalReturns(&expectedUsers, nil)
				})

				It("returns the unmarshalled array of users without error", func() {
					users, err := client.UsersForListID(ListID)

					Expect(err).To(BeNil())
					Expect(users).To(Equal(expectedUsers))
				})
			})
		})

		Context("when ListID > 0", func() {
			ListID := uint(12345)
			expectedUrl := fmt.Sprintf("%s/users?list_id=%d", apiUrl, ListID)

			It("performs GET requests to /users with list_id param", func() {
				fakeJSONHelper.UnmarshalReturns(&[]wundergo.User{}, nil)
				client.UsersForListID(ListID)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})

			Context("when httpHelper.Get returns an error", func() {
				expectedError := errors.New("httpHelper GET error")

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(nil, expectedError)
				})

				It("returns an empty array of users", func() {
					user, _ := client.UsersForListID(ListID)

					Expect(user).To(Equal([]wundergo.User{}))
				})

				It("forwards the error", func() {
					_, err := client.UsersForListID(ListID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				expectedError := errors.New("jsonHelper error")

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns([]byte("invalid json response"), nil)
					fakeJSONHelper.UnmarshalReturns(nil, expectedError)
				})

				It("returns an empty array of users", func() {
					user, _ := client.UsersForListID(ListID)

					Expect(user).To(Equal([]wundergo.User{}))
				})

				It("forwards the error", func() {
					_, err := client.UsersForListID(ListID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when valid response is received", func() {
				expectedUsers := []wundergo.User{}

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns([]byte(""), nil)
					fakeJSONHelper.UnmarshalReturns(&expectedUsers, nil)
				})

				It("returns the unmarshalled array of users without error", func() {
					users, err := client.UsersForListID(ListID)

					Expect(err).To(BeNil())
					Expect(users).To(Equal(expectedUsers))
				})
			})
		})

	})
})
