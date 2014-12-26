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

var _ = Describe("Client - Task operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("Getting tasks for list", func() {
		listID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.TasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /tasks?list_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/tasks?list_id=%d", apiUrl, listID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Task{}, nil)
			client.TasksForListID(listID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TasksForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.TasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.TasksForListID(listID)

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
				_, err := client.TasksForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.TasksForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTasks := &[]wundergo.Task{
				wundergo.Task{
					Title: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTasks, nil)
			})

			It("returns the unmarshalled tasks without error", func() {
				task, err := client.TasksForListID(listID)

				Expect(err).To(BeNil())
				Expect(task).To(Equal(expectedTasks))
			})
		})
	})

	Describe("getting task by ID", func() {
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /tasks/:id", func() {
			expectedUrl := fmt.Sprintf("%s/tasks/%d", apiUrl, taskID)
			fakeJSONHelper.UnmarshalReturns(&wundergo.Task{}, nil)
			client.Task(taskID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Task(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.Task(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.Task(taskID)

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
				_, err := client.Task(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Task(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTask := &wundergo.Task{
				Title: "Test Title",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTask, nil)
			})

			It("returns the unmarshalled task without error", func() {
				task, err := client.Task(taskID)

				Expect(err).To(BeNil())
				Expect(task).To(Equal(expectedTask))
			})
		})
	})

	Describe("creating a new task", func() {
		taskTitle := "newTaskTitle"
		listID := uint(1)

		var assigneeID uint
		var completed bool
		var recurrenceType string
		var recurrenceCount uint
		var dueDate string
		var starred bool

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusCreated
			fakeHTTPHelper.PostReturns(dummyResponse, nil)

			assigneeID = uint(2)
			completed = false
			recurrenceType = "day"
			recurrenceCount = uint(3)
			dueDate = "1970-01-01"
			starred = false
		})

		It("performs POST requests to /tasks with JSONHelper-serialized body", func() {
			expectedUrl := fmt.Sprintf("%s/tasks", apiUrl)
			expectedBody := []byte("some request body")
			fakeJSONHelper.MarshalReturns(expectedBody, nil)

			fakeJSONHelper.UnmarshalReturns(&wundergo.Task{}, nil)
			client.CreateTask(
				taskTitle,
				listID,
				assigneeID,
				completed,
				recurrenceType,
				recurrenceCount,
				dueDate,
				starred,
			)

			Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PostArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when recurrenceType is not provided", func() {
			BeforeEach(func() {
				recurrenceType = ""
			})

			It("does not allow recurrenceCount to be non-zero", func() {
				_, err := client.CreateTask(
					taskTitle,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("recurrenceCount"))
			})
		})

		Context("when recurrenceCount is zero", func() {
			BeforeEach(func() {
				recurrenceCount = 0
			})

			It("does not allow recurrenceType to be provided", func() {
				_, err := client.CreateTask(
					taskTitle,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("recurrenceType"))
			})
		})

		Context("when httpHelper.Post returns an error", func() {
			expectedError := errors.New("httpHelper POST error")

			BeforeEach(func() {
				fakeHTTPHelper.PostReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateTask(
					taskTitle,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CreateTask(
					taskTitle,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PostReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.CreateTask(
					taskTitle,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

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
				_, err := client.CreateTask(
					taskTitle,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateTask(
					taskTitle,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedTask := &wundergo.Task{
				Title: "Test Title",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedTask, nil)
			})

			It("returns the unmarshalled task without error", func() {
				task, err := client.CreateTask(
					taskTitle,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

				Expect(err).To(BeNil())
				Expect(task).To(Equal(expectedTask))
			})
		})
	})
})
