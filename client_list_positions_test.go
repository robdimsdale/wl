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

var _ = Describe("Client - List position operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("getting list positions", func() {

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /list_positions", func() {
			expectedUrl := fmt.Sprintf("%s/list_positions", apiUrl)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Position{}, nil)
			client.ListPositions()

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")
			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.ListPositions()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.ListPositions()

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.ListPositions()

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
				_, err := client.ListPositions()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.ListPositions()

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedLists := &[]wundergo.Position{
				wundergo.Position{
					ID: 1234,
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedLists, nil)
			})

			It("returns the unmarshalled array of lists without error", func() {
				lists, err := client.ListPositions()

				Expect(err).To(BeNil())
				Expect(lists).To(Equal(expectedLists))
			})
		})
	})
})
