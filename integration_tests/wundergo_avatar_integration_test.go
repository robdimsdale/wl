package wundergo_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("avatar functionality", func() {
	It("can request a user's avatar", func() {
		userID := uint(1)
		size := 0
		fallback := false

		expectedURL := "https://avatars.wunderlist.io/uploads/user/avatar/0058/0058_64_4F7DDE.png"

		url, err := client.AvatarURL(userID, size, fallback)

		Expect(err).NotTo(HaveOccurred())
		Expect(url).To(Equal(expectedURL))
	})
})
