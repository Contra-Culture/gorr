package gorr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type (
	NodeHeader struct {
		isParameter bool
		title       string
		match       Matcher
	}
	Handler func(http.ResponseWriter, *http.Request, map[string]string)
	Method  struct {
		title       string
		description string
		handler     Handler
	}
	HTTPMethod int
	Matcher    func(string) bool
	Node       struct {
		header      NodeHeader
		description string
		methods     [10]*Method
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

func (n *Node) Match(method string, chunks *Chunker) (h *Handler, err error) {
	ch := chunks.Chunk()
	if !n.header.match(ch) {
		return
	}
	if n.header.isParameter {
		err = chunks.Set(n.header.title, ch)
		if err != nil {
			h = nil
			return
		}
	}
	hasNext := chunks.Next()
	if !hasNext {
		m := n.methods[StringToMethod(method)]
		if m != nil {
			chunks.Set("$method", m.title)
			h = &m.handler
			return
		}
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
func (p *NodeProxy) Node(h NodeHeader, description string, conf NodeConfFn) {
	if p.err != nil {
		return
	}
	node := &Node{
		header:      h,
		description: description,
		children:    []*Node{},
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
func (n *NodeProxy) Method(t, d string, m HTTPMethod, h Handler) {
	if n.err != nil {
		return
	}
	idx := int(m)
	if n.node.methods[idx] != nil {
		n.err = handlersAlreadySpecifiedErrors[idx]
		return
	}
	n.node.methods[idx] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
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
func (n *Node) toJSON() string {
	var sb = &strings.Builder{}
	sb.WriteRune('{')
	sb.WriteString(fmt.Sprintf("\"title\":\"%s\",", n.header.title))
	sb.WriteString(fmt.Sprintf("\"isParameter\":\"%v\",", n.header.isParameter))
	sb.WriteString(fmt.Sprintf("\"description\":\"%s\",", n.description))
	sb.WriteString("\"children\":[")
	skipCommaAfter := len(n.children) - 1
	for i, ch := range n.children {
		sb.WriteString(ch.toJSON())
		if i < skipCommaAfter {
			sb.WriteRune(',')
		}
	}
	sb.WriteString("]}")
	return sb.String()
}
