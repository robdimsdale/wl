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

var _ = Describe("Client - Note operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("Getting notes for list", func() {
		listID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.NotesForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /notes?list_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/notes?list_id=%d", apiUrl, listID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Note{}, nil)
			client.NotesForListID(listID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.NotesForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.NotesForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.NotesForListID(listID)

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
				_, err := client.NotesForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.NotesForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedNotes := &[]wundergo.Note{
				wundergo.Note{
					Content: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedNotes, nil)
			})

			It("returns the unmarshalled notes without error", func() {
				note, err := client.NotesForListID(listID)

				Expect(err).To(BeNil())
				Expect(note).To(Equal(expectedNotes))
			})
		})
	})

	Describe("Getting notes for task", func() {
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when TaskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.NotesForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /notes?task_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/notes?task_id=%d", apiUrl, taskID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Note{}, nil)
			client.NotesForTaskID(taskID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.NotesForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.NotesForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.NotesForTaskID(taskID)

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
				_, err := client.NotesForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.NotesForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedNotes := &[]wundergo.Note{
				wundergo.Note{
					Content: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedNotes, nil)
			})

			It("returns the unmarshalled notes without error", func() {
				note, err := client.NotesForTaskID(taskID)

				Expect(err).To(BeNil())
				Expect(note).To(Equal(expectedNotes))
			})
		})
	})

	Describe("getting note by ID", func() {
		noteID := uint(1)
		expectedUrl := fmt.Sprintf("%s/notes/%d", apiUrl, noteID)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /notes/:id", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.Note{}, nil)
			client.Note(noteID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Note(noteID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns error", func() {
				_, err := client.Note(noteID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.Note(noteID)

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
				_, err := client.Note(noteID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Note(noteID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedNote := &wundergo.Note{
				Content: "Test Content",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedNote, nil)
			})

			It("returns the unmarshalled note without error", func() {
				note, err := client.Note(noteID)

				Expect(err).To(BeNil())
				Expect(note).To(Equal(expectedNote))
			})
		})
	})
})
