package node

import (
	"errors"
	"net/http"

	"github.com/Contra-Culture/report"
)

type (
	Handler func(http.ResponseWriter, *http.Request, map[string]string) error
	Node    struct {
		parent                         *Node
		isParam                        bool
		title                          string
		description                    string
		methods                        map[HTTPMethod]*Method
		static                         map[string]*Node
		param                          *Node
		wildcard                       *Node
		__notFoundErrorHandler         Handler
		__methodNotAllowedErrorHandler Handler
		__internalServerErrorHandler   Handler
	}
	Method struct {
		title       string
		description string
		handler     Handler
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

func New(cfg func(*NodeCfgr)) (n *Node, r *report.RContext) {
	r = report.New("root")
	n = new(nil, false, r, cfg)
	return
}
func new(p *Node, isProp bool, rctx *report.RContext, cfg func(*NodeCfgr)) (n *Node) {
	n = &Node{
		parent:  p,
		isParam: isProp,
		methods: map[HTTPMethod]*Method{},
		static:  map[string]*Node{},
	}
	cfg(
		&NodeCfgr{
			report: rctx,
			node:   n,
		})
	return
}
func (n *Node) Param() (name string, ok bool) {
	ok = n.isParam
	if ok {
		name = n.title
	}
	return
}
func (n *Node) Handle(w http.ResponseWriter, r *http.Request, params map[string]string) {
	method, ok := n.methods[HTTPMethod(r.Method)]
	if !ok {
		n.handleMethodNotAllowedError(w, r, params)
	}
	err := method.handler(w, r, params)
	if err != nil {
		n.handleInternalServerError(w, r, params)
	}
}
func (n *Node) handleMethodNotAllowedError(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
	var handle Handler
	for {
		handle = n.__methodNotAllowedErrorHandler
		if handle != nil {
			handle(w, r, params)
			return
		}
		n = n.parent
		if n != nil {
			continue
		}
		return errors.New("method not allowed handler is not provided")
	}
}
func (n *Node) HandleNotFoundError(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
	var handle Handler
	for {
		handle = n.__notFoundErrorHandler
		if handle != nil {
			handle(w, r, params)
			return
		}
		n = n.parent
		if n != nil {
			continue
		}
		return errors.New("not found handler is not provided")
	}
}
func (n *Node) handleInternalServerError(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
	var handle Handler
	for {
		handle = n.__internalServerErrorHandler
		if handle != nil {
			handle(w, r, params)
			return
		}
		n = n.parent
		if n != nil {
			continue
		}
		return errors.New("internal server error handler is not provided")
	}
}
func (n *Node) Child(f string) (child *Node, ok bool) {
	child = n.static[f]
	ok = child != nil
	if ok {
		return
	}
	child = n.param
	ok = child != nil
	if ok {
		return
	}
	child = n.wildcard
	ok = child != nil
	return
}
func (m *Method) Title() string {
	return m.title
}
func (m *Method) Decription() string {
	return m.description
}
func (m *Method) Handler() Handler {
	return m.handler
}
