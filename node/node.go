package node

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Contra-Culture/report"
)

type (
	Handler func(http.ResponseWriter, *http.Request, Params) error
	Node    struct {
		typ                            NodeType
		parent                         *Node
		title                          string
		description                    string
		methods                        map[HTTPMethod]*Method
		static                         map[string]*Node
		param                          *Node
		wildcard                       *Node
		matcher                        interface{} // interface{} is string, map[string]bool, func(string) bool, Query
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
	NodeType   int
	Query      func(Params) (interface{}, error)
	params     map[string]interface{}
	Params     interface {
		Get(string) (interface{}, bool)
		Set(string, interface{}) error
	}
)

const (
	_ NodeType = iota
	STATIC
	STRING_PARAM
	ID_PARAM
	VARIANT_PARAM
	WILDCARD
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

func NewParams() Params {
	return params(map[string]interface{}{})
}
func (ps params) Get(pn string) (v interface{}, ok bool) {
	v, ok = ps[pn]
	return
}
func (ps params) Set(pn string, v interface{}) (err error) {
	_, exists := ps[pn]
	if exists {
		err = fmt.Errorf("parameter \"%s\" already setted", pn)
		return
	}
	ps[pn] = v
	return
}
func New(cfg func(*StaticNodeCfgr)) (n *Node, r *report.RContext) {
	r = report.New("root")
	n = new(nil, STATIC)
	c := &StaticNodeCfgr{
		NodeCfgr{
			node:   n,
			report: r,
		},
	}
	cfg(c)
	return
}
func new(p *Node, t NodeType) (n *Node) {
	n = &Node{
		parent:  p,
		typ:     t,
		methods: map[HTTPMethod]*Method{},
		static:  map[string]*Node{},
	}
	return
}
func (n *Node) Param() (string, bool) {
	switch n.typ {
	case STRING_PARAM, ID_PARAM, VARIANT_PARAM:
		return n.title, true
	default:
		return "", false
	}
}
func (n *Node) Handle(w http.ResponseWriter, r *http.Request) {
	fragments := []string{}
	for _, f := range strings.Split(r.URL.Path, "/") {
		if len(f) > 0 {
			fragments = append(fragments, f)
		}
	}
	params := NewParams()
	if len(fragments) == 0 {
		method, ok := n.methods[HTTPMethod(r.Method)]
		if !ok {
			n.handleMethodNotAllowedError(w, r, params)
		}
		err := method.handler(w, r, params)
		if err != nil {
			n.handleInternalServerError(w, r, params)
		}
		return
	}
	var err error
	var ok bool
	var parent = n
	for ; len(fragments) > 0; fragments = fragments[1:] {
		fmt.Printf("\nfragments: %#v", fragments)
		f := fragments[0]
		n, ok = parent.static[f]
		fmt.Printf("\n\tstatic node: %#v", n)
		if ok {
			parent = n
			continue
		}
		n = parent.param
		fmt.Printf("\n\tparam node: %#v", n)
		if n != nil {
			switch matcher := n.matcher.(type) {
			case map[string]bool:
				params.Set(n.title, f)
				if matcher[f] {
					parent = n
					continue
				}
				fmt.Printf("\n\t\tquery not found--1")
				n.HandleNotFoundError(w, r, params)
				return
			case func(string) bool:
				params.Set(n.title, f)
				if matcher(f) {
					parent = n
					continue
				}
				fmt.Printf("\n\t\tquery not found-0")
				n.HandleNotFoundError(w, r, params)
				return
			case Query:
				params.Set(n.title, f)
				v, err := matcher(params)
				if err != nil {
					fmt.Printf("\n\t\tquery not found-1")
					n.HandleNotFoundError(w, r, params)
					return
				}
				params.Set(n.title, v) // should be different key
				parent = n
				continue
			default:
				params.Set(n.title, f)
				fmt.Printf("\n\t\tdefault section")
				parent = n
				continue
			}
		} else {
			n = parent.wildcard
			fmt.Printf("\n\twildcard node: %#v", n)
			if n != nil {
				parent = n
				continue
			}
			fmt.Printf("\n\t\tquery not found-2")
			parent.HandleNotFoundError(w, r, params)
			return
		}
		method, ok := n.methods[HTTPMethod(r.Method)]
		if !ok {
			n.handleMethodNotAllowedError(w, r, params)
		}
		err = method.Handler()(w, r, params)
		if err != nil {
			n.handleMethodNotAllowedError(w, r, params)
		}
		return
	}
	method, ok := n.methods[HTTPMethod(r.Method)]
	if !ok {
		n.handleMethodNotAllowedError(w, r, params)
		return
	}
	err = method.handler(w, r, params)
	if err != nil {
		n.handleInternalServerError(w, r, params)
		return
	}
}
func (n *Node) handleMethodNotAllowedError(w http.ResponseWriter, r *http.Request, params Params) (err error) {
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
func (n *Node) HandleNotFoundError(w http.ResponseWriter, r *http.Request, params Params) (err error) {
	var handle Handler
	for {
		handle = n.__notFoundErrorHandler
		if handle != nil {
			err = handle(w, r, params)
			return
		}
		n = n.parent
		if n != nil {
			continue
		}
		return errors.New("not found handler is not provided")
	}
}
func (n *Node) handleInternalServerError(w http.ResponseWriter, r *http.Request, params Params) (err error) {
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
func (m *Method) Title() string {
	return m.title
}
func (m *Method) Decription() string {
	return m.description
}
func (m *Method) Handler() Handler {
	return m.handler
}
