package gorr

import (
	"github.com/Contra-Culture/report"
)

type (
	NodeCfgr struct {
		node   *Node
		report report.Node
	}
	WildcardNodeCfgr struct {
		NodeCfgr
	}
	StaticNodeCfgr struct {
		NodeCfgr
	}
	StringParamNodeCfgr struct {
		NodeCfgr
	}
	IDParamNodeCfgr struct {
		NodeCfgr
	}
	VariantParamNodeCfgr struct {
		NodeCfgr
	}
	MethodCfgr struct {
		method *Method
		report report.Node
	}
)

func (c *NodeCfgr) Title(t string) {
	if len(c.node.title) > 0 {
		c.report.Error("title already specified")
	}
	c.node.title = t
}
func (c *NodeCfgr) Description(d string) {
	if len(c.node.description) > 0 {
		c.report.Error("description already specified")
	}
	c.node.description = d
}
func (c *NodeCfgr) BeforeDo(h Handler) {
	if c.node.beforeHandler != nil {
		c.report.Error("before hook already specified")
		return
	}
	c.node.beforeHandler = h
}
func (c *NodeCfgr) AfterDo(h Handler) {
	if c.node.afterHandler != nil {
		c.report.Error("after hook already specified")
		return
	}
	c.node.afterHandler = h
}
func (c *NodeCfgr) InheritableBeforeDo(h Handler) {
	if c.node.inheritableBeforeHandler != nil {
		c.report.Error("inheritable before hook already specified")
		return
	}
	c.node.inheritableBeforeHandler = h
}
func (c *NodeCfgr) InheritableAfterDo(h Handler) {
	if c.node.inheritableAfterHandler != nil {
		c.report.Error("inheritable after hook already specified")
		return
	}
	c.node.inheritableAfterHandler = h
}
func (c *NodeCfgr) Wildcard(cfg func(*WildcardNodeCfgr)) {
	if c.node.wildcard != nil {
		c.report.Error("* node already specified")
		return
	}
	rctx := c.report.Structure("*")
	n := new(c.node, WILDCARD)
	nc := NodeCfgr{
		node:   n,
		report: rctx,
	}
	cfg(&WildcardNodeCfgr{nc})
	nc.check()
	c.node.wildcard = n
}
func (c *NodeCfgr) Static(f string, cfg func(*StaticNodeCfgr)) {
	_, exists := c.node.static[f]
	if exists {
		c.report.Error("static \"%s\" node already specified", f)
		return
	}
	rctx := c.report.Structure("%%%s", f)
	n := new(c.node, STATIC)
	nc := NodeCfgr{
		node:   n,
		report: rctx,
	}
	cfg(&StaticNodeCfgr{nc})
	nc.check()
	c.node.static[f] = n
}
func (c *NodeCfgr) StringParam(name string, cfg func(*StringParamNodeCfgr)) {
	if c.node.param != nil {
		c.report.Error("param \":%s\" node already specified", name)
		return
	}
	rctx := c.report.Structure(":%s", name)
	n := new(c.node, STRING_PARAM)
	nc := NodeCfgr{
		node:   n,
		report: rctx,
	}
	cfg(&StringParamNodeCfgr{nc})
	nc.check()
	c.node.param = n
}
func (c *NodeCfgr) IDParam(name string, cfg func(*IDParamNodeCfgr)) {
	if c.node.param != nil {
		c.report.Error("param \":%s\" node already specified", name)
		return
	}
	rctx := c.report.Structure(":%s", name)
	n := new(c.node, ID_PARAM)
	nc := NodeCfgr{
		node:   n,
		report: rctx,
	}
	cfg(&IDParamNodeCfgr{nc})
	nc.check()
	c.node.param = n
}
func (c *NodeCfgr) VariantParam(name string, cfg func(*VariantParamNodeCfgr)) {
	if c.node.param != nil {
		c.report.Error("param \":%s\" node already specified", name)
		return
	}
	rctx := c.report.Structure(":%s", name)
	n := new(c.node, VARIANT_PARAM)
	nc := NodeCfgr{
		node:   n,
		report: rctx,
	}
	cfg(&VariantParamNodeCfgr{nc})
	nc.check()
	c.node.param = n
}
func (c *NodeCfgr) HandleNotFoundErrorWith(h Handler) {
	if c.node.__notFoundErrorHandler != nil {
		c.report.Error("not found error handler already specified")
		return
	}
	c.node.__notFoundErrorHandler = h
}
func (c *NodeCfgr) HandleInternalServerErrorWith(h Handler) {
	if c.node.__internalServerErrorHandler != nil {
		c.report.Error("internal server error handler already specified")
		return
	}
	c.node.__internalServerErrorHandler = h
}
func (c *NodeCfgr) HandleMethodNotAllowedErrorWith(h Handler) {
	if c.node.__methodNotAllowedErrorHandler != nil {
		c.report.Error("method not allowed error handler already specified")
		return
	}
	c.node.__methodNotAllowedErrorHandler = h
}
func (c *NodeCfgr) GET(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[GET]
	if exists {
		c.report.Error("GET handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("GET"),
		})
	c.node.methods[GET] = m
}
func (c *NodeCfgr) POST(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[POST]
	if exists {
		c.report.Error("POST handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("POST"),
		})
	c.node.methods[POST] = m
}
func (c *NodeCfgr) PUT(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[PUT]
	if exists {
		c.report.Error("PUT handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("PUT"),
		})
	c.node.methods[PUT] = m
}
func (c *NodeCfgr) PATCH(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[PATCH]
	if exists {
		c.report.Error("PATCH handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("PATCH"),
		})
	c.node.methods[PATCH] = m
}
func (c *NodeCfgr) DELETE(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[DELETE]
	if exists {
		c.report.Error("DELETE handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("DELETE"),
		})
	c.node.methods[DELETE] = m
}
func (c *NodeCfgr) HEAD(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[HEAD]
	if exists {
		c.report.Error("HEAD handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("HEAD"),
		})
	c.node.methods[HEAD] = m
}
func (c *NodeCfgr) CONNECT(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[CONNECT]
	if exists {
		c.report.Error("CONNECT handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("CONNECT"),
		})
	c.node.methods[CONNECT] = m
}
func (c *NodeCfgr) OPTIONS(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[OPTIONS]
	if exists {
		c.report.Error("OPTIONS handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("OPTIONS"),
		})
	c.node.methods[OPTIONS] = m
}
func (c *NodeCfgr) TRACE(cfg func(*MethodCfgr)) {
	_, exists := c.node.methods[TRACE]
	if exists {
		c.report.Error("TRACE handler already specified")
		return
	}
	m := &Method{}
	cfg(
		&MethodCfgr{
			method: m,
			report: c.report.Structure("PUT"),
		})
	c.node.methods[TRACE] = m
}
func (c *NodeCfgr) check() {
	if len(c.node.title) == 0 {
		c.report.Error("node title is not specified")
	}
	if len(c.node.description) == 0 {
		c.report.Error("node description is not specified")
	}
	if len(c.node.methods) == 0 && len(c.node.static) == 0 && c.node.param == nil && c.node.wildcard == nil {
		c.report.Error("node has neither methods nor children nodes specified")
	}
}
func (c *StaticNodeCfgr) Fragment(f string) {
	if c.node.matcher != nil {
		c.report.Error("fragment already specified")
		return
	}
	c.node.matcher = f
}
func (c *StringParamNodeCfgr) Matcher(m func(string) bool) {
	if c.node.matcher != nil {
		c.report.Error("fragment already specified")
		return
	}
	c.node.matcher = m
}
func (c *VariantParamNodeCfgr) Variant(v string) {
	variants := c.node.matcher.(map[string]bool)
	if variants[v] {
		c.report.Error("variant \"%s\"already specified", v)
		return
	}
	variants[v] = true
}
func (c *IDParamNodeCfgr) Query(q Query) {
	if c.node.matcher != nil {
		c.report.Error("param already specified")
		return
	}
	c.node.matcher = q
}
func (c *MethodCfgr) Title(t string) {
	if len(c.method.title) > 0 {
		c.report.Error("title already specified")
		return
	}
	c.method.title = t
}
func (c *MethodCfgr) Description(d string) {
	if len(c.method.description) > 0 {
		c.report.Error("description already specified")
		return
	}
	c.method.description = d
}
func (c *MethodCfgr) Handler(h Handler) {
	if c.method.handler != nil {
		c.report.Error("handler already specified")
		return
	}
	c.method.handler = h
}
