package gorr

import (
	"fmt"
	"net/http"
	"strings"

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
)

func (c *DispatcherCfgr) Root(d string, cfg func(*node.NodeCfgr)) {
	if c.dispatcher.root != nil {
		c.report.Error("root node already specified")
		return
	}
	root, report := node.New("/", d, cfg)
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
	var (
		path      = r.URL.Path
		current   = d.root
		parent    = current
		ok        bool
		err       error
		fragments = []string{}
		params    = map[string]string{
			"$path": path,
		}
	)
	for _, f := range strings.Split(path, "/") {
		if len(f) > 0 {
			fragments = append(fragments, f)
		}
	}
	for _, fragment := range fragments {
		if err != nil {
			params = nil
			return
		}
		current, ok = parent.Child(fragment)
		if !ok {
			current.HandleNotFoundError(w, r, params)
			return
		}
		pname, ok := current.Param()
		if ok {
			_, exists := params[pname]
			if exists {
				// TODO:
				panic(fmt.Errorf("parameter \"%s\" already marked", pname))
			}
			params[pname] = fragment
		}
		parent = current
	}
	current.Handle(w, r, params)
}
