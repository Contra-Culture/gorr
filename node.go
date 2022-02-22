package gorr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Contra-Culture/report"
)

type (
	Handler func(report.Node, http.ResponseWriter, *http.Request, Params) error
	Node    struct {
		typ                            NodeType
		parent                         *Node
		inheritedBeforeHandlers        []Handler
		inheritableBeforeHandler       Handler
		beforeHandler                  Handler
		afterHandler                   Handler
		inheritableAfterHandler        Handler
		inheritedAfterHandlers         []Handler
		methods                        map[HTTPMethod]*Method
		static                         map[string]*Node // nested static fragment nodes (Priority: 1)
		param                          *Node            // nested param node (Priority: 2)
		wildcard                       *Node            // nested wildcard node (Priority: 3)
		matcher                        interface{}      // interface{} is string, map[string]bool, func(string) bool or Query
		__notFoundErrorHandler         Handler
		__methodNotAllowedErrorHandler Handler
		__internalServerErrorHandler   Handler
		title                          string
		description                    string
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
func new(p *Node, t NodeType) (n *Node) {
	n = &Node{
		parent:  p,
		typ:     t,
		methods: map[HTTPMethod]*Method{},
		static:  map[string]*Node{},
	}
	if p != nil {
		if p.inheritableBeforeHandler != nil {
			n.inheritedBeforeHandlers = append(p.inheritedBeforeHandlers, p.inheritableBeforeHandler)
		}
		if p.inheritableAfterHandler != nil {
			n.inheritedAfterHandlers = append(p.inheritedAfterHandlers, p.inheritableAfterHandler)
		}
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

// Handles request and/or delegates it to its child.
func (n *Node) Handle(rep report.Node, w http.ResponseWriter, r *http.Request) {
	var (
		ok        bool
		params    = NewParams()
		parent    = n
		fragments = []string{}
	)
	for _, f := range strings.Split(r.URL.Path, "/") {
		if len(f) > 0 {
			fragments = append(fragments, f)
		}
	}
	for ; len(fragments) > 0; fragments = fragments[1:] {
		f := fragments[0]
		rep.Info("fragment %s of fragments: %s", f, strings.Join(fragments, "/"))
		n, ok = parent.static[f]
		if ok {
			rep.Info("static node `%s` picked", f)
			parent = n
			continue
		}
		n = parent.param
		if n != nil {
			rep.Info("param node `:%s` picked", f)
			switch matcher := n.matcher.(type) {
			case map[string]bool:
				rep.Info("map[string]bool matcher picked")
				params.Set(n.title, f)
				if matcher[f] {
					parent = n
					continue
				}
				n.handleNotFoundError(rep, w, r, params)
				return
			case func(string) bool:
				rep.Info("func(string) bool matcher picked")
				params.Set(n.title, f)
				if matcher(f) {
					parent = n
					continue
				}
				n.handleNotFoundError(rep, w, r, params)
				return
			case Query:
				rep.Info("Query matcher picked")
				idParamName := fmt.Sprintf("%sID", n.title)
				params.Set(idParamName, f)
				v, err := matcher(params)
				if err != nil {
					n.handleNotFoundError(rep, w, r, params)
					return
				}
				params.Set(n.title, v)
				parent = n
				continue
			default:
				rep.Error("wrong param matcher type")
				n.handleNotFoundError(rep, w, r, params)
				return
			}
		} else {
			n = parent.wildcard
			if n != nil {
				parent = n
				continue
			}
			parent.handleNotFoundError(rep, w, r, params)
			return
		}
	}
	child := rep.Structure("node (%s:%s)", n.title, nodeTypeString(n.typ))
	n.handle(child, w, r, params)
	child.Finalize()
	rep.Finalize()
}
func (n *Node) handle(rep report.Node, w http.ResponseWriter, r *http.Request, params Params) {
	var (
		err error
		h   Handler
	)
	method, ok := n.methods[HTTPMethod(r.Method)]
	if !ok {
		rep.Error("method not allowed %s", r.Method)
		n.handleMethodNotAllowedError(rep, w, r, params)
		return
	}
	for _, h = range n.inheritedBeforeHandlers {
		err = h(rep, w, r, params)
		if err != nil {
			rep.Error("internal server error: %s", err.Error())
			n.handleInternalServerError(rep, w, r, params)
			return
		}
	}
	h = n.inheritableBeforeHandler
	if h != nil {
		err = h(rep, w, r, params)
		if err != nil {
			rep.Error("internal server error: %s", err.Error())
			n.handleInternalServerError(rep, w, r, params)
			return
		}
	}
	h = n.beforeHandler
	if h != nil {
		err = h(rep, w, r, params)
		if err != nil {
			rep.Error("internal server error: %s", err.Error())
			n.handleInternalServerError(rep, w, r, params)
			return
		}
	}
	err = method.Handler()(rep, w, r, params)
	if err != nil {
		rep.Error("internal server error %s", r.Method)
		n.handleInternalServerError(rep, w, r, params)
		return
	}
	h = n.afterHandler
	if h != nil {
		err = h(rep, w, r, params)
		if err != nil {
			rep.Error("internal server error: %s", err.Error())
			n.handleInternalServerError(rep, w, r, params)
			return
		}
	}
	h = n.inheritableAfterHandler
	if h != nil {
		err = h(rep, w, r, params)
		if err != nil {
			rep.Error("internal server error: %s", err.Error())
			n.handleInternalServerError(rep, w, r, params)
			return
		}
	}
	for i := len(n.inheritedAfterHandlers) - 1; i >= 0; i-- {
		h = n.inheritedAfterHandlers[i]
		err = h(rep, w, r, params)
		if err != nil {
			rep.Error("internal server error: %s", err.Error())
			n.handleInternalServerError(rep, w, r, params)
			return
		}
	}
}
func (n *Node) handleMethodNotAllowedError(rep report.Node, w http.ResponseWriter, r *http.Request, params Params) (err error) {
	var handle Handler
	for {
		handle = n.__methodNotAllowedErrorHandler
		if handle != nil {
			rep.Info("method not allowed error handling: handler not found")
			handle(rep, w, r, params)
			return
		}
		n = n.parent
		if n != nil {
			continue
		}
		rep.Info("method not allowed error handling: handler not found")
		return errors.New("method not allowed handler is not provided")
	}
}
func (n *Node) handleNotFoundError(rep report.Node, w http.ResponseWriter, r *http.Request, params Params) (err error) {
	var handle Handler
	for {
		handle = n.__notFoundErrorHandler
		if handle != nil {
			rep.Info("not found error handling: handler not found")
			err = handle(rep, w, r, params)
			return
		}
		n = n.parent
		if n != nil {
			continue
		}
		rep.Info("not found error handling: handler not found")
		return errors.New("not found handler is not provided")
	}
}
func (n *Node) handleInternalServerError(rep report.Node, w http.ResponseWriter, r *http.Request, params Params) (err error) {
	var handle Handler
	for {
		handle = n.__internalServerErrorHandler
		if handle != nil {
			rep.Info("internal server error handling")
			handle(rep, w, r, params)
			return
		}
		n = n.parent
		if n != nil {
			continue
		}
		rep.Info("internal server error handling: handler not found")
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
func nodeTypeString(t NodeType) string {
	switch t {
	case STATIC:
		return "static"
	case STRING_PARAM:
		return "string-param"
	case ID_PARAM:
		return "ID-param"
	case VARIANT_PARAM:
		return "variant-param"
	case WILDCARD:
		return "*"
	default:
		panic(fmt.Sprintf("wrong node type %d", int(t))) // should not occure
	}
}
