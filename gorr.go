package gorr

import (
	"net/http"

	"github.com/Contra-Culture/report"
)

type (
	Dispatcher struct {
		root *Node
	}
	DispatcherCfgr struct {
		dispatcher *Dispatcher
		report     report.Node
	}
	PathHelper  func(map[string]string) string
	PathHelpers map[string]PathHelper
)

func (c *DispatcherCfgr) Root(d string, cfg func(*StaticNodeCfgr)) {
	if c.dispatcher.root != nil {
		c.report.Error("root node already specified")
		return
	}
	r := report.New("root")
	root := new(nil, STATIC)
	nc := &StaticNodeCfgr{
		NodeCfgr{
			node:   root,
			report: r,
		},
	}
	cfg(nc)
	c.dispatcher.root = root
	c.report = r
}
func New(cfg func(*DispatcherCfgr)) (d *Dispatcher, r report.Node) {
	d = &Dispatcher{}
	r = report.New("dispatcher")
	cfg(
		&DispatcherCfgr{
			dispatcher: d,
			report:     r,
		})
	return
}
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.root.Handle(w, r)
}
