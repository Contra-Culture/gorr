package gorr_test

import (
	"net/http"

	. "github.com/Contra-Culture/gorr"
	"github.com/Contra-Culture/gorr/node"
	"github.com/Contra-Culture/report"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dumbHandler(w http.ResponseWriter, r *http.Request, ps map[string]string) {}

var _ = Describe("gorr", func() {
	Describe("router DSL", func() {
		var (
			d *Dispatcher
			r *report.RContext
		)
		Context("when configured properly", func() {
			It("succeed", func() {
				d, r = New(
					func(cfg *DispatcherCfgr) {
						cfg.Root(
							"Test root.",
							func(cfg *node.NodeCfgr) {
							})
					})
				Expect(d).NotTo(BeNil())
				Expect(r).NotTo(BeNil())
				Expect(r.String()).To(Equal("root: dispatcher\n"))
			})
		})
	})
})
