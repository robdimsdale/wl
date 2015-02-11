package wundergo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

type testStruct struct {
	Name    string
	Counter uint
}

var _ = Describe("default JSON Helper", func() {
	var jsonHelper *wundergo.DefaultJSONHelper

	BeforeEach(func() {
		jsonHelper = &wundergo.DefaultJSONHelper{}
	})

	Describe("marshal", func() {
		It("marshals into json without error", func() {
			expectedReturn := []byte(`{"Name":"myName","Counter":2}`)
			ts := testStruct{
				Name:    "myName",
				Counter: 2,
			}

			b, err := jsonHelper.Marshal(ts)

			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(Equal(expectedReturn))
		})
	})

	Describe("unmarshal", func() {
		It("unmarshals from json without error", func() {
			jsonInput := []byte(`{"Name":"myName","Counter":2}`)
			expectedReturn := testStruct{
				Name:    "myName",
				Counter: 2,
			}

			ts := testStruct{}

			returnedInterface, err := jsonHelper.Unmarshal(jsonInput, &ts)

			returned := returnedInterface.(*testStruct)

			Expect(err).NotTo(HaveOccurred())
			Expect(*returned).To(Equal(expectedReturn))
			Expect(ts).To(Equal(expectedReturn))
		})
	})
})
