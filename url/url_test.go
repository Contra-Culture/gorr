package url_test

import (
	"net/url"

	. "github.com/Contra-Culture/gorr/url"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("url", func() {
	u, _ := url.Parse("https://github.com/Contra-Culture/gorr/commit/bf54f7ea25d833b811284220f064b487a0932ed4")
	Describe("Handle", func() {
		Context("when there is no conflict in params", func() {
			It("succeed and returns params map", func() {
				counter := 0
				params, err := Handle(
					"/",
					func(fragment string, mark func(string)) {
						counter++
					})
				Expect(err).NotTo(HaveOccurred())
				Expect(counter).To(Equal(0))
				Expect(params).To(Equal(map[string]string{
					"$path": "/",
				}))
				paramsOrder := []string{
					"orgName",
					"projectName",
					"",
					"commit",
				}
				idx := 0
				params, err = Handle(
					u.Path,
					func(fragment string, mark func(string)) {
						if len(paramsOrder[idx]) > 0 {
							mark(paramsOrder[idx])
						}
						idx++
					})
				Expect(idx).To(Equal(4))
				Expect(err).NotTo(HaveOccurred())
				Expect(params).To(
					Equal(map[string]string{
						"orgName":     "Contra-Culture",
						"projectName": "gorr",
						"commit":      "bf54f7ea25d833b811284220f064b487a0932ed4",
						"$path":       "/Contra-Culture/gorr/commit/bf54f7ea25d833b811284220f064b487a0932ed4",
					}))
			})
		})
		Context("when there is a conflict in params", func() {
			It("fails and returns error", func() {
				paramsOrder := []string{
					"orgName",
					"orgName",
					"",
					"commit",
				}
				idx := 0
				params, err := Handle(
					u.Path,
					func(fragment string, mark func(string)) {
						if len(paramsOrder[idx]) > 0 {
							mark(paramsOrder[idx])
						}
						idx++
					})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("parameter \"orgName\" already marked"))
				Expect(params).To(BeNil())
			})
		})
	})
})
