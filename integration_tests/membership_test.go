package wl_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic membership functionality", func() {
	const muted = true

	var (
		inbox wl.List
		user  wl.User
	)

	BeforeEach(func() {
		var err error

		By("Getting inbox")
		Eventually(func() error {
			inbox, err = client.Inbox()
			return err
		}).Should(Succeed())

		By("Getting user")
		Eventually(func() error {
			user, err = client.User()
			return err
		}).Should(Succeed())
	})

	It("can add members via userID", func() {
		Eventually(func() error {
			_, err := client.AddMemberToListViaUserID(user.ID, inbox.ID, muted)
			return err
		}).Should(Succeed())
	})

	It("can add members via emailAddress", func() {
		Eventually(func() error {
			_, err := client.AddMemberToListViaEmailAddress(user.Email, inbox.ID, muted)
			return err
		}).Should(Succeed())
	})
})
