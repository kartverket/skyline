package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Context("Sender type is valid", func() {
		It("msgraph should return true", func() {
			Expect(SenderType.IsValid(MsGraph)).To(BeTrue())
		})
		It("dummy should return true", func() {
			Expect(SenderType.IsValid(Dummy)).To(BeTrue())
		})
		It("other should return false", func() {
			Expect(SenderType.IsValid(2)).ToNot(BeTrue())
		})
	})
})
