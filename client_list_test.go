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

var _ = Describe("Client - List operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("getting lists", func() {

		BeforeEach(func() {
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

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

			It("forwards the error", func() {
				_, err := client.Lists()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.Lists()

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
				_, err := client.Lists()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Lists()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedLists := &[]wundergo.List{
				wundergo.List{
					Title: "Test Title",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedLists, nil)
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

		BeforeEach(func() {
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

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

			It("forwards the error", func() {
				_, err := client.List(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.List(listID)

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
				_, err := client.List(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.List(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedList := &wundergo.List{
				Title: "Test Title",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedList, nil)
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

		BeforeEach(func() {
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

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

			It("forwards the error", func() {
				_, err := client.ListTaskCount(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.ListTaskCount(listID)

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
				_, err := client.Lists()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.ListTaskCount(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedListTaskCount := &wundergo.ListTaskCount{
				CompletedCount: 1,
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedListTaskCount, nil)
			})

			It("returns the unmarshalled list task count without error", func() {
				listTaskCount, err := client.ListTaskCount(listID)

				Expect(err).To(BeNil())
				Expect(listTaskCount).To(Equal(expectedListTaskCount))
			})
		})

	})

	Describe("creating a new list", func() {
		expectedUrl := fmt.Sprintf("%s/lists", apiUrl)
		listTitle := "newListTitle"
		expectedBody := []byte(fmt.Sprintf(`{"title":"%s"}`, listTitle))

		BeforeEach(func() {
			fakeHTTPHelper.PostReturns(dummyResponse, nil)
		})

		It("performs POST requests to /lists with new list title in body", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.List{}, nil)
			client.CreateList(listTitle)

			Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PostArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when httpHelper.Post returns an error", func() {
			expectedError := errors.New("httpHelper POST error")

			BeforeEach(func() {
				fakeHTTPHelper.PostReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateList(listTitle)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PostReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.CreateList(listTitle)

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
				_, err := client.CreateList(listTitle)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateList(listTitle)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedList := &wundergo.List{
				Title: "Test Title",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedList, nil)
			})

			It("returns the unmarshalled list task count without error", func() {
				list, err := client.CreateList(listTitle)

				Expect(err).To(BeNil())
				Expect(list).To(Equal(expectedList))
			})
		})
	})

	Describe("updating a list", func() {
		list := wundergo.List{
			ID: uint(1),
		}
		expectedBody := []byte{}

		BeforeEach(func() {
			fakeHTTPHelper.PatchReturns(dummyResponse, nil)
		})

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

			It("forwards the error", func() {
				_, err := client.UpdateList(list)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PatchReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.UpdateList(list)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when reading body returns an error", func() {
			expectedError := errors.New("read error")
			BeforeEach(func() {
				dummyResponse.Body = erroringReadCloser{
					readError: expectedError,
				}
				fakeHTTPHelper.PatchReturns(dummyResponse, nil)
			})

			It("forwards the error", func() {
				_, err := client.UpdateList(list)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateList(list)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedList := &wundergo.List{
				Title: "Test Title",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedList, nil)
			})

			It("returns the unmarshalled list without error", func() {
				list, err := client.UpdateList(list)

				Expect(err).To(BeNil())
				Expect(list).To(Equal(expectedList))
			})
		})
	})

	Describe("deleting a list", func() {
		list := wundergo.List{
			ID:       uint(1),
			Revision: 3,
		}

		BeforeEach(func() {
			fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
		})

		It("performs DELETE requests to /lists/:id?revision=:revision", func() {
			expectedUrl := fmt.Sprintf("%s/lists/%d?revision=%d", apiUrl, list.ID, list.Revision)

			client.DeleteList(list)

			Expect(fakeHTTPHelper.DeleteCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.DeleteArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Delete returns an error", func() {
			expectedError := errors.New("httpHelper DELETE error")

			BeforeEach(func() {
				fakeHTTPHelper.DeleteReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				err := client.DeleteList(list)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				err := client.DeleteList(list)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when reading body returns an error", func() {
			expectedError := errors.New("read error")
			BeforeEach(func() {
				dummyResponse.Body = erroringReadCloser{
					readError: expectedError,
				}
				fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
			})

			It("forwards the error", func() {
				err := client.DeleteList(list)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			It("deletes the list without error", func() {
				err := client.DeleteList(list)

				Expect(err).To(BeNil())
			})
		})
	})
})
