package wl_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wl"
)

var _ = Describe("basic user functionality", func() {
	It("can update the user's name", func() {
		By("Creating a new random user name")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newUserName := uuid1.String()

		By("Getting user")
		var user wl.User
		Eventually(func() error {
			user, err = client.User()
			return err
		}).Should(Succeed())

		By("Updating user")
		var updatedUser wl.User
		user.Name = "test-" + newUserName
		Eventually(func() error {
			updatedUser, err = client.UpdateUser(user)
			return err
		}).Should(Succeed())

		Expect(updatedUser.ID).To(Equal(user.ID))
	})
})
