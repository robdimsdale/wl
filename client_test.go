package wundergo_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/fakes"
)

const (
	accessKey = "accessKey"
	clientID  = "clientID"

	apiUrl = "https://a.wunderlist.com/api/v1"
)

var _ = Describe("Client", func() {
	var fakeHTTPHelper fakes.FakeHTTPHelper
	var fakeLogger fakes.FakeLogger

	var client wundergo.Client

	BeforeEach(func() {
		fakeHTTPHelper = fakes.FakeHTTPHelper{}
		fakeLogger = fakes.FakeLogger{}

		wundergo.NewHTTPHelper = func(accessToken string, clientID string) wundergo.HTTPHelper {
			return &fakeHTTPHelper
		}

		wundergo.NewLogger = func() wundergo.Logger {
			return &fakeLogger
		}

		client = wundergo.NewOauthClient(accessKey, clientID)
	})

	Describe("User operations", func() {
		expectedUrl := fmt.Sprintf("%s/user", apiUrl)

		Describe("getting user", func() {
			It("performs GET requests to /user", func() {
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
				badResponse := []byte("invalid json response")
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(badResponse, nil)
				})

				It("returns an empty user", func() {
					user, _ := client.User()

					Expect(user).To(Equal(wundergo.User{}))
				})

				It("forwards the error", func() {
					_, err := client.User()

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when valid response is received", func() {
				validResponse := []byte(`{"name": "newName"}`)
				expectedUser := wundergo.User{
					Name: "newName",
				}
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(validResponse, nil)
				})

				It("returns the unmarshalled user", func() {
					user, _ := client.User()

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
				expectedBody := fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name)

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
				badResponse := []byte("invalid json response")
				BeforeEach(func() {
					fakeHTTPHelper.PutReturns(badResponse, nil)
				})

				It("returns an empty user", func() {
					user, _ := client.UpdateUser(user)

					Expect(user).To(Equal(wundergo.User{}))
				})

				It("forwards the error", func() {
					_, err := client.UpdateUser(user)

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when valid response is received", func() {
				validResponse := []byte(`{"name": "newName"}`)
				expectedUser := wundergo.User{
					Name: "newName",
				}
				BeforeEach(func() {
					fakeHTTPHelper.PutReturns(validResponse, nil)
				})

				It("returns the unmarshalled user", func() {
					user, _ := client.UpdateUser(user)

					Expect(user).To(Equal(expectedUser))
				})
			})
		})
	})

	Describe("Users operations", func() {

		Describe("getting users", func() {
			expectedUrl := fmt.Sprintf("%s/users", apiUrl)

			Context("when ListID is not provided", func() {
				It("performs GET requests to /users", func() {
					client.Users()

					Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
					Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
				})

				Context("when httpHelper.Get returns an error", func() {
					expectedError := errors.New("httpHelper GET error")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(nil, expectedError)
					})

					It("returns an empty list of users", func() {
						user, _ := client.Users()

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("forwards the error", func() {
						_, err := client.Users()

						Expect(err).To(Equal(expectedError))
					})
				})

				Context("when unmarshalling json response returns an error", func() {
					badResponse := []byte("invalid json response")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(badResponse, nil)
					})

					It("returns an empty list of users", func() {
						user, _ := client.Users()

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("forwards the error", func() {
						_, err := client.Users()

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when valid response is received", func() {
					validResponse := []byte(`[{"name": "newName"}]`)
					expectedUsers := []wundergo.User{
						wundergo.User{
							Name: "newName",
						},
					}
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(validResponse, nil)
					})

					It("returns the unmarshalled list of users", func() {
						users, _ := client.Users()

						Expect(users).To(Equal(expectedUsers))
					})
				})
			})

			Context("when ListID == 0", func() {
				ListID := uint(0)
				expectedUrl := fmt.Sprintf("%s/users", apiUrl)

				It("performs GET requests to /users", func() {
					client.UsersForListID(ListID)

					Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
					Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
				})

				Context("when httpHelper.Get returns an error", func() {
					expectedError := errors.New("httpHelper GET error")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(nil, expectedError)
					})

					It("returns an empty list of users", func() {
						user, _ := client.UsersForListID(ListID)

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("forwards the error", func() {
						_, err := client.UsersForListID(ListID)

						Expect(err).To(Equal(expectedError))
					})
				})

				Context("when unmarshalling json response returns an error", func() {
					badResponse := []byte("invalid json response")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(badResponse, nil)
					})

					It("returns an empty list of users", func() {
						user, _ := client.UsersForListID(ListID)

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("forwards the error", func() {
						_, err := client.UsersForListID(ListID)

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when valid response is received", func() {
					validResponse := []byte(`[{"name": "newName"}]`)
					expectedUsers := []wundergo.User{
						wundergo.User{
							Name: "newName",
						},
					}
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(validResponse, nil)
					})

					It("returns the unmarshalled list of users", func() {
						users, _ := client.UsersForListID(ListID)

						Expect(users).To(Equal(expectedUsers))
					})
				})
			})

			Context("when ListID > 0", func() {
				ListID := uint(12345)
				expectedUrl := fmt.Sprintf("%s/users?list_id=%d", apiUrl, ListID)

				It("performs GET requests to /users with list_id param", func() {
					client.UsersForListID(ListID)

					Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
					Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
				})

				Context("when httpHelper.Get returns an error", func() {
					expectedError := errors.New("httpHelper GET error")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(nil, expectedError)
					})

					It("returns an empty list of users", func() {
						user, _ := client.UsersForListID(ListID)

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("forwards the error", func() {
						_, err := client.UsersForListID(ListID)

						Expect(err).To(Equal(expectedError))
					})
				})

				Context("when unmarshalling json response returns an error", func() {
					badResponse := []byte("invalid json response")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(badResponse, nil)
					})

					It("returns an empty list of users", func() {
						user, _ := client.UsersForListID(ListID)

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("forwards the error", func() {
						_, err := client.UsersForListID(ListID)

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when valid response is received", func() {
					validResponse := []byte(`[{"name": "newName"}]`)
					expectedUsers := []wundergo.User{
						wundergo.User{
							Name: "newName",
						},
					}
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(validResponse, nil)
					})

					It("returns the unmarshalled list of users", func() {
						users, _ := client.UsersForListID(ListID)

						Expect(users).To(Equal(expectedUsers))
					})
				})
			})

		})
	})
})
