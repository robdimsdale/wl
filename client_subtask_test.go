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

var _ = Describe("Client - Subtask operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("Getting subtasks for list", func() {
		listID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.SubtasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /subtasks?list_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/subtasks?list_id=%d", apiUrl, listID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Subtask{}, nil)
			client.SubtasksForListID(listID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.SubtasksForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.SubtasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.SubtasksForListID(listID)

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
				_, err := client.SubtasksForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.SubtasksForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedSubtasks := &[]wundergo.Subtask{
				wundergo.Subtask{
					Title: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedSubtasks, nil)
			})

			It("returns the unmarshalled subtasks without error", func() {
				subtask, err := client.SubtasksForListID(listID)

				Expect(err).To(BeNil())
				Expect(subtask).To(Equal(expectedSubtasks))
			})
		})
	})

	Describe("Getting subtasks for task", func() {
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when TaskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.SubtasksForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /subtasks?task_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/subtasks?task_id=%d", apiUrl, taskID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Subtask{}, nil)
			client.SubtasksForTaskID(taskID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.SubtasksForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.SubtasksForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.SubtasksForTaskID(taskID)

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
				_, err := client.SubtasksForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.SubtasksForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedSubtasks := &[]wundergo.Subtask{
				wundergo.Subtask{
					Title: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedSubtasks, nil)
			})

			It("returns the unmarshalled subtasks without error", func() {
				subtask, err := client.SubtasksForTaskID(taskID)

				Expect(err).To(BeNil())
				Expect(subtask).To(Equal(expectedSubtasks))
			})
		})
	})

	Describe("Getting completed subtasks for list", func() {
		listID := uint(1)
		completed := true

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when listID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.CompletedSubtasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /subtasks?list_id=:id&completed=:completed", func() {
			expectedUrl := fmt.Sprintf("%s/subtasks?list_id=%d&completed=%t", apiUrl, listID, completed)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Subtask{}, nil)
			client.CompletedSubtasksForListID(listID, completed)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CompletedSubtasksForListID(listID, completed)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CompletedSubtasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.CompletedSubtasksForListID(listID, completed)

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
				_, err := client.CompletedSubtasksForListID(listID, completed)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CompletedSubtasksForListID(listID, completed)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedSubtasks := &[]wundergo.Subtask{
				wundergo.Subtask{
					Title: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedSubtasks, nil)
			})

			It("returns the unmarshalled subtasks without error", func() {
				subtasks, err := client.CompletedSubtasksForListID(listID, completed)

				Expect(err).To(BeNil())
				Expect(subtasks).To(Equal(expectedSubtasks))
			})
		})
	})

	Describe("Getting completed subtasks for task", func() {
		taskID := uint(1)
		completed := true

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when taskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.CompletedSubtasksForTaskID(taskID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /subtasks?task_id=:id&completed=:completed", func() {
			expectedUrl := fmt.Sprintf("%s/subtasks?task_id=%d&completed=%t", apiUrl, taskID, completed)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Subtask{}, nil)
			client.CompletedSubtasksForTaskID(taskID, completed)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CompletedSubtasksForTaskID(taskID, completed)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CompletedSubtasksForTaskID(taskID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.CompletedSubtasksForTaskID(taskID, completed)

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
				_, err := client.CompletedSubtasksForTaskID(taskID, completed)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CompletedSubtasksForTaskID(taskID, completed)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedSubtasks := &[]wundergo.Subtask{
				wundergo.Subtask{
					Title: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedSubtasks, nil)
			})

			It("returns the unmarshalled subtasks without error", func() {
				subtasks, err := client.CompletedSubtasksForTaskID(taskID, completed)

				Expect(err).To(BeNil())
				Expect(subtasks).To(Equal(expectedSubtasks))
			})
		})
	})

	Describe("Getting subtask by ID", func() {
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /subtasks/:id", func() {
			expectedUrl := fmt.Sprintf("%s/subtasks/%d", apiUrl, taskID)
			fakeJSONHelper.UnmarshalReturns(&wundergo.Subtask{}, nil)
			client.Subtask(taskID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Subtask(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.Subtask(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.Subtask(taskID)

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
				_, err := client.Subtask(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Subtask(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedSubtask := &wundergo.Subtask{
				Title: "Test Title",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedSubtask, nil)
			})

			It("returns the unmarshalled subtask without error", func() {
				subtask, err := client.Subtask(taskID)

				Expect(err).To(BeNil())
				Expect(subtask).To(Equal(expectedSubtask))
			})
		})
	})

	Describe("Creating a new subtask", func() {
		subtaskTitle := "newSubtaskTitle"
		taskID := uint(1)

		var completed bool

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusCreated
			fakeHTTPHelper.PostReturns(dummyResponse, nil)

			completed = false
		})

		It("performs POST requests to /subtasks with new list title in body", func() {
			expectedUrl := fmt.Sprintf("%s/subtasks", apiUrl)
			expectedBody := []byte(fmt.Sprintf(`{"title":"%s","task_id":%d,"completed":%t}`, subtaskTitle, taskID, completed))

			fakeJSONHelper.UnmarshalReturns(&wundergo.Subtask{}, nil)
			client.CreateSubtask(subtaskTitle, taskID, completed)

			Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PostArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when taskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.CreateSubtask(
					subtaskTitle,
					taskID,
					completed,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when httpHelper.Post returns an error", func() {
			expectedError := errors.New("httpHelper POST error")

			BeforeEach(func() {
				fakeHTTPHelper.PostReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateSubtask(
					subtaskTitle,
					taskID,
					completed,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CreateSubtask(
					subtaskTitle,
					taskID,
					completed,
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
				_, err := client.CreateSubtask(
					subtaskTitle,
					taskID,
					completed,
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
				_, err := client.CreateSubtask(
					subtaskTitle,
					taskID,
					completed,
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
				_, err := client.CreateSubtask(
					subtaskTitle,
					taskID,
					completed,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedSubtask := &wundergo.Subtask{
				Title: "Test Title",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedSubtask, nil)
			})

			It("returns the unmarshalled subtask without error", func() {
				subtask, err := client.CreateSubtask(
					subtaskTitle,
					taskID,
					completed,
				)

				Expect(err).To(BeNil())
				Expect(subtask).To(Equal(expectedSubtask))
			})
		})
	})

	Describe("Updating a subtask", func() {
		var subtask wundergo.Subtask

		BeforeEach(func() {
			subtask = wundergo.Subtask{
				ID:       uint(1),
				Revision: 2,
			}

			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.PatchReturns(dummyResponse, nil)
		})

		It("performs PATCH requests to /subtasks/:id", func() {
			expectedBody := []byte{}
			fakeJSONHelper.MarshalReturns(expectedBody, nil)
			fakeJSONHelper.UnmarshalReturns(&wundergo.Subtask{}, nil)
			expectedUrl := fmt.Sprintf("%s/subtasks/%d", apiUrl, subtask.ID)

			client.UpdateSubtask(subtask)

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
				_, err := client.UpdateSubtask(subtask)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when httpHelper.Patch returns an error", func() {
			expectedError := errors.New("httpHelper PATCH error")

			BeforeEach(func() {
				fakeHTTPHelper.PatchReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateSubtask(subtask)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns error", func() {
				_, err := client.UpdateSubtask(subtask)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PatchReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.UpdateSubtask(subtask)

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
				_, err := client.UpdateSubtask(subtask)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateSubtask(subtask)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUpdateSubtask := &wundergo.Subtask{
				Title: "Updated Title",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedUpdateSubtask, nil)
			})

			It("returns the unmarshalled subtask without error", func() {
				subtask, err := client.UpdateSubtask(subtask)

				Expect(err).To(BeNil())
				Expect(subtask).To(Equal(expectedUpdateSubtask))
			})
		})
	})

	Describe("Deleting a subtask", func() {
		subtask := wundergo.Subtask{
			ID:       uint(1),
			Revision: 3,
		}

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusNoContent
			fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
		})

		It("performs DELETE requests to /subtasks/:id?revision=:revision", func() {
			expectedUrl := fmt.Sprintf("%s/subtasks/%d?revision=%d", apiUrl, subtask.ID, subtask.Revision)

			client.DeleteSubtask(subtask)

			Expect(fakeHTTPHelper.DeleteCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.DeleteArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Delete returns an error", func() {
			expectedError := errors.New("httpHelper DELETE error")

			BeforeEach(func() {
				fakeHTTPHelper.DeleteReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				err := client.DeleteSubtask(subtask)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns error", func() {
				err := client.DeleteSubtask(subtask)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				err := client.DeleteSubtask(subtask)

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
				err := client.DeleteSubtask(subtask)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			It("deletes the subtask without error", func() {
				err := client.DeleteSubtask(subtask)

				Expect(err).To(BeNil())
			})
		})
	})
})
