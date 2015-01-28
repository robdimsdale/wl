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

var _ = Describe("Client - TaskComment operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("Getting task comments for list", func() {
		listID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.TaskCommentsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /task_comments?list_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/task_comments?list_id=%d", apiURL, listID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.TaskComment{}, nil)
			client.TaskCommentsForListID(listID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskCommentsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.TaskCommentsForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.TaskCommentsForListID(listID)

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
				_, err := client.TaskCommentsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskCommentsForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTaskComments := &[]wundergo.TaskComment{
				wundergo.TaskComment{
					Text: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTaskComments, nil)
			})

			It("returns the unmarshalled task comment without error", func() {
				taskComments, err := client.TaskCommentsForListID(listID)

				Expect(err).To(BeNil())
				Expect(taskComments).To(Equal(expectedTaskComments))
			})
		})
	})

	Describe("Getting task comments for task", func() {
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when TaskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /task_comments?task_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/task_comments?task_id=%d", apiURL, taskID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.TaskComment{}, nil)
			client.TaskCommentsForTaskID(taskID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.TaskCommentsForTaskID(taskID)

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
				_, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTaskComments := &[]wundergo.TaskComment{
				wundergo.TaskComment{
					Text: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTaskComments, nil)
			})

			It("returns the unmarshalled task comments without error", func() {
				taskComments, err := client.TaskCommentsForTaskID(taskID)

				Expect(err).To(BeNil())
				Expect(taskComments).To(Equal(expectedTaskComments))
			})
		})
	})

	Describe("creating a new task comment", func() {
		text := "someText"
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusCreated
			fakeHTTPHelper.PostReturns(dummyResponse, nil)
		})

		It("performs POST requests to /task_comments with new task_comment content in body", func() {
			expectedUrl := fmt.Sprintf("%s/task_comments", apiURL)
			expectedBody := []byte(fmt.Sprintf(`{"text":"%s","task_id":%d}`, text, taskID))

			fakeJSONHelper.UnmarshalReturns(&wundergo.TaskComment{}, nil)
			client.CreateTaskComment(text, taskID)

			Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PostArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when taskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when httpHelper.Post returns an error", func() {
			expectedError := errors.New("httpHelper POST error")

			BeforeEach(func() {
				fakeHTTPHelper.PostReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PostReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.CreateTaskComment(text, taskID)

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
				_, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTaskComment := &wundergo.TaskComment{
				Text: "Test Text",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTaskComment, nil)
			})

			It("returns the unmarshalled task comment without error", func() {
				expectedTaskComment, err := client.CreateTaskComment(text, taskID)

				Expect(err).To(BeNil())
				Expect(expectedTaskComment).To(Equal(expectedTaskComment))
			})
		})
	})

	Describe("getting task comment by ID", func() {
		taskCommentID := uint(1)
		expectedUrl := fmt.Sprintf("%s/task_comments/%d", apiURL, taskCommentID)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /task_comments/:id", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.TaskComment{}, nil)
			client.TaskComment(taskCommentID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskComment(taskCommentID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.TaskComment(taskCommentID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.TaskComment(taskCommentID)

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
				_, err := client.TaskComment(taskCommentID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TaskComment(taskCommentID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTaskComment := &wundergo.TaskComment{
				Text: "Test Content",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTaskComment, nil)
			})

			It("returns the unmarshalled task comment without error", func() {
				taskComment, err := client.TaskComment(taskCommentID)

				Expect(err).To(BeNil())
				Expect(taskComment).To(Equal(expectedTaskComment))
			})
		})
	})

	Describe("deleting a task comment", func() {
		taskComment := wundergo.TaskComment{
			ID:       1,
			Revision: 3,
		}

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusNoContent
			fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
		})

		It("performs DELETE requests to /task_comments/:id?revision=:revision", func() {
			expectedUrl := fmt.Sprintf("%s/task_comments/%d?revision=%d", apiURL, taskComment.ID, taskComment.Revision)

			client.DeleteTaskComment(taskComment)

			Expect(fakeHTTPHelper.DeleteCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.DeleteArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Delete returns an error", func() {
			expectedError := errors.New("httpHelper DELETE error")

			BeforeEach(func() {
				fakeHTTPHelper.DeleteReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				err := client.DeleteTaskComment(taskComment)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				err := client.DeleteTaskComment(taskComment)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				err := client.DeleteTaskComment(taskComment)

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
				err := client.DeleteTaskComment(taskComment)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			It("deletes the note without error", func() {
				err := client.DeleteTaskComment(taskComment)

				Expect(err).To(BeNil())
			})
		})
	})
})
