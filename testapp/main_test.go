package main_test

import (
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("web app", func() {
	var (
		p   *os.Process
		err error
	)
	BeforeSuite(func() {
		p, err = os.StartProcess("go", []string{"run ./main.go"}, &os.ProcAttr{})
		Expect(err).NotTo(HaveOccurred())
	})
	AfterSuite(func() {
		p.Kill()
	})
	It("runs successfully", func() {
		resp, err := http.Get("http://localhost:8080/articles/")
		defer (func() {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		})()
		Expect(err).NotTo(HaveOccurred())
		bytes := make([]byte, 1000)
		_, err = resp.Body.Read(bytes)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(bytes)).To(BeEquivalentTo(""))
	})

	It("works fine", func() {

	})
})
