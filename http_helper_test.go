package wundergo_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/fakes"
)

var _ = Describe("Client", func() {
	var fakeHTTPTransport fakes.FakeHTTPTransport

	var httpHelper wundergo.HTTPHelper

	var dummyRequest *http.Request

	BeforeEach(func() {
		fakeHTTPTransport = fakes.FakeHTTPTransport{}

		wundergo.NewHTTPTransport = func() wundergo.HTTPTransport {
			return &fakeHTTPTransport
		}

		dummyRequest = &http.Request{
			Header: make(http.Header),
		}

		httpHelper = wundergo.NewOauthClientHTTPHelper(dummyAccessToken, dummyClientID)
	})

	verifyAuthHeaders := func() {
		accessTokenHeader := dummyRequest.Header.Get("X-Access-Token")
		clientIDHeader := dummyRequest.Header.Get("X-Client-ID")

		Expect(accessTokenHeader).To(Equal(dummyAccessToken))
		Expect(clientIDHeader).To(Equal(dummyClientID))
	}

	Describe("GET requests", func() {
		Context("when httpTrasport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransport error")

			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(nil, expectedError)
			})

			It("returns nil byte array", func() {
				b, _ := httpHelper.Get("someUrl")

				Expect(b).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Get("someUrl")

				Expect(err).To(Equal(expectedError))
			})

			Context("when request creation is successful", func() {
				BeforeEach(func() {
					fakeHTTPTransport.NewRequestReturns(dummyRequest, nil)
				})

				It("adds authentication authorization headers to request", func() {
					httpHelper.Get("someUrl")

					verifyAuthHeaders()
				})

				Context("when httpTrasport.DoRequest returns with error", func() {
					expectedError := errors.New("fakeHTTPTransport error")

					BeforeEach(func() {
						fakeHTTPTransport.DoRequestReturns(nil, expectedError)
					})

					It("returns nil byte array", func() {
						b, _ := httpHelper.Get("someUrl")

						Expect(b).To(BeNil())
					})

					It("forwards the error", func() {
						_, err := httpHelper.Get("someUrl")

						Expect(err).To(Equal(expectedError))
					})
				})

				Context("when httpTrasport.DoRequest returns with nil response", func() {
					BeforeEach(func() {
						fakeHTTPTransport.DoRequestReturns(nil, nil)
					})

					It("returns nil byte array", func() {
						b, _ := httpHelper.Get("someUrl")

						Expect(b).To(BeNil())
					})

					It("returns an error", func() {
						_, err := httpHelper.Get("someUrl")

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when httpTrasport.DoRequest returns with valid response", func() {
					var validResponse *http.Response

					BeforeEach(func() {
						validResponse = &http.Response{}
						fakeHTTPTransport.DoRequestReturns(validResponse, nil)
					})

					Context("when the response body is nil", func() {
						It("returns empty byte array without error", func() {
							b, err := httpHelper.Get("someUrl")

							Expect(err).To(BeNil())
							Expect(b).To(Equal([]byte{}))
						})
					})

					Context("when the response body is non nil", func() {
						expectedResponseBody := []byte("expectedResponseBody")

						BeforeEach(func() {
							validResponse.Body = ioutil.NopCloser(bytes.NewReader(expectedResponseBody))
						})

						It("returns the body wihout error", func() {
							b, err := httpHelper.Get("someUrl")

							Expect(err).To(BeNil())
							Expect(b).To(Equal(expectedResponseBody))
						})
					})
				})
			})
		})
	})

	Describe("POST requests", func() {
		Context("when httpTrasport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransport error")

			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(nil, expectedError)
			})

			It("returns nil byte array", func() {
				b, _ := httpHelper.Post("someUrl", []byte("someRequestBody"))

				Expect(b).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when request creation is successful", func() {
			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(dummyRequest, nil)
			})

			It("adds authentication authorization headers to request", func() {
				httpHelper.Post("someUrl", []byte("someRequestBody"))

				verifyAuthHeaders()
			})

			It("adds 'Content-Type: application/json' header", func() {
				httpHelper.Post("someUrl", []byte("someRequestBody"))
				contentTypeHeader := dummyRequest.Header.Get("Content-Type")

				Expect(contentTypeHeader).To(Equal("application/json"))
			})

			Context("when body is not empty", func() {
				It("adds body to request", func() {
					httpHelper.Post("someUrl", []byte("someRequestBody"))

					body := dummyRequest.Body
					bodyContent, err := ioutil.ReadAll(body)
					Expect(err).To(BeNil())
					Expect(bodyContent).To(Equal([]byte([]byte("someRequestBody"))))
				})
			})

			Context("when body is nil", func() {
				It("does not add body to request", func() {
					httpHelper.Post("someUrl", nil)

					body := dummyRequest.Body
					Expect(body).To(BeNil())
				})
			})

			Context("when body is empty", func() {
				It("does not add body to request", func() {
					httpHelper.Post("someUrl", []byte(""))

					body := dummyRequest.Body
					Expect(body).To(BeNil())
				})
			})

			Context("when httpTrasport.DoRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransport error")

				BeforeEach(func() {
					fakeHTTPTransport.DoRequestReturns(nil, expectedError)
				})

				It("returns nil byte array", func() {
					b, _ := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(b).To(BeNil())
				})

				It("forwards the error", func() {
					_, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when httpTrasport.DoRequest returns with nil response", func() {
				BeforeEach(func() {
					fakeHTTPTransport.DoRequestReturns(nil, nil)
				})

				It("returns nil byte array", func() {
					b, _ := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(b).To(BeNil())
				})

				It("returns an error", func() {
					_, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when httpTrasport.DoRequest returns with valid response", func() {
				var validResponse *http.Response

				BeforeEach(func() {
					validResponse = &http.Response{}
					fakeHTTPTransport.DoRequestReturns(validResponse, nil)
				})

				Context("when the response body is nil", func() {
					It("returns empty byte array without error", func() {
						b, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

						Expect(err).To(BeNil())
						Expect(b).To(Equal([]byte{}))
					})
				})

				Context("when the response body is non nil", func() {
					expectedResponseBody := []byte("expectedResponseBody")

					BeforeEach(func() {
						validResponse.Body = ioutil.NopCloser(bytes.NewReader(expectedResponseBody))
					})

					It("returns the body wihout error", func() {
						b, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

						Expect(err).To(BeNil())
						Expect(b).To(Equal(expectedResponseBody))
					})
				})
			})
		})
	})

	Describe("PUT requests", func() {
		Context("when httpTrasport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransport error")

			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(nil, expectedError)
			})

			It("returns nil byte array", func() {
				b, _ := httpHelper.Put("someUrl", []byte("someRequestBody"))

				Expect(b).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Put("someUrl", []byte("someRequestBody"))

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when request creation is successful", func() {
			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(dummyRequest, nil)
			})

			It("adds authentication authorization headers to request", func() {
				httpHelper.Put("someUrl", []byte("someRequestBody"))

				verifyAuthHeaders()
			})

			It("adds 'Content-Type: application/x-www-form-urlencoded' header", func() {
				httpHelper.Put("someUrl", []byte("someRequestBody"))
				contentTypeHeader := dummyRequest.Header.Get("Content-Type")

				Expect(contentTypeHeader).To(Equal("application/x-www-form-urlencoded"))
			})

			Context("when body is not empty", func() {
				It("adds body to request", func() {
					httpHelper.Put("someUrl", []byte("someRequestBody"))

					body := dummyRequest.Body
					bodyContent, err := ioutil.ReadAll(body)
					Expect(err).To(BeNil())
					Expect(bodyContent).To(Equal([]byte([]byte("someRequestBody"))))
				})
			})

			Context("when body is nil", func() {
				It("does not add body to request", func() {
					httpHelper.Put("someUrl", nil)

					body := dummyRequest.Body
					Expect(body).To(BeNil())
				})
			})

			Context("when body is empty", func() {
				It("does not add body to request", func() {
					httpHelper.Put("someUrl", []byte(""))

					body := dummyRequest.Body
					Expect(body).To(BeNil())
				})
			})

			Context("when httpTrasport.DoRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransport error")

				BeforeEach(func() {
					fakeHTTPTransport.DoRequestReturns(nil, expectedError)
				})

				It("returns nil byte array", func() {
					b, _ := httpHelper.Put("someUrl", []byte("someRequestBody"))

					Expect(b).To(BeNil())
				})

				It("forwards the error", func() {
					_, err := httpHelper.Put("someUrl", []byte("someRequestBody"))

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when httpTrasport.DoRequest returns with nil response", func() {
				BeforeEach(func() {
					fakeHTTPTransport.DoRequestReturns(nil, nil)
				})

				It("returns nil byte array", func() {
					b, _ := httpHelper.Put("someUrl", []byte("someRequestBody"))

					Expect(b).To(BeNil())
				})

				It("returns an error", func() {
					_, err := httpHelper.Put("someUrl", []byte("someRequestBody"))

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when httpTrasport.DoRequest returns with valid response", func() {
				var validResponse *http.Response

				BeforeEach(func() {
					validResponse = &http.Response{}
					fakeHTTPTransport.DoRequestReturns(validResponse, nil)
				})

				Context("when the response body is nil", func() {
					It("returns empty byte array without error", func() {
						b, err := httpHelper.Put("someUrl", []byte("someRequestBody"))

						Expect(err).To(BeNil())
						Expect(b).To(Equal([]byte{}))
					})
				})

				Context("when the response body is non nil", func() {
					expectedResponseBody := []byte("expectedResponseBody")

					BeforeEach(func() {
						validResponse.Body = ioutil.NopCloser(bytes.NewReader(expectedResponseBody))
					})

					It("returns the body wihout error", func() {
						b, err := httpHelper.Put("someUrl", []byte("someRequestBody"))

						Expect(err).To(BeNil())
						Expect(b).To(Equal(expectedResponseBody))
					})
				})
			})
		})
	})

	Describe("PATCH requests", func() {
		Context("when httpTrasport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransport error")

			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(nil, expectedError)
			})

			It("returns nil byte array", func() {
				b, _ := httpHelper.Patch("someUrl", []byte("someRequestBody"))

				Expect(b).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Patch("someUrl", []byte("someRequestBody"))

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when request creation is successful", func() {
			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(dummyRequest, nil)
			})

			It("adds authentication authorization headers to request", func() {
				httpHelper.Patch("someUrl", []byte("someRequestBody"))

				verifyAuthHeaders()
			})

			It("adds 'Content-Type: application/json' header", func() {
				httpHelper.Patch("someUrl", []byte("someRequestBody"))
				contentTypeHeader := dummyRequest.Header.Get("Content-Type")

				Expect(contentTypeHeader).To(Equal("application/json"))
			})

			Context("when body is not empty", func() {
				It("adds body to request", func() {
					httpHelper.Patch("someUrl", []byte("someRequestBody"))

					body := dummyRequest.Body
					bodyContent, err := ioutil.ReadAll(body)
					Expect(err).To(BeNil())
					Expect(bodyContent).To(Equal([]byte([]byte("someRequestBody"))))
				})
			})

			Context("when body is nil", func() {
				It("does not add body to request", func() {
					httpHelper.Patch("someUrl", nil)

					body := dummyRequest.Body
					Expect(body).To(BeNil())
				})
			})

			Context("when body is empty", func() {
				It("does not add body to request", func() {
					httpHelper.Patch("someUrl", []byte(""))

					body := dummyRequest.Body
					Expect(body).To(BeNil())
				})
			})

			Context("when httpTrasport.DoRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransport error")

				BeforeEach(func() {
					fakeHTTPTransport.DoRequestReturns(nil, expectedError)
				})

				It("returns nil byte array", func() {
					b, _ := httpHelper.Patch("someUrl", []byte("someRequestBody"))

					Expect(b).To(BeNil())
				})

				It("forwards the error", func() {
					_, err := httpHelper.Patch("someUrl", []byte("someRequestBody"))

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when httpTrasport.DoRequest returns with nil response", func() {
				BeforeEach(func() {
					fakeHTTPTransport.DoRequestReturns(nil, nil)
				})

				It("returns nil byte array", func() {
					b, _ := httpHelper.Patch("someUrl", []byte("someRequestBody"))

					Expect(b).To(BeNil())
				})

				It("returns an error", func() {
					_, err := httpHelper.Patch("someUrl", []byte("someRequestBody"))

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when httpTrasport.DoRequest returns with valid response", func() {
				var validResponse *http.Response

				BeforeEach(func() {
					validResponse = &http.Response{}
					fakeHTTPTransport.DoRequestReturns(validResponse, nil)
				})

				Context("when the response body is nil", func() {
					It("returns empty byte array without error", func() {
						b, err := httpHelper.Patch("someUrl", []byte("someRequestBody"))

						Expect(err).To(BeNil())
						Expect(b).To(Equal([]byte{}))
					})
				})

				Context("when the response body is non nil", func() {
					expectedResponseBody := []byte("expectedResponseBody")

					BeforeEach(func() {
						validResponse.Body = ioutil.NopCloser(bytes.NewReader(expectedResponseBody))
					})

					It("returns the body wihout error", func() {
						b, err := httpHelper.Patch("someUrl", []byte("someRequestBody"))

						Expect(err).To(BeNil())
						Expect(b).To(Equal(expectedResponseBody))
					})
				})
			})
		})
	})

	Describe("DELETE requests", func() {
		Context("when httpTrasport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransport error")

			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(nil, expectedError)
			})

			It("forwards the error", func() {
				err := httpHelper.Delete("someUrl")

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when request creation is successful", func() {
			BeforeEach(func() {
				fakeHTTPTransport.NewRequestReturns(dummyRequest, nil)
			})

			It("adds authentication authorization headers to request", func() {
				httpHelper.Delete("someUrl")

				verifyAuthHeaders()
			})

			Context("when httpTrasport.DoRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransport error")

				BeforeEach(func() {
					fakeHTTPTransport.DoRequestReturns(nil, expectedError)
				})

				It("forwards the error", func() {
					err := httpHelper.Delete("someUrl")

					Expect(err).To(Equal(expectedError))
				})
			})
		})
	})
})
