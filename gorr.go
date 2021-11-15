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

	}
	root, report := node.New(t, d, cfg)
	c.dispatcher.root = root
	c.report = report
}
func New(cfg func(*DispatcherCfgr)) (d *Dispatcher) {
	d = &Dispatcher{}
	cfg(
		&DispatcherCfgr{
			dispatcher: d,
		})
	return
}
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fragments := url.New(r.URL)
	if !fragments.HasNext() {
		handle := d.root.Handler(node.HTTPMethod(r.Method))
		handle(w, r, fragments.Params())
	}
	node := d.root
	for fragments.Next() {
		fragment := fragments.Current()
		node.Child(fragment)
	}
}
