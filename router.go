package gorr

import (
	"errors"
	"net/http"
	"strings"
)

type (
	Router struct {
		isRootSet    bool
		root         *Node
		before       *http.HandlerFunc
		beforeMethod *Handler
		afterMethod  *Handler
		after        *http.HandlerFunc
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
	rootNodeNotSpecifiedError             = errors.New("root node not specified")
	rootNodeAlreadySpecifiedError         = errors.New("root node already specified")
	beforeHookAlreadySpecifiedError       = errors.New("before hook already specified")
	afterHookAlreadySpecifiedError        = errors.New("after hook already specified")
	beforeMethodHookAlreadySpecifiedError = errors.New("beforeMethod hook already specified")
	afterMethodHookAlreadySpecifiedError  = errors.New("afterMethod hook already specified")
	rootMatcher                           = MatchesOneOf([]string{"", "/"})
	errorHandlerAlreadySpecifiedErrors    = [4]error{
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
func New(conf RouterConfFn) (rr *Router, err error) {
	root := &Node{
		header: NodeHeader{
			title: "/",
			match: rootMatcher,
		},
		description: "root node",
		children:    []*Node{},
	}
	rr = &Router{
		isRootSet:     false,
		root:          root,
		errorHandlers: errorHandlers{},
	}
	proxy := &RouterProxy{router: rr}
	conf(proxy)
	err = proxy.err
	if err != nil {
		rr = nil
		return
	}
	if !rr.isRootSet {
		rr = nil
		err = rootNodeNotSpecifiedError
		return
	}
	idx := int(NotFoundError)
	if rr.errorHandlers[idx] == nil {
		rr = nil
		err = errorHandlerNotSpecifiedErrors[idx]
		return
	}
	idx = int(MethodNotAllowedError)
	if rr.errorHandlers[idx] == nil {
		rr = nil
		err = errorHandlerNotSpecifiedErrors[idx]
		return
	}
	idx = int(InternalServerError)
	if rr.errorHandlers[idx] == nil {
		rr = nil
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
func (r *RouterProxy) Before(h http.HandlerFunc) {
	if r.err != nil {
		return
	}
	if r.router.before != nil {
		r.err = beforeHookAlreadySpecifiedError
	}
	r.router.before = &h
}
func (r *RouterProxy) After(h http.HandlerFunc) {
	if r.err != nil {
		return
	}
	if r.router.after != nil {
		r.err = afterHookAlreadySpecifiedError
	}
	r.router.after = &h
}
func (r *RouterProxy) BeforeMethod(h Handler) {
	if r.err != nil {
		return
	}
	if r.router.beforeMethod != nil {
		r.err = beforeMethodHookAlreadySpecifiedError
	}
	r.router.beforeMethod = &h
}
func (r *RouterProxy) AfterMethod(h Handler) {
	if r.err != nil {
		return
	}
	if r.router.afterMethod != nil {
		r.err = afterMethodHookAlreadySpecifiedError
	}
	r.router.afterMethod = &h
}

// Defines router's root node.
// Router can have only one root node.
func (r *RouterProxy) Root(title, description string, conf NodeConfFn) {
	if r.err != nil {
		return
	}
	if r.router.isRootSet {
		r.err = rootNodeAlreadySpecifiedError
		return
	}
	node := r.router.root
	node.header.title = title
	node.description = description
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
	r.router.isRootSet = true
}
func (r *RouterProxy) ShowRoutes(t string, cond func() bool) {
	if r.err != nil {
		return
	}
	n := &NodeProxy{node: r.router.root}
	n.Node(StaticConditional(t, cond), "routing documentation", func(n *NodeProxy) {
		n.Method("get-routes", "provides app's routing information", GET, func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
			w.Write([]byte(r.router.toJSON()))
			w.WriteHeader(http.StatusOK)
		})
	})
	r.err = n.err
}

func (r *Router) toJSON() string {
	var sb = &strings.Builder{}
	sb.WriteByte('{')
	sb.WriteString(r.root.toJSON())
	sb.WriteByte('}')
	return sb.String()
}
func (rr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(*rr.before)(w, r)
	chunks := NewChunker(r.URL)
	handler, err := rr.root.Match(r.Method, chunks)
	if err != nil {
		//	return
	}
	if handler == nil {
		handler = &rr.errorHandlers[int(NotFoundError)]
	}
	ps := chunks.Params()
	(*rr.beforeMethod)(w, r, ps)
	(*handler)(w, r, ps)
	(*rr.afterMethod)(w, r, ps)
	(*rr.after)(w, r)
}
