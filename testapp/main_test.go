package main_test

import (
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("web app", func() {
	var (
		path    string
		command *exec.Cmd
		session *gexec.Session
		err     error
	)
	BeforeSuite(func() {
		path, err = gexec.BuildWithEnvironment("./main.go", []string{})
		Expect(err).NotTo(HaveOccurred())
		command = exec.Command(path)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})
	AfterSuite(func() {
		session.Interrupt()
		gexec.CleanupBuildArtifacts()
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
