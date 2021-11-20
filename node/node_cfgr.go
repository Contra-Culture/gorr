package node

import (
	"fmt"

	"github.com/Contra-Culture/report"
)

type (
	NodeCfgr struct {
		node   *Node
		report *report.RContext
	}
	MethodCfgr struct {
		method *Method
		report *report.RContext
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
func (c *NodeCfgr) Wildcard(cfg func(*NodeCfgr)) {
	if c.node.wildcard != nil {
		c.report.Error("* node already specified")
		return
	}
	rctx := c.report.Context("*")
	n := new(c.node, false, rctx, cfg)
	c.node.wildcard = n
}
func (c *NodeCfgr) Static(chunk string, cfg func(*NodeCfgr)) {
	_, exists := c.node.static[chunk]
	if exists {
		c.report.Error(fmt.Sprintf("static \"%s\" node already specified", chunk))
		return
	}
	rctx := c.report.Context(fmt.Sprintf("%%%s", chunk))
	n := new(c.node, false, rctx, cfg)
	c.node.static[chunk] = n
}
func (c *NodeCfgr) Param(name string, cfg func(*NodeCfgr)) {
	if c.node.param != nil {
		c.report.Error(fmt.Sprintf("param \":%s\" node already specified", name))
		return
	}
	rctx := c.report.Context(fmt.Sprintf(":%s", name))
	c.node.param = new(c.node, true, rctx, cfg)
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
			report: c.report.Context("GET"),
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
			report: c.report.Context("POST"),
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
			report: c.report.Context("PUT"),
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
			report: c.report.Context("PATCH"),
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
			report: c.report.Context("DELETE"),
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
			report: c.report.Context("HEAD"),
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
			report: c.report.Context("CONNECT"),
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
			report: c.report.Context("OPTIONS"),
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
			report: c.report.Context("PUT"),
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
