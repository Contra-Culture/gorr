package node

import (
	"net/http"

	"github.com/Contra-Culture/report"
)

type (
	Node struct {
		isParameter bool
		title       string
		matcher     interface{} // interface{} is string, []string or func(string) bool
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
	n = &Node{
		isParameter: false,
		title:       t,
		description: d,
		matcher:     t,
	}
	r = report.New("routing tree")
	cfg(
		&NodeCfgr{
			report: r,
			node:   n,
		})
	return
}
func (n *Node) Handler(m HTTPMethod) func(http.ResponseWriter, *http.Request, map[string]string) {
	method := n.methods[m]
	if method != nil {
		return method.handler
	}
	return n.methodNotAllowedErrorHandler()
}
func (n *Node) methodNotAllowedErrorHandler() func(http.ResponseWriter, *http.Request, map[string]string) {
	return nil
}
func (n *Node) notFoundErrorHandler() func(http.ResponseWriter, *http.Request, map[string]string) {
	return nil
}
func (n *Node) internalServerErrorHandler() func(http.ResponseWriter, *http.Request, map[string]string) {
	return nil
}
func (n *Node) Child(f string) *Node {
	child, ok := n.static[f]
	if ok {
		return child
	}
	return nil
}
