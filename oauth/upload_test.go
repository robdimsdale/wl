package oauth_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/robdimsdale/wl/oauth"
)

var _ = Describe("client - Upload operations", func() {
	Describe("creating a new upload", func() {
		var (
			localFilePath  string
			remoteFileName string
			contentType    string
			md5sum         string

			tempDirPath string
			fileSize    int
		)

		BeforeEach(func() {
			uuid1, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			remoteFileName = uuid1.String()

			tempDirPath, err = ioutil.TempDir(os.TempDir(), "wl-integration-test")
			Expect(err).NotTo(HaveOccurred())

			localFilePath = filepath.Join(tempDirPath, "test-file")

			fileContent := []byte("some-text")
			err = ioutil.WriteFile(localFilePath, fileContent, os.ModePerm)

			contentType = "text"
			md5sum = ""
			fileSize = 9 // size of "some-text" with no newline
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDirPath)
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("Initial upload creation", func() {
			It("performs POST requests with correct headers to /uploads", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/uploads"),
						ghttp.VerifyHeader(http.Header{
							"X-Access-Token": []string{dummyAccessToken},
							"X-Client-ID":    []string{dummyClientID},
						}),
						ghttp.VerifyJSON(fmt.Sprintf(
							`{"content_type":"%s","file_name":"%s","file_size":%d}`,
							contentType,
							remoteFileName,
							fileSize,
						)),
					),
				)

				client.UploadFile(
					localFilePath,
					remoteFileName,
					contentType,
					md5sum,
				)

				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})

			Context("when localFilePath is empty", func() {
				BeforeEach(func() {
					localFilePath = ""
				})

				It("returns an error", func() {
					_, err := client.UploadFile(
						localFilePath,
						remoteFileName,
						contentType,
						md5sum,
					)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when remoteFileName is empty", func() {
				BeforeEach(func() {
					remoteFileName = ""
				})

				It("returns an error", func() {
					_, err := client.UploadFile(
						localFilePath,
						remoteFileName,
						contentType,
						md5sum,
					)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when creating request fails with error", func() {
				BeforeEach(func() {
					client = oauth.NewClient("", "", "", testLogger)
				})

				It("forwards the error", func() {
					_, err := client.UploadFile(
						localFilePath,
						remoteFileName,
						contentType,
						md5sum,
					)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when executing request fails with error", func() {
				BeforeEach(func() {
					client = oauth.NewClient("", "", "http://not-a-real-url.com", testLogger)
				})

				It("forwards the error", func() {
					_, err := client.UploadFile(
						localFilePath,
						remoteFileName,
						contentType,
						md5sum,
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

					_, err := client.UploadFile(
						localFilePath,
						remoteFileName,
						contentType,
						md5sum,
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

					_, err := client.UploadFile(
						localFilePath,
						remoteFileName,
						contentType,
						md5sum,
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

					_, err := client.UploadFile(
						localFilePath,
						remoteFileName,
						contentType,
						md5sum,
					)

					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
