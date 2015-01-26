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

var _ = Describe("Client - File operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("Getting files for list", func() {
		listID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when ListID == 0", func() {
			listID := uint(0)

			It("returns an error", func() {
				_, err := client.FilesForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /files?list_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/files?list_id=%d", apiURL, listID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.File{}, nil)
			client.FilesForListID(listID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.FilesForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.FilesForListID(listID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.FilesForListID(listID)

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
				_, err := client.FilesForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.FilesForListID(listID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedFiles := &[]wundergo.File{
				wundergo.File{
					URL: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedFiles, nil)
			})

			It("returns the unmarshalled files without error", func() {
				file, err := client.FilesForListID(listID)

				Expect(err).To(BeNil())
				Expect(file).To(Equal(expectedFiles))
			})
		})
	})

	Describe("Getting files for task", func() {
		taskID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		Context("when TaskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.FilesForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		It("performs GET requests to /files?task_id=:id", func() {
			expectedUrl := fmt.Sprintf("%s/files?task_id=%d", apiURL, taskID)

			fakeJSONHelper.UnmarshalReturns(&[]wundergo.File{}, nil)
			client.FilesForTaskID(taskID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.FilesForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.FilesForTaskID(taskID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.FilesForTaskID(taskID)

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
				_, err := client.FilesForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.FilesForTaskID(taskID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedFiles := &[]wundergo.File{
				wundergo.File{
					URL: "testy",
				},
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedFiles, nil)
			})

			It("returns the unmarshalled f without error", func() {
				file, err := client.FilesForTaskID(taskID)

				Expect(err).To(BeNil())
				Expect(file).To(Equal(expectedFiles))
			})
		})
	})

	Describe("getting file by ID", func() {
		fileID := uint(1)
		expectedUrl := fmt.Sprintf("%s/files/%d", apiURL, fileID)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /files/:id", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.File{}, nil)
			client.File(fileID)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.GetArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper GET error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.File(fileID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.File(fileID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.File(fileID)

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
				_, err := client.File(fileID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.File(fileID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedFile := &wundergo.File{
				URL: "testy",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedFile, nil)
			})

			It("returns the unmarshalled file without error", func() {
				file, err := client.File(fileID)

				Expect(err).To(BeNil())
				Expect(file).To(Equal(expectedFile))
			})
		})
	})

	Describe("creating a file", func() {
		var uploadID uint
		var taskID uint
		var localCreatedAt string

		BeforeEach(func() {
			uploadID = 1
			taskID = 1
			localCreatedAt = ""

			dummyResponse.StatusCode = http.StatusCreated
			fakeHTTPHelper.PostReturns(dummyResponse, nil)
		})

		It("performs POST requests to /files?upload_id=:uploadID&task_id=:taskID", func() {
			expectedUrl := fmt.Sprintf("%s/files?upload_id=%d&task_id=%d", apiURL, uploadID, taskID)

			fakeJSONHelper.UnmarshalReturns(&wundergo.File{}, nil)
			client.CreateFile(uploadID, taskID, localCreatedAt)

			Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
			arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
		})

		Context("when uploadID == 0", func() {
			uploadID := uint(0)

			It("returns an error", func() {
				_, err := client.CreateFile(uploadID, taskID, localCreatedAt)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when taskID == 0", func() {
			taskID := uint(0)

			It("returns an error", func() {
				_, err := client.CreateFile(uploadID, taskID, localCreatedAt)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when localCreatedAt is empty", func() {
			It("does not include local_created_at in the url params", func() {
				expectedUrl := fmt.Sprintf("%s/files?upload_id=%d&task_id=%d", apiURL, uploadID, taskID)
				fakeJSONHelper.UnmarshalReturns(&wundergo.File{}, nil)

				client.CreateFile(uploadID, taskID, localCreatedAt)

				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))
			})
		})

		Context("when localCreatedAt is not empty", func() {
			BeforeEach(func() {
				localCreatedAt = "some_time"
			})

			It("includes local_created_at in the url params", func() {
				fakeJSONHelper.UnmarshalReturns(&wundergo.File{}, nil)
				expectedUrl := fmt.Sprintf("%s/files?upload_id=%d&task_id=%d&local_created_at=%s", apiURL, uploadID, taskID, localCreatedAt)

				client.CreateFile(uploadID, taskID, localCreatedAt)
				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))

			})
		})

		Context("when httpHelper.Post returns an error", func() {
			expectedError := errors.New("httpHelper PATCH error")

			BeforeEach(func() {
				fakeHTTPHelper.PostReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateFile(uploadID, taskID, localCreatedAt)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CreateFile(uploadID, taskID, localCreatedAt)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.CreateFile(uploadID, taskID, localCreatedAt)

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
				_, err := client.CreateFile(uploadID, taskID, localCreatedAt)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateFile(uploadID, taskID, localCreatedAt)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedFile := &wundergo.File{
				URL: "url",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedFile, nil)
			})

			It("returns the unmarshalled upload without error", func() {
				file, err := client.CreateFile(uploadID, taskID, localCreatedAt)

				Expect(err).To(BeNil())
				Expect(file).To(Equal(expectedFile))
			})
		})
	})

	Describe("destroying a file", func() {
		file := wundergo.File{
			ID:       1,
			Revision: 3,
		}

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusNoContent
			fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
		})

		It("performs DELETE requests to /files/:id?revision=:revision", func() {
			expectedUrl := fmt.Sprintf("%s/files/%d?revision=%d", apiURL, file.ID, file.Revision)

			client.DestroyFile(file)

			Expect(fakeHTTPHelper.DeleteCallCount()).To(Equal(1))
			Expect(fakeHTTPHelper.DeleteArgsForCall(0)).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Delete returns an error", func() {
			expectedError := errors.New("httpHelper DELETE error")

			BeforeEach(func() {
				fakeHTTPHelper.DeleteReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				err := client.DestroyFile(file)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				err := client.DestroyFile(file)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.DeleteReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				err := client.DestroyFile(file)

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
				err := client.DestroyFile(file)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			It("deletes the note without error", func() {
				err := client.DestroyFile(file)

				Expect(err).To(BeNil())
			})
		})
	})
})
