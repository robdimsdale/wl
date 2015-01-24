package wundergo_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("Basic membership functionality", func() {
	muted := true

	It("can add members via userID", func() {
		var lists []wundergo.List
		Eventually(func() error {
			l, err := client.Lists()
			lists = *l
			return err
		}).Should(Succeed())
		list := lists[0]

		var user wundergo.User
		Eventually(func() error {
			u, err := client.User()
			user = *u
			return err
		}).Should(Succeed())

		// Adding a member to a list they are already a member of
		// should return 201. This is odd, but it's the way it works

		Eventually(func() error {
			_, err := client.AddMemberToListViaUserID(user.ID, list.ID, muted)
			return err
		}).Should(Succeed())
	})

	It("can add members via emailAddress", func() {
		var lists []wundergo.List
		Eventually(func() error {
			l, err := client.Lists()
			lists = *l
			return err
		}).Should(Succeed())
		list := lists[0]

		var user wundergo.User
		Eventually(func() error {
			u, err := client.User()
			user = *u
			return err
		}).Should(Succeed())

		// Adding a member to a list they are already a member of
		// should return 201. This is odd, but it's the way it works

		Eventually(func() error {
			_, err := client.AddMemberToListViaEmailAddress(user.Email, list.ID, muted)
			return err
		}).Should(Succeed())
	})

})
