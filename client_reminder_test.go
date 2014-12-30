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

var _ = Describe("Client - Reminder operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("Getting reminders for list", func() {
		listID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.RemindersForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /reminders?list_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/reminders?list_id=%d", apiUrl, listID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Reminder{}, nil)
			client.RemindersForListID(listID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.RemindersForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.RemindersForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.RemindersForListID(listID)

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
				_, err := client.RemindersForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.RemindersForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedReminders := &[]wundergo.Reminder{
				wundergo.Reminder{
					Date: "some-date",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedReminders, nil)
			})

			It("returns the unmarshalled reminders without error", func() {
				Reminder, err := client.RemindersForListID(listID)

				Expect(err).To(BeNil())
				Expect(Reminder).To(Equal(expectedReminders))
			})
		})
	})

	Describe("Getting reminders for task", func() {
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when TaskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.RemindersForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /reminders?task_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/reminders?task_id=%d", apiUrl, taskID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.Reminder{}, nil)
			client.RemindersForTaskID(taskID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.RemindersForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.RemindersForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.RemindersForTaskID(taskID)

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
				_, err := client.RemindersForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.RemindersForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedReminders := &[]wundergo.Reminder{
				wundergo.Reminder{
					Date: "some-date",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedReminders, nil)
			})

			It("returns the unmarshalled reminders without error", func() {
				Reminder, err := client.RemindersForTaskID(taskID)

				Expect(err).To(BeNil())
				Expect(Reminder).To(Equal(expectedReminders))
			})
		})
	})

	Describe("Getting Reminder by ID", func() {
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /reminders/:id", func() {
			expectedUrl := fmt.Sprintf("%s/reminders/%d", apiUrl, taskID)
			fakeJSONHelper.UnmarshalReturns(&wundergo.Reminder{}, nil)
			client.Reminder(taskID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Reminder(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.Reminder(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.Reminder(taskID)

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
				_, err := client.Reminder(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.Reminder(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedReminder := &wundergo.Reminder{
				Date: "some-date",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedReminder, nil)
			})

			It("returns the unmarshalled Reminder without error", func() {
				Reminder, err := client.Reminder(taskID)

				Expect(err).To(BeNil())
				Expect(Reminder).To(Equal(expectedReminder))
			})
		})
	})

	Describe("Creating a new Reminder", func() {
		reminderDate := "new-reminder-date"
		var createdByDeviceUdid string
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusCreated
			fakeHTTPHelper.PostReturns(dummyResponse, nil)

			createdByDeviceUdid = "some-device-id"
		})

		It("performs POST requests to /reminders with new list title in body", func() {
			expectedUrl := fmt.Sprintf("%s/reminders", apiUrl)
			expectedBody := []byte(fmt.Sprintf(`{"date":"%s","task_id":%d,"created_by_device_udid":%s}`, reminderDate, taskID, createdByDeviceUdid))

			fakeJSONHelper.UnmarshalReturns(&wundergo.Reminder{}, nil)
			client.CreateReminder(
				reminderDate,
				taskID,
				createdByDeviceUdid,
			)

			Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
			arg0, arg1 := fakeHTTPHelper.PostArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
			Expect(arg1).To(Equal(expectedBody))
		})

		Context("when taskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.CreateReminder(
					reminderDate,
					taskID,
					createdByDeviceUdid,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when createdByDeviceUdid is empty", func() {
			BeforeEach(func() {
				createdByDeviceUdid = ""
			})

			It("does not include created_by_device_udid in the params", func() {
				expectedBody := []byte(fmt.Sprintf(`{"date":"%s","task_id":%d}`, reminderDate, taskID))

				fakeJSONHelper.UnmarshalReturns(&wundergo.Reminder{}, nil)
				client.CreateReminder(
					reminderDate,
					taskID,
					createdByDeviceUdid,
				)

				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				_, arg1 := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg1).To(Equal(expectedBody))
			})
		})

		Context("when httpHelper.Post returns an error", func() {
			expectedError := errors.New("httpHelper POST error")

			BeforeEach(func() {
				fakeHTTPHelper.PostReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateReminder(
					reminderDate,
					taskID,
					createdByDeviceUdid,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CreateReminder(
					reminderDate,
					taskID,
					createdByDeviceUdid,
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
				_, err := client.CreateReminder(
					reminderDate,
					taskID,
					createdByDeviceUdid,
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
				_, err := client.CreateReminder(
					reminderDate,
					taskID,
					createdByDeviceUdid,
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
				_, err := client.CreateReminder(
					reminderDate,
					taskID,
					createdByDeviceUdid,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedReminder := &wundergo.Reminder{
				Date: "some-date",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedReminder, nil)
			})

			It("returns the unmarshalled Reminder without error", func() {
				reminder, err := client.CreateReminder(
					reminderDate,
					taskID,
					createdByDeviceUdid,
				)

				Expect(err).To(BeNil())
				Expect(reminder).To(Equal(expectedReminder))
			})
		})
	})

	Describe("Updating a Reminder", func() {
		var Reminder wundergo.Reminder

		BeforeEach(func() {
			Reminder = wundergo.Reminder{
				ID:       uint(1),
				Revision: 2,
			}

			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.PatchReturns(dummyResponse, nil)
		})

		It("performs PATCH requests to /reminders/:id", func() {
			expectedBody := []byte{}
			fakeJSONHelper.MarshalReturns(expectedBody, nil)
			fakeJSONHelper.UnmarshalReturns(&wundergo.Reminder{}, nil)
			expectedUrl := fmt.Sprintf("%s/reminders/%d", apiUrl, Reminder.ID)

			client.UpdateReminder(Reminder)

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
				_, err := client.UpdateReminder(Reminder)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when httpHelper.Patch returns an error", func() {
			expectedError := errors.New("httpHelper PATCH error")

			BeforeEach(func() {
				fakeHTTPHelper.PatchReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateReminder(Reminder)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns error", func() {
				_, err := client.UpdateReminder(Reminder)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PatchReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.UpdateReminder(Reminder)

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
				_, err := client.UpdateReminder(Reminder)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.UpdateReminder(Reminder)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUpdateReminder := &wundergo.Reminder{
				Date: "some-updated-date",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedUpdateReminder, nil)
			})

			It("returns the unmarshalled Reminder without error", func() {
				Reminder, err := client.UpdateReminder(Reminder)

				Expect(err).To(BeNil())
				Expect(Reminder).To(Equal(expectedUpdateReminder))
			})
		})
	})

	Describe("Deleting a Reminder", func() {
		reminder := wundergo.Reminder{
			ID:       uint(1),
			Revision: 3,
		}

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusNoContent
			fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
		})

		It("performs DELETE requests to /reminders/:id?revision=:revision", func() {
			expectedUrl := fmt.Sprintf("%s/reminders/%d?revision=%d", apiUrl, reminder.ID, reminder.Revision)

			client.DeleteReminder(reminder)

			Expect(fakeHTTPHelper.DeleteCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.DeleteArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Delete returns an error", func() {
			expectedError := errors.New("httpHelper DELETE error")

			BeforeEach(func() {
				fakeHTTPHelper.DeleteReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				err := client.DeleteReminder(reminder)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns error", func() {
				err := client.DeleteReminder(reminder)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				err := client.DeleteReminder(reminder)

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
				err := client.DeleteReminder(reminder)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			It("deletes the Reminder without error", func() {
				err := client.DeleteReminder(reminder)

				Expect(err).To(BeNil())
			})
		})
	})
})
