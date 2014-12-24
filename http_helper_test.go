package wundergo_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/fakes"
)

var _ = Describe("HTTPHelper", func() {
	var fakeHTTPTransportHelper fakes.FakeHTTPTransportHelper

	var httpHelper wundergo.HTTPHelper

	var dummyRequest *http.Request

	BeforeEach(func() {
		fakeHTTPTransportHelper = fakes.FakeHTTPTransportHelper{}

		wundergo.NewHTTPTransportHelper = func() wundergo.HTTPTransportHelper {
			return &fakeHTTPTransportHelper
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
		Context("when httpTransport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransportHelper error")

			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(nil, expectedError)
			})

			It("returns nil response", func() {
				resp, _ := httpHelper.Get("someUrl")

				Expect(resp).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Get("someUrl")

				Expect(err).To(Equal(expectedError))
			})

			Context("when request creation is successful", func() {
				BeforeEach(func() {
					fakeHTTPTransportHelper.NewRequestReturns(dummyRequest, nil)
				})

				It("adds authentication authorization headers to request", func() {
					httpHelper.Get("someUrl")

					verifyAuthHeaders()
				})

				Context("when httpTransport.DoRequest returns with error", func() {
					expectedError := errors.New("fakeHTTPTransportHelper error")

					BeforeEach(func() {
						fakeHTTPTransportHelper.DoRequestReturns(nil, expectedError)
					})

					It("returns nil response", func() {
						resp, _ := httpHelper.Get("someUrl")

						Expect(resp).To(BeNil())
					})

					It("forwards the error", func() {
						_, err := httpHelper.Get("someUrl")

						Expect(err).To(Equal(expectedError))
					})
				})

				Context("when httpTransport.DoRequest returns with nil response", func() {
					BeforeEach(func() {
						fakeHTTPTransportHelper.DoRequestReturns(nil, nil)
					})

					It("returns nil response", func() {
						resp, _ := httpHelper.Get("someUrl")

						Expect(resp).To(BeNil())
					})

					It("returns an error", func() {
						_, err := httpHelper.Get("someUrl")

						Expect(err).ToNot(BeNil())
					})
				})

				Context("when httpTransport.DoRequest returns with valid response", func() {
					var validResponse *http.Response

					BeforeEach(func() {
						validResponse = &http.Response{}
						fakeHTTPTransportHelper.DoRequestReturns(validResponse, nil)
					})

					It("returns the response wihout error", func() {
						resp, err := httpHelper.Get("someUrl")

						Expect(err).To(BeNil())
						Expect(resp).To(Equal(validResponse))
					})
				})
			})
		})
	})

	Describe("POST requests", func() {
		Context("when httpTransport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransportHelper error")

			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(nil, expectedError)
			})

			It("returns nil response", func() {
				resp, _ := httpHelper.Post("someUrl", []byte("someRequestBody"))

				Expect(resp).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when request creation is successful", func() {
			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(dummyRequest, nil)
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

			Context("when httpTransport.DoRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransportHelper error")

				BeforeEach(func() {
					fakeHTTPTransportHelper.DoRequestReturns(nil, expectedError)
				})

				It("returns nil response", func() {
					resp, _ := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(resp).To(BeNil())
				})

				It("forwards the error", func() {
					_, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when httpTransport.DoRequest returns with nil response", func() {
				BeforeEach(func() {
					fakeHTTPTransportHelper.DoRequestReturns(nil, nil)
				})

				It("returns nil response", func() {
					resp, _ := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(resp).To(BeNil())
				})

				It("returns an error", func() {
					_, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when httpTransport.DoRequest returns with valid response", func() {
				var validResponse *http.Response

				BeforeEach(func() {
					validResponse = &http.Response{}
					fakeHTTPTransportHelper.DoRequestReturns(validResponse, nil)
				})

				It("returns the response wihout error", func() {
					resp, err := httpHelper.Post("someUrl", []byte("someRequestBody"))

					Expect(err).To(BeNil())
					Expect(resp).To(Equal(validResponse))
				})
			})
		})
	})

	Describe("PUT requests", func() {
		Context("when httpTransport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransportHelper error")

			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(nil, expectedError)
			})

			It("returns nil response", func() {
				resp, _ := httpHelper.Put("someUrl", []byte("someRequestBody"))

				Expect(resp).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Put("someUrl", []byte("someRequestBody"))

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when request creation is successful", func() {
			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(dummyRequest, nil)
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

			Context("when httpTransport.DoRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransportHelper error")

				BeforeEach(func() {
					fakeHTTPTransportHelper.DoRequestReturns(nil, expectedError)
				})

				It("returns nil response", func() {
					resp, _ := httpHelper.Put("someUrl", []byte("someRequestBody"))

					Expect(resp).To(BeNil())
				})

				It("forwards the error", func() {
					_, err := httpHelper.Put("someUrl", []byte("someRequestBody"))

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when httpTransport.DoRequest returns with nil response", func() {
				BeforeEach(func() {
					fakeHTTPTransportHelper.DoRequestReturns(nil, nil)
				})

				It("returns nil response", func() {
					resp, _ := httpHelper.Put("someUrl", []byte("someRequestBody"))

					Expect(resp).To(BeNil())
				})

				It("returns an error", func() {
					_, err := httpHelper.Put("someUrl", []byte("someRequestBody"))

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when httpTransport.DoRequest returns with valid response", func() {
				var validResponse *http.Response

				BeforeEach(func() {
					validResponse = &http.Response{}
					fakeHTTPTransportHelper.DoRequestReturns(validResponse, nil)
				})

				It("returns the response wihout error", func() {
					resp, err := httpHelper.Put("someUrl", []byte("someResponseBody"))

					Expect(err).To(BeNil())
					Expect(resp).To(Equal(validResponse))
				})
			})
		})
	})

	Describe("PATCH requests", func() {
		Context("when httpTransport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransportHelper error")

			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(nil, expectedError)
			})

			It("returns nil response", func() {
				resp, _ := httpHelper.Patch("someUrl", []byte("someRequestBody"))

				Expect(resp).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Patch("someUrl", []byte("someRequestBody"))

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when request creation is successful", func() {
			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(dummyRequest, nil)
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

			Context("when httpTransport.DoRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransportHelper error")

				BeforeEach(func() {
					fakeHTTPTransportHelper.DoRequestReturns(nil, expectedError)
				})

				It("returns nil response", func() {
					resp, _ := httpHelper.Patch("someUrl", []byte("someRequestBody"))

					Expect(resp).To(BeNil())
				})

				It("forwards the error", func() {
					_, err := httpHelper.Patch("someUrl", []byte("someRequestBody"))

					Expect(err).To(Equal(expectedError))
				})
			})

			Context("when httpTransport.DoRequest returns with nil response", func() {
				BeforeEach(func() {
					fakeHTTPTransportHelper.DoRequestReturns(nil, nil)
				})

				It("returns nil response", func() {
					resp, _ := httpHelper.Patch("someUrl", []byte("someRequestBody"))

					Expect(resp).To(BeNil())
				})

				It("returns an error", func() {
					_, err := httpHelper.Patch("someUrl", []byte("someRequestBody"))

					Expect(err).ToNot(BeNil())
				})
			})

			Context("when httpTransport.DoRequest returns with valid response", func() {
				var validResponse *http.Response

				BeforeEach(func() {
					validResponse = &http.Response{}
					fakeHTTPTransportHelper.DoRequestReturns(validResponse, nil)
				})

				It("returns the response wihout error", func() {
					resp, err := httpHelper.Patch("someUrl", []byte("someResponseBody"))

					Expect(err).To(BeNil())
					Expect(resp).To(Equal(validResponse))
				})
			})
		})
	})

	Describe("DELETE requests", func() {
		Context("when httpTransport.NewRequest returns with error", func() {
			expectedError := errors.New("fakeHTTPTransportHelper error")

			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(nil, expectedError)
			})

			It("returns nil response", func() {
				resp, _ := httpHelper.Delete("someUrl")

				Expect(resp).To(BeNil())
			})

			It("forwards the error", func() {
				_, err := httpHelper.Delete("someUrl")

				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when request creation is successful", func() {
			BeforeEach(func() {
				fakeHTTPTransportHelper.NewRequestReturns(dummyRequest, nil)
			})

			It("adds authentication authorization headers to request", func() {
				httpHelper.Delete("someUrl")

				verifyAuthHeaders()
			})

			Context("when httpTransport.DoRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransportHelper error")

				BeforeEach(func() {
					fakeHTTPTransportHelper.DoRequestReturns(nil, expectedError)
				})

				It("returns nil response", func() {
					resp, _ := httpHelper.Delete("someUrl")

					Expect(resp).To(BeNil())
				})

				It("forwards the error", func() {
					_, err := httpHelper.Delete("someUrl")

					Expect(err).To(Equal(expectedError))
				})
			})
		})
	})
})
