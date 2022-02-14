package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var err error
var mainPath string

var _ = Describe("web app", func() {
	It("works", func() {
		buildcmd := exec.Command("go", "build", "main.go")
		err = buildcmd.Run()
		Expect(err).NotTo(HaveOccurred())
		runcmd := exec.Command("./main")
		session, err := gexec.Start(runcmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gbytes.Say(``))
	})
})
