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

var _ = Describe("client - Upload operations", func() {
	var dummyResponse *http.Response

	BeforeEach(func() {
		initializeFakes()
		initializeClient()

		dummyResponse = &http.Response{}
		dummyResponse.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	})

	Describe("creating an upload", func() {
		var contentType string
		var fileName string
		var fileSize uint
		var partNumber uint
		var md5Sum string

		BeforeEach(func() {
			contentType = "contentType"
			fileName = "filename"
			fileSize = 1
			partNumber = 0
			md5Sum = ""

			dummyResponse.StatusCode = http.StatusCreated
			fakeHTTPHelper.PostReturns(dummyResponse, nil)
		})

		It("performs POST requests to /uploads?content_type=:contentType&file_name=:fileName&file_size=:fileSize", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)
			expectedUrl := fmt.Sprintf("%s/uploads?content_type=%s&file_name=%s&file_size=%d", apiURL, contentType, fileName, fileSize)

			client.CreateUpload(
				contentType,
				fileName,
				fileSize,
				partNumber,
				md5Sum,
			)

			Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
			arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
		})

		Context("when contentType is empty", func() {
			BeforeEach(func() {
				contentType = ""
			})

			It("returns an error", func() {
				_, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when fileName is empty", func() {
			BeforeEach(func() {
				fileName = ""
			})

			It("returns an error", func() {
				_, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when fileSize == 0", func() {
			BeforeEach(func() {
				fileSize = 0
			})

			It("returns an error", func() {
				_, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when partNumber is 0", func() {
			BeforeEach(func() {
				partNumber = 0
			})

			It("does not include part_number in the url params", func() {
				fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)

				client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).ShouldNot(ContainSubstring("part_number"))
			})
		})

		Context("when partNumber is > 0", func() {
			BeforeEach(func() {
				partNumber = 1
			})

			It("includes part_number in the url params", func() {
				fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)
				expectedUrl := fmt.Sprintf("%s/uploads?content_type=%s&file_name=%s&file_size=%d&part_number=%d", apiURL, contentType, fileName, fileSize, partNumber)

				client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)
				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))
			})
		})

		Context("when md5sum is empty", func() {
			It("does not include md5_sum in the url params", func() {
				fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)

				client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).ShouldNot(ContainSubstring("md5_sum"))
			})
		})

		Context("when md5Sum is not empty", func() {
			BeforeEach(func() {
				md5Sum = "md5sum"
			})

			It("includes md5_sum in the url params", func() {
				fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)
				expectedUrl := fmt.Sprintf("%s/uploads?content_type=%s&file_name=%s&file_size=%d&md5_sum=%s", apiURL, contentType, fileName, fileSize, md5Sum)

				client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)
				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))
			})
		})

		Context("when partNumber is > 0 and md5Sum is not empty", func() {
			BeforeEach(func() {
				partNumber = 1
				md5Sum = "md5Sum"
			})

			It("includes part_number and md5_sum in the url params", func() {
				fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)
				expectedUrl := fmt.Sprintf("%s/uploads?content_type=%s&file_name=%s&file_size=%d&part_number=%d&md5_sum=%s", apiURL, contentType, fileName, fileSize, partNumber, md5Sum)

				client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)
				Expect(fakeHTTPHelper.PostCallCount()).To(Equal(1))
				arg0, _ := fakeHTTPHelper.PostArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))
			})
		})

		Context("when httpHelper.Post returns an error", func() {
			expectedError := errors.New("httpHelper POST error")

			BeforeEach(func() {
				fakeHTTPHelper.PostReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
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
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("forwards the error", func() {
				_, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
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
				_, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUpload := &wundergo.Upload{
				State: "testy",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedUpload, nil)
			})

			It("returns the unmarshalled upload without error", func() {
				upload, err := client.CreateUpload(
					contentType,
					fileName,
					fileSize,
					partNumber,
					md5Sum,
				)

				Expect(err).To(BeNil())
				Expect(upload).To(Equal(expectedUpload))
			})
		})
	})

	Describe("chunked upload part", func() {
		var uploadID uint
		var partNumber uint
		var md5Sum string

		BeforeEach(func() {
			uploadID = 1
			partNumber = 1
			md5Sum = ""

			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.GetReturns(dummyResponse, nil)
		})

		It("performs GET requests to /uploads/:id/parts?part_number=:partNumber", func() {
			fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)
			expectedUrl := fmt.Sprintf("%s/uploads/%d/parts?part_number=%d", apiURL, uploadID, partNumber)

			client.ChunkedUploadPart(
				uploadID,
				partNumber,
				md5Sum,
			)

			Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
			arg0 := fakeHTTPHelper.GetArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
		})

		Context("when uploadID == 0", func() {
			BeforeEach(func() {
				uploadID = 0
			})

			It("returns an error", func() {
				_, err := client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when partNumber == 0", func() {
			BeforeEach(func() {
				partNumber = 0
			})

			It("returns an error", func() {
				_, err := client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when md5sum is empty", func() {
			It("does not include md5_sum in the url params", func() {
				fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)

				client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
				)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				arg0 := fakeHTTPHelper.GetArgsForCall(0)
				Expect(arg0).ShouldNot(ContainSubstring("md5_sum"))
			})
		})

		Context("when md5Sum is not empty", func() {
			BeforeEach(func() {
				md5Sum = "md5sum"
			})

			It("includes md5_sum in the url params", func() {
				fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)
				expectedUrl := fmt.Sprintf("%s/uploads/%d/parts?part_number=%d&md5_sum=%s", apiURL, uploadID, partNumber, md5Sum)

				client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
				)

				Expect(fakeHTTPHelper.GetCallCount()).To(Equal(1))
				arg0 := fakeHTTPHelper.GetArgsForCall(0)
				Expect(arg0).To(Equal(expectedUrl))
			})
		})

		Context("when httpHelper.Get returns an error", func() {
			expectedError := errors.New("httpHelper PATCH error")

			BeforeEach(func() {
				fakeHTTPHelper.GetReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
				)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
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
				fakeHTTPHelper.GetReturns(dummyResponse, nil)
			})

			It("forwards the error", func() {
				_, err := client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
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
				_, err := client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
				)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUpload := &wundergo.Upload{
				State: "finished",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedUpload, nil)
			})

			It("returns the unmarshalled upload without error", func() {
				upload, err := client.ChunkedUploadPart(
					uploadID,
					partNumber,
					md5Sum,
				)

				Expect(err).To(BeNil())
				Expect(upload).To(Equal(expectedUpload))
			})
		})
	})

	Describe("marking upload complete", func() {
		uploadID := uint(1)

		BeforeEach(func() {
			dummyResponse.StatusCode = http.StatusOK
			fakeHTTPHelper.PatchReturns(dummyResponse, nil)
		})

		It("performs PATCH requests to /uploads/:id?state=finished", func() {
			expectedBody := []byte{}
			fakeJSONHelper.MarshalReturns(expectedBody, nil)
			fakeJSONHelper.UnmarshalReturns(&wundergo.Upload{}, nil)
			expectedUrl := fmt.Sprintf("%s/uploads/%d?state=finished", apiURL, uploadID)

			client.MarkUploadComplete(uploadID)

			Expect(fakeHTTPHelper.PatchCallCount()).To(Equal(1))
			arg0, _ := fakeHTTPHelper.PatchArgsForCall(0)
			Expect(arg0).To(Equal(expectedUrl))
		})

		Context("when httpHelper.Patch returns an error", func() {
			expectedError := errors.New("httpHelper PATCH error")

			BeforeEach(func() {
				fakeHTTPHelper.PatchReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.MarkUploadComplete(uploadID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when response status code is unexpected", func() {
			BeforeEach(func() {
				dummyResponse.StatusCode = http.StatusBadRequest
			})

			It("returns an error", func() {
				_, err := client.MarkUploadComplete(uploadID)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body is nil", func() {
			BeforeEach(func() {
				dummyResponse.Body = nil
				fakeHTTPHelper.PatchReturns(dummyResponse, nil)
			})

			It("returns an error", func() {
				_, err := client.MarkUploadComplete(uploadID)

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
				_, err := client.MarkUploadComplete(uploadID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when unmarshalling json response returns an error", func() {
			expectedError := errors.New("jsonHelper error")

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				_, err := client.MarkUploadComplete(uploadID)

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when valid response is received", func() {
			expectedUpload := &wundergo.Upload{
				State: "finished",
			}

			BeforeEach(func() {
				fakeJSONHelper.UnmarshalReturns(expectedUpload, nil)
			})

			It("returns the unmarshalled note without error", func() {
				upload, err := client.MarkUploadComplete(uploadID)

				Expect(err).To(BeNil())
				Expect(upload).To(Equal(expectedUpload))
			})
		})
	})
})
