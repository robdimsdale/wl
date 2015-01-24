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
			expectedUrl := fmt.Sprintf("%s/tasks?list_id=%d", apiURL, listID)

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

	Describe("Getting completed tasks for list", func() {
		listID := uint(1)
		completed := true

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /tasks?list_id=:id&completed=:completed", func() {
			expectedUrl := fmt.Sprintf("%s/tasks?list_id=%d&completed=%t", apiURL, listID, completed)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Task{}, nil)
			client.CompletedTasksForListID(listID, completed)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.CompletedTasksForListID(listID, completed)

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
				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CompletedTasksForListID(listID, completed)

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
				tasks, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(BeNil())
				Expect(tasks).To(Equal(expectedTasks))
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
			expectedUrl := fmt.Sprintf("%s/tasks/%d", apiURL, taskID)
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
			expectedUrl := fmt.Sprintf("%s/tasks", apiURL)
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

		Context("when ListID == 0", func() {
			listID := uint(0)

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

		Context("when marshalling task-creation config struct returns an error", func() {
			expectedError := errors.New("JSONHelper marshal error")

			BeforeEach(func() {
				fakeJSONHelper.MarshalReturns(nil, expectedError)
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

	Describe("updating a task", func() {
		var expectedGetTask *wundergo.Task
		var task wundergo.Task

		var dummyGetResponse *http.Response

		BeforeEach(func() {
			task = wundergo.Task{
				ID:       uint(1),
				Revision: 2,
			}

			expectedGetTask = &wundergo.Task{
				Title: "testy",
			}

			dummyGetResponse = &http.Response{}
			dummyGetResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
			dummyGetResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyGetResponse, nil)

			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.PatchReturns(dummyResponse, nil)

			fakeJSONHelper.UnmarshalReturns(expectedGetTask, nil)
		})

		Context("when recurrenceType is not provided", func() {
			BeforeEach(func() {
				task.RecurrenceType = ""
			})

			It("does not allow recurrenceCount to be non-zero", func() {
				task.RecurrenceCount = 1
				_, err := client.UpdateTask(task)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("recurrenceCount"))
			})
		})

		Context("when recurrenceCount is zero", func() {
			BeforeEach(func() {
				task.RecurrenceCount = 0
			})

			It("does not allow recurrenceType to be provided", func() {
				task.RecurrenceType = "day"
				_, err := client.UpdateTask(task)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("recurrenceType"))
			})
		})

		Describe("during initial GET request", func() {
			It("performs GET request to /tasks/:id", func() {
				expectedUrl := fmt.Sprintf("%s/tasks/%d", apiURL, task.ID)
				fakeJSONHelper.UnmarshalReturns(&wundergo.Task{}, nil)
				client.UpdateTask(task)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
			})

			Context("when initial GET request returns with error", func() {
				expectedError := errors.New("httpHelper GET error")

				BeforeEach(func() {
					fakeHTTPHelper.GetReturns(nil, expectedError)
				})

				It("forwards the error", func() {
					_, err := client.UpdateTask(task)

					Expect(err).To(Equal(expectedError))
				})
			})
		})

		Describe("assigneeID updates", func() {
			Context("when original task had an assignee", func() {
				BeforeEach(func() {
					expectedGetTask.AssigneeID = 1
				})

				Context("and new task has no assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = 0
					})

					It("adds 'assignee_id' to fields to remove", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.Remove).To(Equal([]string{"assignee_id"}))
					})

					It("sets assigneeID to 0", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.AssigneeID).To(Equal(task.AssigneeID))
					})
				})

				Context("and new task has same assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = expectedGetTask.AssigneeID
					})

					It("calls JSONHelper.marshal with same value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.AssigneeID).To(Equal(task.AssigneeID))
					})
				})

				Context("and new task has different assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = 2
					})

					It("calls JSONHelper.marshal with new value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.AssigneeID).To(Equal(task.AssigneeID))
					})
				})
			})

			Context("when original task had no assignee", func() {
				BeforeEach(func() {
					expectedGetTask.AssigneeID = 0
				})

				Context("and new task has no assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = 0
					})

					It("sets assigneeID to 0", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.AssigneeID).To(Equal(task.AssigneeID))
					})
				})

				Context("and new task has an assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = 1
					})

					It("calls JSONHelper.marshal with new value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.AssigneeID).To(Equal(task.AssigneeID))
					})
				})
			})
		})

		Describe("dueDate updates", func() {
			Context("when original task had a due date", func() {
				BeforeEach(func() {
					expectedGetTask.DueDate = "1970-01-01"
				})

				Context("and new task has empty due date", func() {
					BeforeEach(func() {
						task.DueDate = ""
					})

					It("adds 'due_date' to fields to remove", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.Remove).To(Equal([]string{"due_date"}))
					})

					It("sets dueDate to empty string", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.DueDate).To(Equal(task.DueDate))
					})
				})

				Context("and new task has same due date", func() {
					BeforeEach(func() {
						task.DueDate = expectedGetTask.DueDate
					})

					It("calls JSONHelper.marshal with same value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.DueDate).To(Equal(task.DueDate))
					})
				})

				Context("and new task has different due date", func() {
					BeforeEach(func() {
						task.DueDate = "1971-01-01"
					})

					It("calls JSONHelper.marshal with new value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.DueDate).To(Equal(task.DueDate))
					})
				})
			})

			Context("when original task had no due date", func() {
				BeforeEach(func() {
					expectedGetTask.DueDate = ""
				})

				Context("and new task has empty due date", func() {
					BeforeEach(func() {
						task.DueDate = ""
					})

					It("sets dueDate to empty", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.DueDate).To(Equal(task.DueDate))
					})
				})

				Context("and new task has a due date", func() {
					BeforeEach(func() {
						task.DueDate = "1971-01-01"
					})

					It("calls JSONHelper.marshal with new value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.DueDate).To(Equal(task.DueDate))
					})
				})
			})
		})

		Describe("recurrence updates", func() {
			Context("when recurrence was previously set", func() {
				BeforeEach(func() {
					expectedGetTask.RecurrenceType = "day"
					expectedGetTask.RecurrenceCount = uint(1)
				})

				Context("and recurrence is now not set", func() {
					BeforeEach(func() {
						task.RecurrenceType = ""
						task.RecurrenceCount = 0
					})

					It("adds 'recurrence_type' and 'recurrence_count' to fields to remove", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.Remove).To(Equal([]string{"recurrence_type", "recurrence_count"}))
					})

					It("unsets recurrence", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.RecurrenceType).To(Equal(task.RecurrenceType))
						Expect(actualTuc.RecurrenceCount).To(Equal(task.RecurrenceCount))
					})
				})

				Context("and recurrence type is now set to a different value", func() {
					BeforeEach(func() {
						task.RecurrenceType = "week"
						task.RecurrenceCount = uint(1)
					})

					It("sets recurrence to new value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.RecurrenceType).To(Equal(task.RecurrenceType))
						Expect(actualTuc.RecurrenceCount).To(Equal(task.RecurrenceCount))
					})
				})

				Context("and recurrence count is now set to a different value", func() {
					BeforeEach(func() {
						task.RecurrenceType = "day"
						task.RecurrenceCount = uint(2)
					})

					It("sets recurrence to new value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.RecurrenceType).To(Equal(task.RecurrenceType))
						Expect(actualTuc.RecurrenceCount).To(Equal(task.RecurrenceCount))
					})
				})

				Context("and recurrence is now set to the same value", func() {
					BeforeEach(func() {
						task.RecurrenceType = "day"
						task.RecurrenceCount = uint(1)
					})

					It("sets recurrence to same value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.RecurrenceType).To(Equal(task.RecurrenceType))
						Expect(actualTuc.RecurrenceCount).To(Equal(task.RecurrenceCount))
					})
				})
			})

			Context("when recurrence was not previously set", func() {
				BeforeEach(func() {
					expectedGetTask.RecurrenceType = ""
					expectedGetTask.RecurrenceCount = 0
				})

				Context("and recurrence is still not set", func() {
					BeforeEach(func() {
						task.RecurrenceType = ""
						task.RecurrenceCount = 0
					})

					It("leaves recurrence unset", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.RecurrenceType).To(Equal(task.RecurrenceType))
						Expect(actualTuc.RecurrenceCount).To(Equal(task.RecurrenceCount))
					})
				})

				Context("and recurrence is now set", func() {
					BeforeEach(func() {
						task.RecurrenceType = "day"
						task.RecurrenceCount = uint(1)
					})

					It("sets recurrence to new value", func() {
						client.UpdateTask(task)

						Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
						arg0 := fakeJSONHelper.MarshalArgsForCall(0)
						actualTuc := arg0.(wundergo.TaskUpdateConfig)
						Expect(actualTuc.RecurrenceType).To(Equal(task.RecurrenceType))
						Expect(actualTuc.RecurrenceCount).To(Equal(task.RecurrenceCount))
					})
				})
			})
		})

		Describe("completed updates", func() {
			Context("when completed state changes", func() {
				BeforeEach(func() {
					expectedGetTask.Completed = true
					task.Completed = false
				})

				It("calls JSONHelper.marshal with new value", func() {
					client.UpdateTask(task)

					Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
					arg0 := fakeJSONHelper.MarshalArgsForCall(0)
					actualTuc := arg0.(wundergo.TaskUpdateConfig)
					Expect(actualTuc.Completed).To(Equal(task.Completed))
				})
			})

			Context("when completed state is unchanged", func() {
				BeforeEach(func() {
					expectedGetTask.Completed = true
					task.Completed = true
				})

				It("calls JSONHelper.marshal with new value", func() {
					client.UpdateTask(task)

					Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
					arg0 := fakeJSONHelper.MarshalArgsForCall(0)
					actualTuc := arg0.(wundergo.TaskUpdateConfig)
					Expect(actualTuc.Completed).To(Equal(task.Completed))
				})
			})
		})

		Describe("starred updates", func() {
			Context("when starred state changes", func() {
				BeforeEach(func() {
					expectedGetTask.Starred = true
					task.Starred = false
				})

				It("calls JSONHelper.marshal with new value", func() {
					client.UpdateTask(task)

					Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
					arg0 := fakeJSONHelper.MarshalArgsForCall(0)
					actualTuc := arg0.(wundergo.TaskUpdateConfig)
					Expect(actualTuc.Starred).To(Equal(task.Starred))
				})
			})

			Context("when starred state is unchanged", func() {
				BeforeEach(func() {
					expectedGetTask.Starred = true
					task.Starred = true
				})

				It("calls JSONHelper.marshal with new value", func() {
					client.UpdateTask(task)

					Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
					arg0 := fakeJSONHelper.MarshalArgsForCall(0)
					actualTuc := arg0.(wundergo.TaskUpdateConfig)
					Expect(actualTuc.Starred).To(Equal(task.Starred))
				})
			})
		})

		Describe("title updates", func() {
			Context("when title changes", func() {
				BeforeEach(func() {
					expectedGetTask.Title = "Old Title"
					task.Title = "Old Title"
				})

				It("calls JSONHelper.marshal with new value", func() {
					client.UpdateTask(task)

					Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
					arg0 := fakeJSONHelper.MarshalArgsForCall(0)
					actualTuc := arg0.(wundergo.TaskUpdateConfig)
					Expect(actualTuc.Title).To(Equal(task.Title))
				})
			})

			Context("when title is unchanged", func() {
				BeforeEach(func() {
					expectedGetTask.Title = "Old Title"
					task.Title = "New Title"
				})

				It("calls JSONHelper.marshal with new value", func() {
					client.UpdateTask(task)

					Expect(fakeJSONHelper.MarshalCallCount()).To(Equal(1))
					arg0 := fakeJSONHelper.MarshalArgsForCall(0)
					actualTuc := arg0.(wundergo.TaskUpdateConfig)
					Expect(actualTuc.Title).To(Equal(task.Title))
				})
			})
		})

		It("performs PATCH requests to /tasks/:id", func() {
			expectedBody := []byte{}
			fakeJSONHelper.MarshalReturns(expectedBody, nil)
			fakeJSONHelper.UnmarshalReturns(&wundergo.Task{}, nil)
			expectedUrl := fmt.Sprintf("%s/tasks/%d", apiURL, task.ID)

			client.UpdateTask(task)

			Expect(fakeHTTPHelper.PatchCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PatchArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when marshalling update body returns an error", func() {
			expectedError := errors.New("JSONHelper marshal error")

			BeforeEach(func() {
				fakeJSONHelper.MarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when httpHelper.Patch returns an error", func() {
			expectedError := errors.New("httpHelper PATCH error")

			BeforeEach(func() {
				fakeHTTPHelper.PatchReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PatchReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.UpdateTask(task)

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
				_, err := client.UpdateTask(task)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				callCount := 0
				fakeJSONHelper.UnmarshalStub = func([]byte, interface{}) (interface{}, error) {
					callCount++
					switch callCount {
					case 1:
						return expectedGetTask, nil
					default:
						return nil, expectedError
					}
				}
			})

			It("forwards the error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUpdateTask := &wundergo.Task{
				Title: "Updated Title",
			}

			BeforeEach(func() {
				callCount := 0
				fakeJSONHelper.UnmarshalStub = func([]byte, interface{}) (interface{}, error) {
					callCount++
					switch callCount {
					case 1:
						return expectedGetTask, nil
					default:
						return expectedUpdateTask, nil
					}
				}
			})

			It("returns the unmarshalled task without error", func() {
				task, err := client.UpdateTask(task)

				Expect(err).To(BeNil())
				Expect(task).To(Equal(expectedUpdateTask))
			})
		})
	})

	Describe("deleting a task", func() {
		task := wundergo.Task{
			ID:       uint(1),
			Revision: 3,
		}

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusNoContent
			fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
		})

		It("performs DELETE requests to /tasks/:id?revision=:revision", func() {
			expectedUrl := fmt.Sprintf("%s/tasks/%d?revision=%d", apiURL, task.ID, task.Revision)

			client.DeleteTask(task)

			Expect(fakeHTTPHelper.DeleteCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.DeleteArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Delete returns an error", func() {
			expectedError := errors.New("httpHelper DELETE error")

			BeforeEach(func() {
				fakeHTTPHelper.DeleteReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				err := client.DeleteTask(task)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns error", func() {
				err := client.DeleteTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				err := client.DeleteTask(task)

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
				err := client.DeleteTask(task)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			It("deletes the task without error", func() {
				err := client.DeleteTask(task)

				Expect(err).To(BeNil())
			})
		})
	})
})
