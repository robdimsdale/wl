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

var _ = Describe("Client - Task position operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("getting task positions for list", func() {
		listID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.TaskPositionsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /task_positions?list_id=%d", func() {
			expectedUrl := fmt.Sprintf("%s/task_positions?list_id=%d", apiURL, listID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Position{}, nil)
			client.TaskPositionsForListID(listID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")
			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskPositionsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.TaskPositionsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.TaskPositionsForListID(listID)

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
				_, err := client.TaskPositionsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskPositionsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTaskPositions := &[]wundergo.Position{
				wundergo.Position{
					ID: 1234,
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTaskPositions, nil)
			})

			It("returns the unmarshalled array of task positions without error", func() {
				taskPositions, err := client.TaskPositionsForListID(listID)

				Expect(err).To(BeNil())
				Expect(taskPositions).To(Equal(expectedTaskPositions))
			})
		})
	})

	Describe("getting taskPosition by ID", func() {
		taskPositionID := uint(1)
		expectedUrl := fmt.Sprintf("%s/task_positions/%d", apiURL, taskPositionID)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /task_positions/:id", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.Position{}, nil)
			client.TaskPosition(taskPositionID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskPosition(taskPositionID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.TaskPosition(taskPositionID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.TaskPosition(taskPositionID)

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
				_, err := client.TaskPosition(taskPositionID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskPosition(taskPositionID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTaskPosition := &wundergo.Position{
				ID: 1234,
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTaskPosition, nil)
			})

			It("returns the unmarshalled task position without error", func() {
				taskPosition, err := client.TaskPosition(taskPositionID)

				Expect(err).To(BeNil())
				Expect(taskPosition).To(Equal(expectedTaskPosition))
			})
		})
	})

	Describe("updating a task position", func() {
		taskPosition := wundergo.Position{
			ID: 1,
		}

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.PatchReturns(dummyResponse, nil)
		})

		It("performs PATCH requests to /task_positions/:id", func() {
			expectedBody := []byte{}
			fakeJSONHelper.MarshalReturns(expectedBody, nil)
			fakeJSONHelper.UnmarshalReturns(&wundergo.Position{}, nil)
			expectedUrl := fmt.Sprintf("%s/task_positions/%d", apiURL, taskPosition.ID)

			client.UpdateTaskPosition(taskPosition)

			Expect(fakeHTTPHelper.PatchCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PatchArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when marshalling task returns an error", func() {
			expectedError := errors.New("JSONHelper marshal error")

			BeforeEach(func() {
				fakeJSONHelper.MarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateTaskPosition(taskPosition)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when httpHelper.Patch returns an error", func() {
			expectedError := errors.New("httpHelper PATCH error")

			BeforeEach(func() {
				fakeHTTPHelper.PatchReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateTaskPosition(taskPosition)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.UpdateTaskPosition(taskPosition)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PatchReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.UpdateTaskPosition(taskPosition)

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
				_, err := client.UpdateTaskPosition(taskPosition)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateTaskPosition(taskPosition)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTaskPosition := &wundergo.Position{
				Values: []uint{1, 2},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTaskPosition, nil)
			})

			It("returns the unmarshalled task without error", func() {
				taskPosition, err := client.UpdateTaskPosition(*expectedTaskPosition)

				Expect(err).To(BeNil())
				Expect(taskPosition).To(Equal(expectedTaskPosition))
			})
		})
	})

})
