package wl_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

const (
	processorType = "generic"
	configuration = ""
)

var _ = Describe("basic webhook functionality", func() {
	var (
		newList wl.List
	)

	BeforeEach(func() {
		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		Eventually(func() error {
			newList, err = client.CreateList(newListTitle)
			return err
		}).Should(Succeed())
	})

	AfterEach(func() {
		By("Deleting new list")
		Eventually(func() error {
			l, err := client.List(newList.ID)
			if err != nil {
				return err
			}
			return client.DeleteList(l)
		}).Should(Succeed())

		Eventually(func() (bool, error) {
			lists, err := client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})

	It("can list, create and delete webhooks", func() {
		var err error

		By("Listing existing webhooks")
		var webhooks []wl.Webhook
		Eventually(func() error {
			webhooks, err = client.WebhooksForListID(newList.ID)
			return err
		}).Should(Succeed())
		Expect(len(webhooks)).To(BeZero())

		By("Creating a new webhook")
		url := "https://some-fake-url.com"

		var newWebhook wl.Webhook
		Eventually(func() error {
			newWebhook, err = client.CreateWebhook(
				newList.ID,
				url,
				processorType,
				configuration,
			)
			return err
		}).Should(Succeed())
		Expect(newWebhook.URL).To(Equal(url))

		By("Validating the new webhook is present in all webhooks")
		Eventually(func() bool {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			allWebhooks, _ := client.Webhooks()
			return webhooksContain(allWebhooks, newWebhook)
		}).Should(BeTrue())

		By("Validating the new webhook is present in webhooks for list")
		Eventually(func() (bool, error) {
			webhooks, err := client.WebhooksForListID(newList.ID)
			return webhooksContain(webhooks, newWebhook), err
		}).Should(BeTrue())

		By("Validating the new webhook can be retrieved")
		var aWebhook wl.Webhook
		Eventually(func() wl.Webhook {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			aWebhook, _ = client.Webhook(newWebhook.ID)
			return aWebhook
		}).Should(Equal(newWebhook))

		By("Deleting the new webhook")
		Eventually(func() error {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			w, _ := client.Webhook(newWebhook.ID)
			return client.DeleteWebhook(w)
		}).Should(Succeed())

		By("Validating the new webhook is not present in list")
		Eventually(func() bool {
			// It is statistically probable that one of the lists will
			// be deleted, so we ignore error here.
			webhooks, _ := client.WebhooksForListID(newList.ID)
			return webhooksContain(webhooks, newWebhook)
		}).Should(BeFalse())
	})
})

func webhooksContain(webhooks []wl.Webhook, webhook wl.Webhook) bool {
	for _, w := range webhooks {
		if w.ID == webhook.ID {
			return true
		}
	}
	return false
}
