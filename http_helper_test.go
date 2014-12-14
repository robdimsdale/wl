package wundergo_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
	"github.com/robdimsdale/wundergo/fakes"
)

var _ = Describe("Client", func() {
	var fakeHTTPTransport fakes.FakeHTTPTransport

	var httpHelper wundergo.HTTPHelper

	BeforeEach(func() {
		fakeHTTPTransport = fakes.FakeHTTPTransport{}

		wundergo.NewHTTPTransport = func() wundergo.HTTPTransport {
			return &fakeHTTPTransport
		}

		httpHelper = wundergo.NewOauthClientHTTPHelper(dummyAccessToken, dummyClientID)
	})

	Describe("GET requests", func() {
		Context("When httpTrasport.NewRequest returns with error", func() {
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
		})

		Describe("POST requests", func() {
			Context("When httpTrasport.NewRequest returns with error", func() {
				expectedError := errors.New("fakeHTTPTransport error")

				BeforeEach(func() {
					fakeHTTPTransport.NewRequestReturns(nil, expectedError)
				})

				It("returns nil byte array", func() {
					b, _ := httpHelper.Post("someUrl", "someBody")

					Expect(b).To(BeNil())
				})

				It("forwards the error", func() {
					_, err := httpHelper.Post("someUrl", "someBody")

					Expect(err).To(Equal(expectedError))
				})
			})

			Describe("PUT requests", func() {
				Context("When httpTrasport.NewRequest returns with error", func() {
					expectedError := errors.New("fakeHTTPTransport error")

					BeforeEach(func() {
						fakeHTTPTransport.NewRequestReturns(nil, expectedError)
					})

					It("returns nil byte array", func() {
						b, _ := httpHelper.Put("someUrl", "someBody")

						Expect(b).To(BeNil())
					})

					It("forwards the error", func() {
						_, err := httpHelper.Put("someUrl", "someBody")

						Expect(err).To(Equal(expectedError))
					})
				})
			})
		})
	})
})
