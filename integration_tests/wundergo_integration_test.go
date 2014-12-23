package wundergo_test

import (
	"log"
	"os"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var (
	client wundergo.Client
)

func contains(lists []wundergo.List, list wundergo.List) bool {
	for _, l := range lists {
		if l == list {
			return true
		}
	}
	return false
}

var _ = Describe("Wundergo library", func() {
	BeforeEach(func() {
		accessToken := os.Getenv("WL_ACCESS_TOKEN")
		clientID := os.Getenv("WL_CLIENT_ID")

		if accessToken == "" {
			log.Fatal("Error - WL_ACCESS_TOKEN must be provided")
		}

		if clientID == "" {
			log.Fatal("Error - WL_CLIENT_ID must be provided")
		}

		client = wundergo.NewOauthClient(accessToken, clientID)
	})

	Describe("Basic functionality", func() {
		It("creates, updates and deletes new list", func() {
			uuid1, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			newListTitle1 := uuid1.String()

			uuid2, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			newListTitle2 := uuid2.String()

			originalLists, err := client.Lists()
			Expect(err).NotTo(HaveOccurred())

			newList, err := client.CreateList(newListTitle1)
			Expect(err).NotTo(HaveOccurred())

			newLists, err := client.Lists()
			Expect(err).NotTo(HaveOccurred())
			Expect(contains(newLists, newList)).To(BeTrue())

			newList.Title = newListTitle2
			updatedList, err := client.UpdateList(newList)
			Expect(err).NotTo(HaveOccurred())
			newList.Revision = newList.Revision + 1
			Expect(updatedList).To(Equal(newList))

			err = client.DeleteList(newList)
			Expect(err).NotTo(HaveOccurred())

			afterDeleteLists, err := client.Lists()
			Expect(err).NotTo(HaveOccurred())

			Expect(afterDeleteLists).To(Equal(originalLists))
		})
	})
})
