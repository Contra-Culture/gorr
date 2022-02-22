package gorr_test

import (
	"net/http"
	"strings"
	"time"

	. "github.com/Contra-Culture/gorr"
	"github.com/Contra-Culture/report"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func dumbHandler(w http.ResponseWriter, r *http.Request, ps map[string]string) {}

var _ = Describe("gorr", func() {
	Describe("router DSL", func() {
		var (
			d *Dispatcher
			r report.Node
		)
		Context("when configured properly", func() {
			It("succeed", func() {
				d, r = New(
					func(cfg *DispatcherCfgr) {
						cfg.Root(
							"Test root.",
							func(cfg *StaticNodeCfgr) {
							})
					})
				Expect(d).NotTo(BeNil())
				Expect(r).NotTo(BeNil())
				var sb strings.Builder
				fn := func(path []int, k report.Kind, _ time.Time, _ *time.Duration, s string) (err error) {
					for range path {
						_, err = sb.WriteRune('\t')
						if err != nil {
							return
						}
					}
					_, err = sb.WriteString(s)
					if err != nil {
						return
					}
					_, err = sb.WriteString("\n\n")
					return
				}
				Expect(r.Traverse(fn)).NotTo(HaveOccurred())
				Expect(sb.String()).To(Equal("dispatcher\n\n"))
			})
		})
	})
})
