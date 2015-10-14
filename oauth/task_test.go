package oauth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wl"
	"github.com/robdimsdale/wl/oauth"
)

var _ = Describe("client - Task operations", func() {
	Describe("getting tasks for list", func() {
		var (
			listID uint
		)

		BeforeEach(func() {
			listID = 1234
		})

		It("performs GET requests with correct headers to /tasks", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/tasks", "list_id=1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.TasksForListID(listID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedTasks := []wl.Task{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedTasks)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedTasks)
				Expect(err).NotTo(HaveOccurred())

				// Actually marshal into the transport type
				expectedBody, err = json.Marshal(transportsFromTasks(expectedTasks))
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)
				tasks, err := client.TasksForListID(listID)
				Expect(err).NotTo(HaveOccurred())

				Expect(tasks).To(Equal(expectedTasks))
			})
		})

		Context("when ListID == 0", func() {
			BeforeEach(func() {
				listID = 0
			})

			It("returns an error", func() {
				_, err := client.TasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.TasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.TasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response status code is unexpected", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNotFound, nil),
					),
				)

				_, err := client.TasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)

				_, err := client.TasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, "invalid json response"),
					),
				)

				_, err := client.TasksForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting completed tasks for list", func() {
		var (
			listID    uint
			completed bool
		)

		BeforeEach(func() {
			listID = 1234
			completed = true
		})

		It("performs GET requests with correct headers to /tasks", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/tasks", "list_id=1234&completed=true"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.CompletedTasksForListID(listID, completed)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedTasks := []wl.Task{{ID: 2345}}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedTasks)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedTasks)
				Expect(err).NotTo(HaveOccurred())

				// Actually marshal into the transport type
				expectedBody, err = json.Marshal(transportsFromTasks(expectedTasks))
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				tasks, err := client.CompletedTasksForListID(listID, completed)
				Expect(err).NotTo(HaveOccurred())

				Expect(tasks).To(Equal(expectedTasks))
			})
		})

		Context("when ListID == 0", func() {
			BeforeEach(func() {
				listID = 0
			})

			It("returns an error", func() {
				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response status code is unexpected", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNotFound, nil),
					),
				)

				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)

				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, "invalid json response"),
					),
				)

				_, err := client.CompletedTasksForListID(listID, completed)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("getting task by ID", func() {
		var taskID uint

		BeforeEach(func() {
			taskID = 1234
		})

		It("performs DELETE requests with correct headers to /task/:id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/tasks/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.Task(taskID)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedTask := wl.Task{
					ID: taskID,
				}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedTask)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedTask)
				Expect(err).NotTo(HaveOccurred())

				// Actually marshal into the transport type
				expectedBody, err = json.Marshal(transportFromTask(expectedTask))
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, expectedBody),
					),
				)

				task, err := client.Task(taskID)
				Expect(err).NotTo(HaveOccurred())

				Expect(task).To(Equal(expectedTask))
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Task(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.Task(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response status code is unexpected", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNotFound, nil),
					),
				)

				_, err := client.Task(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)

				_, err := client.Task(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, "invalid json response"),
					),
				)

				_, err := client.Task(taskID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("creating a new task", func() {
		var (
			title           string
			listID          uint
			assigneeID      uint
			completed       bool
			recurrenceType  string
			recurrenceCount uint
			dueDate         time.Time
			starred         bool
		)

		BeforeEach(func() {
			title = "newTaskTitle"
			listID = uint(1)

			assigneeID = uint(2)
			completed = true
			recurrenceType = "day"
			recurrenceCount = uint(3)
			dueDate = time.Date(1968, 1, 2, 0, 0, 0, 0, time.UTC)
			starred = true
		})

		It("performs POST requests with correct headers to /tasks", func() {
			expectedBody := fmt.Sprintf(`{
				"title":"%s",
				"list_id":%d,
				"completed":%t,
				"assignee_id":%d,
				"recurrence_type":"%s",
				"recurrence_count":%d,
				"due_date":"%s",
				"starred":%t
			}`,
				title,
				listID,
				completed,
				assigneeID,
				recurrenceType,
				recurrenceCount,
				"1968-01-02",
				starred,
			)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/tasks"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSON(expectedBody),
				),
			)

			client.CreateTask(
				title,
				listID,
				assigneeID,
				completed,
				recurrenceType,
				recurrenceCount,
				dueDate,
				starred,
			)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				expectedTask := wl.Task{ID: 2345}

				// Marshal and unmarshal to ensure exact object is returned
				// - this avoids odd behavior with the time fields
				expectedBody, err := json.Marshal(expectedTask)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(expectedBody, &expectedTask)
				Expect(err).NotTo(HaveOccurred())

				// Actually marshal into the transport type
				expectedBody, err = json.Marshal(transportFromTask(expectedTask))
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusCreated, expectedBody),
					),
				)

				actualTask, err := client.CreateTask(
					title,
					listID,
					assigneeID,
					completed,
					recurrenceType,
					recurrenceCount,
					dueDate,
					starred,
				)

				Expect(err).NotTo(HaveOccurred())

				Expect(actualTask).To(Equal(expectedTask))
			})
		})

		Context("when listID == 0", func() {
			BeforeEach(func() {
				listID = 0
			})

			It("returns an error", func() {
				_, err := client.CreateTask(
					title,
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

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("returns an error", func() {
				_, err := client.CreateTask(
					title,
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

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.CreateTask(
					title,
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

		Context("when response status code is unexpected", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNotFound, nil),
					),
				)

				_, err := client.CreateTask(
					title,
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
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)

				_, err := client.CreateTask(
					title,
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

		Context("when unmarshalling json response returns an error", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, "invalid json response"),
					),
				)

				_, err := client.CreateTask(
					title,
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
	})

	Describe("updating a task", func() {
		var (
			originalTask transportTask
			expectedTask transportTask

			task   wl.Task
			taskID uint

			title           string
			revision        uint
			assigneeID      uint
			completed       bool
			recurrenceType  string
			recurrenceCount uint
			dueDate         time.Time
			starred         bool

			expectedTaskUpdateConfig oauth.TaskUpdateConfig
		)

		BeforeEach(func() {
			taskID = 1234
			title = "newTaskTitle"
			revision = 23
			assigneeID = uint(2)
			completed = false
			recurrenceType = "day"
			recurrenceCount = uint(3)
			dueDate = time.Date(1968, 1, 2, 0, 0, 0, 0, time.UTC)
			starred = false

			task = wl.Task{
				ID:              taskID,
				Title:           title,
				Revision:        revision,
				AssigneeID:      assigneeID,
				Completed:       completed,
				RecurrenceType:  recurrenceType,
				RecurrenceCount: recurrenceCount,
				DueDate:         dueDate,
				Starred:         starred,
			}

			originalTask = transportFromTask(task)
			expectedTask = transportFromTask(task)

			expectedTaskUpdateConfig = oauth.TaskUpdateConfig{
				Title:           title,
				Revision:        revision,
				AssigneeID:      assigneeID,
				Completed:       completed,
				RecurrenceType:  recurrenceType,
				RecurrenceCount: recurrenceCount,
				DueDate:         "1968-01-02",
				Starred:         starred,
				Remove:          []string{},
			}
		})

		JustBeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/tasks/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusOK, originalTask),
				),
			)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "/tasks/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
					ghttp.VerifyJSONRepresenting(expectedTaskUpdateConfig),
					ghttp.RespondWithJSONEncoded(http.StatusOK, expectedTask),
				),
			)
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
			Context("when initial GET request returns with error", func() {
				BeforeEach(func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.RespondWith(http.StatusOK, nil),
						),
					)
				})

				It("returns an error", func() {
					_, err := client.UpdateTask(task)

					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("assigneeID updates", func() {
			Context("when original task had an assignee", func() {
				BeforeEach(func() {
					originalTask.AssigneeID = 1
				})

				Context("and new task has no assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = 0

						expectedTaskUpdateConfig.AssigneeID = 0
						expectedTaskUpdateConfig.Remove = []string{"assignee_id"}

						expectedTask.AssigneeID = 0
					})

					It("removes assigneeID", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.AssigneeID).To(Equal(uint(0)))
					})
				})

				Context("and new task has same assignee", func() {
					It("retains the same value of assigneeID", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.AssigneeID).To(Equal(task.AssigneeID))
					})
				})

				Context("and new task has different assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = 2

						expectedTaskUpdateConfig.AssigneeID = 2

						expectedTask.AssigneeID = 2
					})

					It("changes assigneeID", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.AssigneeID).To(Equal(uint(2)))
					})
				})
			})

			Context("when original task had no assignee", func() {
				BeforeEach(func() {
					originalTask.AssigneeID = 0
				})

				Context("and new task has no assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = 0

						expectedTaskUpdateConfig.AssigneeID = 0

						expectedTask.AssigneeID = 0
					})

					It("sets assigneeID to 0", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.AssigneeID).To(Equal(uint(0)))
					})
				})

				Context("and new task has an assignee", func() {
					BeforeEach(func() {
						task.AssigneeID = 2

						expectedTaskUpdateConfig.AssigneeID = 2

						expectedTask.AssigneeID = 2
					})

					It("adds assigneeID", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.AssigneeID).To(Equal(uint(2)))
					})
				})
			})
		})

		Describe("dueDate updates", func() {
			Context("when original task had a due date", func() {
				BeforeEach(func() {
					originalTask.DueDate = "1970-01-01"
				})

				Context("and new task has empty due date", func() {
					BeforeEach(func() {
						task.DueDate = time.Time{}

						expectedTaskUpdateConfig.DueDate = ""
						expectedTaskUpdateConfig.Remove = []string{"due_date"}

						expectedTask.DueDate = ""
					})

					It("removes dueDate", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.DueDate).To(Equal(time.Time{}))
					})
				})

				Context("and new task has same due date", func() {
					It("retains the same value of dueDate", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.DueDate).To(Equal(task.DueDate))
					})
				})

				Context("and new task has different due date", func() {
					BeforeEach(func() {
						task.DueDate = time.Date(1921, 12, 24, 0, 0, 0, 0, time.UTC)

						expectedTaskUpdateConfig.DueDate = "1921-12-24"

						expectedTask.DueDate = "1921-12-24"
					})

					It("changes dueDate", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.DueDate).To(Equal(task.DueDate))
					})
				})
			})

			Context("when original task had no due date", func() {
				BeforeEach(func() {
					originalTask.DueDate = ""
				})

				Context("and new task has empty due date", func() {
					BeforeEach(func() {
						task.DueDate = time.Time{}

						expectedTaskUpdateConfig.DueDate = ""

						expectedTask.DueDate = ""
					})

					It("retains no dueDate", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.DueDate).To(Equal(time.Time{}))
					})
				})

				Context("and new task has different due date", func() {
					BeforeEach(func() {
						task.DueDate = time.Date(1921, 12, 24, 0, 0, 0, 0, time.UTC)

						expectedTaskUpdateConfig.DueDate = "1921-12-24"

						expectedTask.DueDate = "1921-12-24"
					})

					It("changes dueDate", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.DueDate).To(Equal(task.DueDate))
					})
				})
			})
		})

		Describe("recurrence updates", func() {
			Context("when original task had recurrence", func() {
				BeforeEach(func() {
					originalTask.RecurrenceType = "day"
					originalTask.RecurrenceCount = uint(1)
				})

				Context("and recurrence is unset", func() {
					BeforeEach(func() {
						task.RecurrenceType = ""
						expectedTaskUpdateConfig.RecurrenceType = ""
						expectedTask.RecurrenceType = ""

						task.RecurrenceCount = 0
						expectedTaskUpdateConfig.RecurrenceCount = 0
						expectedTask.RecurrenceCount = 0

						expectedTaskUpdateConfig.Remove = []string{"recurrence_type", "recurrence_count"}
					})

					It("removes recurrence", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.RecurrenceType).To(Equal(""))
						Expect(actualTask.RecurrenceCount).To(Equal(uint(0)))
					})
				})

				Context("and recurrence type is now set to a different value", func() {
					BeforeEach(func() {
						task.RecurrenceType = "week"
						expectedTaskUpdateConfig.RecurrenceType = "week"
						expectedTask.RecurrenceType = "week"

						// Unchanged
						task.RecurrenceCount = originalTask.RecurrenceCount
						expectedTaskUpdateConfig.RecurrenceCount = originalTask.RecurrenceCount
						expectedTask.RecurrenceCount = originalTask.RecurrenceCount
					})

					It("updates recurrence type", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.RecurrenceType).To(Equal("week"))
					})

					It("does not change recurrence count", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.RecurrenceCount).To(Equal(uint(task.RecurrenceCount)))
					})
				})

				Context("and recurrence count is now set to a different value", func() {
					BeforeEach(func() {
						task.RecurrenceCount = uint(2)
						expectedTaskUpdateConfig.RecurrenceCount = 2
						expectedTask.RecurrenceCount = 2

						// Unchanged
						task.RecurrenceType = originalTask.RecurrenceType
						expectedTaskUpdateConfig.RecurrenceType = originalTask.RecurrenceType
						expectedTask.RecurrenceType = originalTask.RecurrenceType
					})

					It("updates recurrence count", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.RecurrenceCount).To(Equal(uint(2)))
					})

					It("does not change recurrence type", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.RecurrenceType).To(Equal(task.RecurrenceType))
					})
				})

				Context("and recurrence is now set to the same value", func() {
					BeforeEach(func() {
						task.RecurrenceCount = originalTask.RecurrenceCount
						expectedTaskUpdateConfig.RecurrenceCount = originalTask.RecurrenceCount
						expectedTask.RecurrenceCount = originalTask.RecurrenceCount

						task.RecurrenceType = originalTask.RecurrenceType
						expectedTaskUpdateConfig.RecurrenceType = originalTask.RecurrenceType
						expectedTask.RecurrenceType = originalTask.RecurrenceType
					})

					It("does not change recurrence count or type", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.RecurrenceCount).To(Equal(task.RecurrenceCount))
						Expect(actualTask.RecurrenceType).To(Equal(task.RecurrenceType))
					})
				})
			})

			Context("when original task had no recurrence", func() {
				BeforeEach(func() {
					originalTask.RecurrenceType = ""
					originalTask.RecurrenceCount = 0
				})

				Context("and recurrence is now set", func() {
					BeforeEach(func() {
						task.RecurrenceCount = 12
						expectedTaskUpdateConfig.RecurrenceCount = 12
						expectedTask.RecurrenceCount = 12

						task.RecurrenceType = "week"
						expectedTaskUpdateConfig.RecurrenceType = "week"
						expectedTask.RecurrenceType = "week"
					})

					It("updates recurrence count and type", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.RecurrenceCount).To(Equal(task.RecurrenceCount))
						Expect(actualTask.RecurrenceType).To(Equal(task.RecurrenceType))
					})
				})

				Context("and recurrence is still unset", func() {
					BeforeEach(func() {
						task.RecurrenceCount = originalTask.RecurrenceCount
						expectedTaskUpdateConfig.RecurrenceCount = originalTask.RecurrenceCount
						expectedTask.RecurrenceCount = originalTask.RecurrenceCount

						task.RecurrenceType = originalTask.RecurrenceType
						expectedTaskUpdateConfig.RecurrenceType = originalTask.RecurrenceType
						expectedTask.RecurrenceType = originalTask.RecurrenceType
					})

					It("does not change recurrence count or type", func() {
						actualTask, err := client.UpdateTask(task)
						Expect(err).NotTo(HaveOccurred())

						Expect(actualTask.RecurrenceCount).To(Equal(task.RecurrenceCount))
						Expect(actualTask.RecurrenceType).To(Equal(task.RecurrenceType))
					})
				})
			})
		})

		Describe("completed updates", func() {
			Context("when completed state changes", func() {
				BeforeEach(func() {
					originalTask.Completed = false

					task.Completed = true
					expectedTaskUpdateConfig.Completed = true
					expectedTask.Completed = true
				})

				It("changes completed state", func() {
					actualTask, err := client.UpdateTask(task)
					Expect(err).NotTo(HaveOccurred())

					Expect(actualTask.Completed).To(Equal(true))
				})
			})

			Context("when completed state is unchanged", func() {
				BeforeEach(func() {
					originalTask.Completed = true

					task.Completed = true
					expectedTaskUpdateConfig.Completed = true
					expectedTask.Completed = true
				})

				It("does not change completed state", func() {
					actualTask, err := client.UpdateTask(task)
					Expect(err).NotTo(HaveOccurred())

					Expect(actualTask.Completed).To(Equal(true))
				})
			})
		})

		Describe("starred updates", func() {
			Context("when starred state changes", func() {
				BeforeEach(func() {
					originalTask.Starred = false

					task.Starred = true
					expectedTaskUpdateConfig.Starred = true
					expectedTask.Starred = true
				})

				It("changes starred state", func() {
					actualTask, err := client.UpdateTask(task)
					Expect(err).NotTo(HaveOccurred())

					Expect(actualTask.Starred).To(Equal(true))
				})
			})

			Context("when starred state is unchanged", func() {
				BeforeEach(func() {
					originalTask.Starred = true

					task.Starred = true
					expectedTaskUpdateConfig.Starred = true
					expectedTask.Starred = true
				})

				It("does not change starred state", func() {
					actualTask, err := client.UpdateTask(task)
					Expect(err).NotTo(HaveOccurred())

					Expect(actualTask.Starred).To(Equal(true))
				})
			})
		})

		Describe("title updates", func() {
			Context("when title changes", func() {
				BeforeEach(func() {
					originalTask.Title = "Old Title"

					task.Title = "new title"
					expectedTaskUpdateConfig.Title = "new title"
					expectedTask.Title = "new title"
				})

				It("changes title", func() {
					actualTask, err := client.UpdateTask(task)
					Expect(err).NotTo(HaveOccurred())

					Expect(actualTask.Title).To(Equal("new title"))
				})
			})

			Context("when title is unchanged", func() {
				BeforeEach(func() {
					originalTask.Title = "Old Title"

					task.Title = originalTask.Title
					expectedTaskUpdateConfig.Title = originalTask.Title
					expectedTask.Title = originalTask.Title
				})

				It("does not change title", func() {
					actualTask, err := client.UpdateTask(task)
					Expect(err).NotTo(HaveOccurred())

					Expect(actualTask.Title).To(Equal(originalTask.Title))
				})
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/tasks/1234"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, originalTask),
					),
				)

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PATCH", "/tasks/1234"),
						ghttp.RespondWith(http.StatusNotFound, nil),
					),
				)
			})

			It("returns an error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/tasks/1234"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, originalTask),
					),
				)

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PATCH", "/tasks/1234"),
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)
			})

			It("returns an error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/tasks/1234"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, originalTask),
					),
				)

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PATCH", "/tasks/1234"),
						ghttp.RespondWith(http.StatusOK, "invalid json response"),
					),
				)
			})

			It("returns an error", func() {
				_, err := client.UpdateTask(task)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("deleting a task", func() {
		var task wl.Task

		BeforeEach(func() {
			task = wl.Task{ID: 1234, Revision: 23}
		})

		It("performs DELETE requests with correct headers to /task/:id", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/tasks/1234"),
					ghttp.VerifyHeader(http.Header{
						"X-Access-Token": []string{dummyAccessToken},
						"X-Client-ID":    []string{dummyClientID},
					}),
				),
			)

			client.DeleteTask(task)

			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("when the request is valid", func() {
			It("returns successfully", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNoContent, nil),
					),
				)

				err := client.DeleteTask(task)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when creating request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "", testLogger)
			})

			It("forwards the error", func() {
				err := client.DeleteTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when executing request fails with error", func() {
			BeforeEach(func() {
				client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
			})

			It("forwards the error", func() {
				err := client.DeleteTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response status code is unexpected", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusNotFound, nil),
					),
				)

				err := client.DeleteTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)

				err := client.DeleteTask(task)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			It("returns an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, "invalid json response"),
					),
				)

				err := client.DeleteTask(task)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})

