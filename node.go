package gorr

import (
	"errors"
	"net/http"
)

type (
	Handler    func(http.ResponseWriter, *http.Request, map[string]string)
	HTTPMethod int
	Matcher    func(string) bool
	Node       struct {
		title       string
		description string
		match       Matcher
		methods     [10]*Handler
		children    []*Node
	}

	NodeConfFn func(*NodeProxy)
	NodeProxy  struct {
		node *Node
		err  error
	}
)

const (
	WRONG_METHOD = HTTPMethod(iota)
	GET
	HEAD
	POST
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
	PATCH
)

var (
	methodsAlreadyDefinedError     = errors.New("methods already defined")
	obsoleteNodeError              = errors.New("node should have at least one method/handler or a child node")
	handlersAlreadySpecifiedErrors = [10]error{
		nil,
		errors.New("GET handler already specified"),
		errors.New("HEAD handler already specified"),
		errors.New("POST handler already specified"),
		errors.New("PUT handler already specified"),
		errors.New("DELETE handler already specified"),
		errors.New("CONNECT handler already specified"),
		errors.New("OPTIONS handler already specified"),
		errors.New("TRACE handler already specified"),
		errors.New("PATCH handler already specified"),
	}
)

func StringToMethod(s string) (m HTTPMethod) {
	switch s {
	case "GET":
		m = GET
	case "HEAD":
		m = HEAD
	case "POST":
		m = POST
	case "PUT":
		m = PUT
	case "DELETE":
		m = DELETE
	case "CONNECT":
		m = CONNECT
	case "OPTIONS":
		m = OPTIONS
	case "TRACE":
		m = TRACE
	case "PATCH":
		m = PATCH
	default:
		m = WRONG_METHOD
	}
	return
}
func Matches(expected string) Matcher {
	return func(v string) bool {
		return expected == v
	}
}
func MatchesOneOf(samples []string) Matcher {
	return func(v string) bool {
		for _, s := range samples {
			if s == v {
				return true
			}
		}
		return false
	}
}
func (n *Node) Match(method string, chunks *Chunker) (h *Handler, err error) {
	ch := chunks.Chunk()
	if !n.match(ch) {
		return
	}
	err = chunks.Set(n.title, ch)
	if err != nil {
		h = nil
		return
	}
	hasNext := chunks.Next()
	if !hasNext {
		h = n.methods[StringToMethod(method)]
		return
	}
	for _, ch := range n.children {
		h, err = ch.Match(method, chunks)
		if err != nil {
			h = nil
			return
		}
		if h != nil {
			return
		}
	}
	return
}
func (p *NodeProxy) Node(title, description string, match Matcher, conf NodeConfFn) {
	if p.err != nil {
		return
	}
	node := &Node{
		title:       title,
		description: description,
		match:       match,
	}
	proxy := &NodeProxy{node: node}
	conf(proxy)
	if proxy.err != nil {
		p.err = proxy.err
		return
	}
	if node.isEmpty() {
		p.err = obsoleteNodeError
		return
	}
	p.node.children = append(p.node.children, node)
}
func (n *NodeProxy) Method(m HTTPMethod, h Handler) {
	if n.err != nil {
		return
	}
	idx := int(m)
	if n.node.methods[idx] != nil {
		n.err = handlersAlreadySpecifiedErrors[idx]
		return
	}
	n.node.methods[idx] = &h
}
func (n *Node) isEmpty() bool {
	if len(n.children) > 0 {
		return false
	}
	for _, elem := range n.methods {
		if elem != nil {
			return false
		}
	}
	return true
}
