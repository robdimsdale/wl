package wundergo_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Client - List operations", func() {

	BeforeEach(func() {
		initializeFakes()
		initializeClient()
	})

	Describe("getting lists", func() {
		It("performs GET requests to /lists", func() {
			expectedUrl := fmt.Sprintf("%s/lists", apiUrl)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.List{}, nil)
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
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns([]byte("invalid json response"), nil)
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
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

		Context("when valid response is received", func() {
			expectedLists := []wundergo.List{}
			BeforeEach(func() {
				fakeHTTPHelper.GetReturns([]byte(""), nil)
				fakeJSONHelper.UnmarshalReturns(&expectedLists, nil)
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
			fakeJSONHelper.UnmarshalReturns(&wundergo.List{}, nil)
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
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns([]byte("invalid json response"), nil)
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
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

		Context("when valid response is received", func() {
			expectedList := wundergo.List{}

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns([]byte(""), nil)
				fakeJSONHelper.UnmarshalReturns(&expectedList, nil)
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
			fakeJSONHelper.UnmarshalReturns(&wundergo.ListTaskCount{}, nil)
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
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns([]byte("invalid json response"), nil)
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
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

		Context("when valid response is received", func() {
			expectedListTaskCount := wundergo.ListTaskCount{}

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns([]byte(""), nil)
				fakeJSONHelper.UnmarshalReturns(&expectedListTaskCount, nil)
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
				fakeJSONHelper.UnmarshalReturns(&wundergo.List{}, nil)
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
				expectedError := errors.New("jsonHelper error")

				BeforeEach(func() {
					fakeHTTPHelper.PostReturns([]byte("invalid json response"), nil)
					fakeJSONHelper.UnmarshalReturns(nil, expectedError)
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

			Context("when valid response is received", func() {
				expectedList := wundergo.List{}

				BeforeEach(func() {
					fakeHTTPHelper.PostReturns([]byte(""), nil)
					fakeJSONHelper.UnmarshalReturns(&expectedList, nil)
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
			ID: uint(1),
		}
		expectedBody := []byte{}

		BeforeEach(func() {
			fakeJSONHelper.MarshalReturns(expectedBody, nil)
		})

		It("performs PATCH requests to /lists/:id", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.List{}, nil)
			expectedUrl := fmt.Sprintf("%s/lists/%d", apiUrl, list.ID)

			client.UpdateList(list)

			Expect(fakeHTTPHelper.PatchCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PatchArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when marshalling list returns an error", func() {
			expectedError := errors.New("JSONHelper marshal error")

			BeforeEach(func() {
				fakeJSONHelper.MarshalReturns(nil, expectedError)
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

		Context("when httpHelper.Patch returns an error", func() {
			expectedError := errors.New("httpHelper PATCH error")

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
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeHTTPHelper.PatchReturns([]byte("invalid json response"), nil)
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
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

		Context("when valid response is received", func() {
			expectedList := wundergo.List{}

			BeforeEach(func() {
				fakeHTTPHelper.PatchReturns([]byte(""), nil)
				fakeJSONHelper.UnmarshalReturns(&expectedList, nil)
			})

			It("returns the unmarshalled list without error", func() {
				list, err := client.UpdateList(list)

				Expect(err).To(BeNil())
				Expect(list).To(Equal(expectedList))
			})
		})
	})
})
