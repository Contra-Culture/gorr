package gorr

import (
	"net/http"

	"github.com/Contra-Culture/gorr/node"
	"github.com/Contra-Culture/gorr/url"
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
)

func (c *DispatcherCfgr) Root(t, d string, cfg func(*node.NodeCfgr)) {
	if c.dispatcher.root != nil {
		c.report.Error("root node already specified")
		return
	}
	root, report := node.New(t, d, cfg)
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
	n := d.root
	params, err := url.Handle(
		r.URL.Path,
		func(f string, fn func(string)) {
			n = n.Child(f)
		})
	if err != nil {
		return // TODO:
	}
	m := n.Handler(node.HTTPMethod(r.Method))
	handle := m.Handler()
	handle(w, r, params)
}
