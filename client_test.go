package wundergo_test

import (
	"encoding/json"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/fakes"
)

const (
	dummyAccessToken = "dummyAccessToken"
	dummyClientID    = "dummyClientID"

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

		client = wundergo.NewOauthClient(dummyAccessToken, dummyClientID)
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

				It("returns an error", func() {
					_, err := client.User()

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when valid response is received", func() {
				var expectedUser wundergo.User

				BeforeEach(func() {
					validResponse := []byte(`{"name": "newName"}`)
					fakeHTTPHelper.GetReturns(validResponse, nil)

					err := json.Unmarshal(validResponse, &expectedUser)
					Expect(err).To(BeNil())
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
				badResponse := []byte("invalid json response")
				BeforeEach(func() {
					fakeHTTPHelper.PutReturns(badResponse, nil)
				})

				It("returns an empty user", func() {
					user, _ := client.UpdateUser(user)

					Expect(user).To(Equal(wundergo.User{}))
				})

				It("returns an error", func() {
					_, err := client.UpdateUser(user)

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when valid response is received", func() {
				var expectedUser wundergo.User

				BeforeEach(func() {
					validResponse := []byte(`{"name": "newName"}`)
					fakeHTTPHelper.PutReturns(validResponse, nil)

					err := json.Unmarshal(validResponse, &expectedUser)
					Expect(err).To(BeNil())
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
					badResponse := []byte("invalid json response")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(badResponse, nil)
					})

					It("returns an empty array of users", func() {
						user, _ := client.Users()

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("returns an error", func() {
						_, err := client.Users()

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when valid response is received", func() {
					var expectedUsers []wundergo.User
					BeforeEach(func() {
						validResponse := []byte(`[{"name": "newName"}]`)
						fakeHTTPHelper.GetReturns(validResponse, nil)

						err := json.Unmarshal(validResponse, &expectedUsers)
						Expect(err).To(BeNil())
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
					badResponse := []byte("invalid json response")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(badResponse, nil)
					})

					It("returns an empty array of users", func() {
						user, _ := client.UsersForListID(ListID)

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("returns an error", func() {
						_, err := client.UsersForListID(ListID)

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when valid response is received", func() {
					var expectedUsers []wundergo.User
					BeforeEach(func() {
						validResponse := []byte(`[{"name": "newName"}]`)
						fakeHTTPHelper.GetReturns(validResponse, nil)

						err := json.Unmarshal(validResponse, &expectedUsers)
						Expect(err).To(BeNil())
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
					badResponse := []byte("invalid json response")
					BeforeEach(func() {
						fakeHTTPHelper.GetReturns(badResponse, nil)
					})

					It("returns an empty array of users", func() {
						user, _ := client.UsersForListID(ListID)

						Expect(user).To(Equal([]wundergo.User{}))
					})

					It("returns an error", func() {
						_, err := client.UsersForListID(ListID)

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when valid response is received", func() {
					var expectedUsers []wundergo.User

					BeforeEach(func() {
						validResponse := []byte(`[{"name": "newName"}]`)
						fakeHTTPHelper.GetReturns(validResponse, nil)

						err := json.Unmarshal(validResponse, &expectedUsers)
						Expect(err).To(BeNil())
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

	Describe("List operations", func() {
		Describe("getting lists", func() {
			expectedUrl := fmt.Sprintf("%s/lists", apiUrl)

			It("performs GET requests to /lists", func() {
				client.Lists()

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})

			Context("when httpHelper.Get returns an error", func() {
				expectedError := errors.New("httpHelper GET error")
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(nil, expectedError)
				})

				It("returns an empty array of lists", func() {
					lists, _ := client.Lists()

					Expect(lists).To(Equal([]wundergo.List{}))
				})

				It("forwards the error", func() {
					_, err := client.Lists()

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				badResponse := []byte("invalid json response")
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(badResponse, nil)
				})

				It("returns an empty array of lists", func() {
					lists, _ := client.Lists()

					Expect(lists).To(Equal([]wundergo.List{}))
				})

				It("returns an error", func() {
					_, err := client.Lists()

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when valid response is received", func() {
				validResponse := []byte(`[{
          "id": 83526310,
          "created_at": "2013-08-30T08:29:46.203Z",
          "title": "Read Later",
          "list_type": "list",
          "type": "list",
          "revision": 10
          }]`)
				var expectedLists []wundergo.List
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(validResponse, nil)

					err := json.Unmarshal(validResponse, &expectedLists)
					Expect(err).To(BeNil())
				})

				It("returns the unmarshalled array of lists without error", func() {
					lists, err := client.Lists()

					Expect(err).To(BeNil())
					Expect(lists).To(Equal(expectedLists))
				})
			})
		})

		Describe("getting list by ID", func() {
			listID := uint(1)
			expectedUrl := fmt.Sprintf("%s/lists/%d", apiUrl, listID)

			It("performs GET requests to /lists/:id", func() {
				client.List(listID)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})

			Context("when httpHelper.Get returns an error", func() {
				expectedError := errors.New("httpHelper GET error")
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(nil, expectedError)
				})

				It("returns an empty list", func() {
					list, _ := client.List(listID)

					Expect(list).To(Equal(wundergo.List{}))
				})

				It("forwards the error", func() {
					_, err := client.List(listID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				badResponse := []byte("invalid json response")
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(badResponse, nil)
				})

				It("returns an empty list", func() {
					list, _ := client.List(listID)

					Expect(list).To(Equal(wundergo.List{}))
				})

				It("returns an error", func() {
					_, err := client.List(listID)

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when valid response is received", func() {
				validResponse := []byte(`{
                                "id": 83526310,
                                "created_at": "2013-08-30T08:29:46.203Z",
                                "title": "Read Later",
                                "list_type": "list",
                                "type": "list",
                                "revision": 10
                                }`)
				var expectedList wundergo.List
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(validResponse, nil)

					err := json.Unmarshal(validResponse, &expectedList)
					Expect(err).To(BeNil())
				})

				It("returns the unmarshalled list without error", func() {
					list, err := client.List(listID)

					Expect(err).To(BeNil())
					Expect(list).To(Equal(expectedList))
				})
			})
		})

		Describe("getting a list's task count", func() {
			listID := uint(1)
			expectedUrl := fmt.Sprintf("%s/lists/tasks_count?list_id=%d", apiUrl, listID)

			It("performs GET requests to /lists/tasks_count?list_id=:id", func() {
				client.ListTaskCount(listID)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})

			Context("when httpHelper.Get returns an error", func() {
				expectedError := errors.New("httpHelper GET error")
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(nil, expectedError)
				})

				It("returns an empty list task count", func() {
					listTaskCount, _ := client.ListTaskCount(listID)

					Expect(listTaskCount).To(Equal(wundergo.ListTaskCount{}))
				})

				It("forwards the error", func() {
					_, err := client.ListTaskCount(listID)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				badResponse := []byte("invalid json response")
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(badResponse, nil)
				})

				It("returns an empty list task count", func() {
					listTaskCount, _ := client.ListTaskCount(listID)

					Expect(listTaskCount).To(Equal(wundergo.ListTaskCount{}))
				})

				It("returns an error", func() {
					_, err := client.ListTaskCount(listID)

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when valid response is received", func() {
				validResponse := []byte(`{
          "id": 83526310,
          "completed_count": 100,
          "uncompleted_count": 200,
          "type": "tasks_count"
          }`)
				var expectedListTaskCount wundergo.ListTaskCount
				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(validResponse, nil)

					err := json.Unmarshal(validResponse, &expectedListTaskCount)
					Expect(err).To(BeNil())
				})

				It("returns the unmarshalled list task count without error", func() {
					listTaskCount, err := client.ListTaskCount(listID)

					Expect(err).To(BeNil())
					Expect(listTaskCount).To(Equal(expectedListTaskCount))
				})
			})

			Describe("creating a new list", func() {
				expectedUrl := fmt.Sprintf("%s/lists", apiUrl)
				listTitle := "newListTitle"
				expectedBody := []byte(fmt.Sprintf(`{"title":"%s"}`, listTitle))

				It("performs POST requests to /lists with new list title in body", func() {
					client.CreateList(listTitle)

					Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
					arg0, arg1 := fakeHTTPHelper.PostArgsForCall(0)
					Expect(arg0).To(Equal(expectedUrl))
					Expect(arg1).To(Equal(expectedBody))
				})

				Context("when httpHelper.Get returns an error", func() {
					expectedError := errors.New("httpHelper GET error")
					BeforeEach(func() {
						fakeHTTPHelper.PostReturns(nil, expectedError)
					})

					It("returns an empty list", func() {
						list, _ := client.CreateList(listTitle)

						Expect(list).To(Equal(wundergo.List{}))
					})

					It("forwards the error", func() {
						_, err := client.CreateList(listTitle)

						Expect(err).To(Equal(expectedError))
					})
				})

				Context("when unmarshalling json response returns an error", func() {
					badResponse := []byte("invalid json response")
					BeforeEach(func() {
						fakeHTTPHelper.PostReturns(badResponse, nil)
					})

					It("returns an empty list", func() {
						list, _ := client.CreateList(listTitle)

						Expect(list).To(Equal(wundergo.List{}))
					})

					It("returns an error", func() {
						_, err := client.CreateList(listTitle)

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when valid response is received", func() {
					validResponse := []byte(`{
            "id": 83526310,
            "created_at": "2013-08-30T08:29:46.203Z",
            "title": "Read Later",
            "revision": 1000,
            "type": "list"
            }`)
					var expectedList wundergo.List
					BeforeEach(func() {
						fakeHTTPHelper.PostReturns(validResponse, nil)

						err := json.Unmarshal(validResponse, &expectedList)
						Expect(err).To(BeNil())
					})

					It("returns the unmarshalled list task count without error", func() {
						list, err := client.CreateList(listTitle)

						Expect(err).To(BeNil())
						Expect(list).To(Equal(expectedList))
					})
				})
			})
		})

		Describe("updating a list", func() {
			list := wundergo.List{
				ID:       uint(1),
				Revision: 1,
				Title:    "newTitle",
			}
			expectedUrl := fmt.Sprintf("%s/lists/%d", apiUrl, list.ID)
			var expectedBody []byte

			BeforeEach(func() {
				var err error
				expectedBody, err = json.Marshal(list)

				Expect(err).To(BeNil())
			})

			It("performs PATCH requests to /lists/:id", func() {
				client.UpdateList(list)

				Expect(fakeHTTPHelper.PatchCallCount()).To(Equal(1))
				arg0, arg1 := fakeHTTPHelper.PatchArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))
				Expect(arg1).To(Equal(expectedBody))
			})

			Context("when httpHelper.Patch returns an error", func() {
				expectedError := errors.New("httpHelper GET error")
				BeforeEach(func() {
					fakeHTTPHelper.PatchReturns(nil, expectedError)
				})

				It("returns an empty list", func() {
					list, _ := client.UpdateList(list)

					Expect(list).To(Equal(wundergo.List{}))
				})

				It("forwards the error", func() {
					_, err := client.UpdateList(list)

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when unmarshalling json response returns an error", func() {
				badResponse := []byte("invalid json response")
				BeforeEach(func() {
					fakeHTTPHelper.PatchReturns(badResponse, nil)
				})

				It("returns an empty list", func() {
					list, _ := client.UpdateList(list)

					Expect(list).To(Equal(wundergo.List{}))
				})

				It("returns an error", func() {
					_, err := client.UpdateList(list)

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when valid response is received", func() {
				validResponse := []byte(`{
                                "id": 83526310,
                                "created_at": "2013-08-30T08:29:46.203Z",
                                "title": "Read Later",
                                "list_type": "list",
                                "type": "list",
                                "revision": 10
                                }`)
				var expectedList wundergo.List
				BeforeEach(func() {
					fakeHTTPHelper.PatchReturns(validResponse, nil)

					err := json.Unmarshal(validResponse, &expectedList)
					Expect(err).To(BeNil())
				})

				It("returns the unmarshalled list without error", func() {
					list, err := client.UpdateList(list)

					Expect(err).To(BeNil())
					Expect(list).To(Equal(expectedList))
				})
			})
		})
	})
})
