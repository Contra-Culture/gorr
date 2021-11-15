package node

import (
	"fmt"
	"net/http"

	"github.com/Contra-Culture/report"
)

type (
	Node struct {
		matcher     interface{} // interface{} is string, []string or func(string) bool
		title       string
		description string
		methods     map[HTTPMethod]*Method
		static      map[string]*Node
		param       *Node
		wildcard    *Node
	}
	Method struct {
		title       string
		description string
		handler     func(http.ResponseWriter, *http.Request, map[string]string)
	}
	HTTPMethod string
)

const (
	GET     = "GET"
	HEAD    = "HEAD"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	CONNECT = "CONNECT"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
	PATCH   = "PATCH"
)

func New(t, d string, cfg func(*NodeCfgr)) (n *Node, r *report.RContext) {
	r = report.New(t)
	n = new(t, d, r, cfg)
	return
}
func new(t, d string, rctx *report.RContext, cfg func(*NodeCfgr)) (n *Node) {
	n = &Node{
		title:       t,
		description: d,
		matcher:     t,
		methods:     map[HTTPMethod]*Method{},
		static:      map[string]*Node{},
	}
	cfg(
		&NodeCfgr{
			report: rctx,
			node:   n,
		})
	return
}
func (n *Node) Handler(m HTTPMethod) *Method {
	return n.methods[m]
}
func (n *Node) Child(f string) (child *Node, ok bool) {
	fmt.Printf("\n\nnode.Child() parent node: %#v\n\n", n)
	child, ok = n.static[f]
	if ok {
		return
	}
	child = n.param
	if child != nil {
		return
	}
	child = n.wildcard
	return child, n.wildcard != nil
}
func (m *Method) Title() string {
	return m.title
}
func (m *Method) Decription() string {
	return m.description
}
func (m *Method) Handler() func(http.ResponseWriter, *http.Request, map[string]string) {
	return m.handler
}
