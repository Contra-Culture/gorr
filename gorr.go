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
	fmt.Printf("\n\nroot node specified c.dispatcher.root %#v\n\n", c.dispatcher.root)
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
	fmt.Printf("\n\nServeHTTP dispatcher %#v, %s\n\n", d, r.URL.String())
	var current = d.root
	var ok bool
	params, err := handle(
		r.URL.Path,
		func(f string, fn func(string)) {
			if current != nil {
				fmt.Printf("\n\nServeHTTP iteration over path (%s) `%s` -> node %#v\n\n", r.URL.Path, f, current)
				current, ok = current.Child(f)
				if !ok {
					w.Write([]byte("not found"))
					w.WriteHeader(404)
					return
				}
				pname, ok := current.Param()
				if ok {
					fn(pname)
				}
				return
			}
		})
	if err != nil {
		w.Write([]byte("not-found 2"))
		w.WriteHeader(404)
		return
	}
	fmt.Printf("\n\nServeHTTP iteration over path (%s) end node %#v\n\n", r.URL.Path, current)
	m := current.Handler(node.HTTPMethod(r.Method))
	handle := m.Handler()
	handle(w, r, params)
}
func handle(path string, iterBlck func(string, func(string))) (params map[string]string, err error) {
	params = map[string]string{
		"$path": path,
	}
	fragments := []string{}
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
		iterBlck(
			fragment,
			func(k string) {
				_, exists := params[k]
				if exists {
					err = fmt.Errorf("parameter \"%s\" already marked", k)
					return
				}
				params[k] = fragment
			})
	}
	return
}