type transportTask struct {
	ID              uint      `json:"id" yaml:"id"`
	AssigneeID      uint      `json:"assignee_id" yaml:"assignee_id"`
	AssignerID      uint      `json:"assigner_id" yaml:"assigner_id"`
	CreatedAt       time.Time `json:"created_at" yaml:"created_at"`
	CreatedByID     uint      `json:"created_by_id" yaml:"created_by_id"`
	DueDate         string    `json:"due_date" yaml:"due_date"`
	ListID          uint      `json:"list_id" yaml:"list_id"`
	Revision        uint      `json:"revision" yaml:"revision"`
	Starred         bool      `json:"starred" yaml:"starred"`
	Title           string    `json:"title" yaml:"title"`
	Completed       bool      `json:"completed" yaml:"completed"`
	CompletedAt     time.Time `json:"completed_at" yaml:"completed_at"`
	CompletedByID   uint      `json:"completed_by" yaml:"completed_by"`
	RecurrenceType  string    `json:"recurrence_type" yaml:"recurrence_type"`
	RecurrenceCount uint      `json:"recurrence_count" yaml:"recurrence_count"`
}

func transportsFromTasks(tasks []wl.Task) []transportTask {
	transportTasks := make([]transportTask, len(tasks))

	for i, t := range tasks {
		transportTasks[i] = transportFromTask(t)
	}

	return transportTasks
}

func transportFromTask(t wl.Task) transportTask {
	return transportTask{
		ID:              t.ID,
		AssigneeID:      t.AssigneeID,
		AssignerID:      t.AssignerID,
		CreatedAt:       t.CreatedAt,
		CreatedByID:     t.CreatedByID,
		DueDate:         dueDateToString(t.DueDate),
		ListID:          t.ListID,
		Revision:        t.Revision,
		Starred:         t.Starred,
		Title:           t.Title,
		Completed:       t.Completed,
		CompletedAt:     t.CompletedAt,
		CompletedByID:   t.CompletedByID,
		RecurrenceType:  t.RecurrenceType,
		RecurrenceCount: t.RecurrenceCount,
	}
}

func dueDateToString(dueDate time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", dueDate.Year(), dueDate.Month(), dueDate.Day())
}
