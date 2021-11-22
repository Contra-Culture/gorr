package gorr

import (
	"net/http"

	"github.com/Contra-Culture/gorr/node"
	"github.com/Contra-Culture/report"
)

type (
	Dispatcher struct {
		root *node.Node
	}
	DispatcherCfgr struct {
		dispatcher *Dispatcher
		report     *report.RContext
	}
	PathHelper  func(map[string]string) string
	PathHelpers map[string]PathHelper
)

func (c *DispatcherCfgr) Root(d string, cfg func(*node.StaticNodeCfgr)) {
	if c.dispatcher.root != nil {
		c.report.Error("root node already specified")
		return
	}
	root, report := node.New(cfg)
	c.dispatcher.root = root
	c.report = report
}
func New(cfg func(*DispatcherCfgr)) (d *Dispatcher, r *report.RContext) {
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
