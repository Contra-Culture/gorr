package gorr_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGorr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gorr Suite")
}
