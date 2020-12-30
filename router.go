package gorr

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	Router struct {
		root *Node
		errorHandlers
	}
	errorHandlers    [4]Handler
	NodeProxyContext interface {
		setError(error)
	}
	RouterConfFn func(*RouterProxy)
	RouterProxy  struct {
		router *Router
		err    error
	}
	RouterError int
)

const (
	_ = RouterError(iota)
	NotFoundError
	MethodNotAllowedError
	InternalServerError
)

var (
	rootNodeNotSpecifiedError          = errors.New("root node not specified")
	rootNodeAlreadySpecifiedError      = errors.New("root node already specified")
	rootMatcher                        = MatchesOneOf([]string{"", "/"})
	errorHandlerAlreadySpecifiedErrors = [4]error{
		nil,
		errors.New("`Not Found` handler already specified"),
		errors.New("`Method Not Allowed` handler already specified"),
		errors.New("`Internal Server Error` handler already specified"),
	}
	errorHandlerNotSpecifiedErrors = [4]error{
		nil,
		errors.New("`Not Found` handler not specified"),
		errors.New("`Method Not Allowed` handler not specified"),
		errors.New("`Internal Server Error` handler not specified"),
	}
)

// Returns new router.
func New(conf RouterConfFn) (r *Router, err error) {
	r = &Router{errorHandlers: errorHandlers{}}
	proxy := &RouterProxy{router: r}
	conf(proxy)
	err = proxy.err
	if err != nil {
		r = nil
		return
	}
	if r.root == nil {
		r = nil
		err = rootNodeNotSpecifiedError
		return
	}
	idx := int(NotFoundError)
	if r.errorHandlers[idx] == nil {
		r = nil
		err = errorHandlerNotSpecifiedErrors[idx]
		return
	}
	idx = int(MethodNotAllowedError)
	if r.errorHandlers[idx] == nil {
		r = nil
		err = errorHandlerNotSpecifiedErrors[idx]
		return
	}
	idx = int(InternalServerError)
	if r.errorHandlers[idx] == nil {
		r = nil
		err = errorHandlerNotSpecifiedErrors[idx]
		return
	}
	return
}
func (r *RouterProxy) OnError(re RouterError, h Handler) {
	if r.err != nil {
		return
	}
	idx := int(re)
	if r.router.errorHandlers[idx] != nil {
		r.err = errorHandlerAlreadySpecifiedErrors[idx]
		return
	}
	r.router.errorHandlers[idx] = h
}

// Defines router's root node.
// Router can have only one root node.
func (r *RouterProxy) Root(title, description string, conf NodeConfFn) {
	if r.err != nil {
		return
	}
	if r.router.root != nil {
		r.err = rootNodeAlreadySpecifiedError
		return
	}
	node := &Node{
		title:       title,
		description: description,
		match:       rootMatcher,
	}
	proxy := &NodeProxy{
		node: node,
	}
	conf(proxy)
	if proxy.err != nil {
		r.err = proxy.err
		return
	}
	if node.isEmpty() {
		r.err = obsoleteNodeError
		return
	}
	r.router.root = node
}
func (rr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nserves: %s: %s\n\n", r.Method, r.URL.String())
	chunks := NewChunker(r.URL)
	handler, err := rr.root.Match(r.Method, chunks)
	if err != nil {
		return
	}
	(*handler)(w, r, chunks.Params())
}
