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
)

func (c *NodeCfgr) Wildcard(t, d string, cfg func(*NodeCfgr)) {
	if c.node.wildcard != nil {
		c.report.Error("* node already specified")
		return
	}
	rctx := c.report.Context(fmt.Sprintf("*%s", t))
	n := new(c.node, t, d, false, rctx, cfg)
	c.node.wildcard = n
}
func (c *NodeCfgr) Static(t, d string, cfg func(*NodeCfgr)) {
	_, exists := c.node.static[t]
	if exists {
		c.report.Error(fmt.Sprintf("static \"%s\" node already specified", t))
		return
	}
	rctx := c.report.Context(fmt.Sprintf("%%%s", t))
	n := new(c.node, t, d, false, rctx, cfg)
	c.node.static[t] = n
}
func (c *NodeCfgr) Param(t, d string, cfg func(*NodeCfgr)) {
	if c.node.param != nil {
		c.report.Error(fmt.Sprintf("param \":%s\" node already specified", t))
		return
	}
	rctx := c.report.Context(fmt.Sprintf(":%s", t))
	n := new(c.node, t, d, true, rctx, cfg)
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
func (c *NodeCfgr) GET(t, d string, h Handler) {
	_, exists := c.node.methods[GET]
	if exists {
		c.report.Error("GET handler already specified")
		return
	}
	c.node.methods[GET] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
func (c *NodeCfgr) POST(t, d string, h Handler) {
	_, exists := c.node.methods[POST]
	if exists {
		c.report.Error("POST handler already specified")
		return
	}
	c.node.methods[POST] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
func (c *NodeCfgr) PUT(t, d string, h Handler) {
	_, exists := c.node.methods[PUT]
	if exists {
		c.report.Error("PUT handler already specified")
		return
	}
	c.node.methods[PUT] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
func (c *NodeCfgr) PATCH(t, d string, h Handler) {
	_, exists := c.node.methods[PATCH]
	if exists {
		c.report.Error("PATCH handler already specified")
		return
	}
	c.node.methods[PATCH] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
func (c *NodeCfgr) DELETE(t, d string, h Handler) {
	_, exists := c.node.methods[DELETE]
	if exists {
		c.report.Error("DELETE handler already specified")
		return
	}
	c.node.methods[DELETE] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
func (c *NodeCfgr) HEAD(t, d string, h Handler) {
	_, exists := c.node.methods[HEAD]
	if exists {
		c.report.Error("HEAD handler already specified")
		return
	}
	c.node.methods[HEAD] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
func (c *NodeCfgr) CONNECT(t, d string, h Handler) {
	_, exists := c.node.methods[CONNECT]
	if exists {
		c.report.Error("CONNECT handler already specified")
		return
	}
	c.node.methods[CONNECT] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
func (c *NodeCfgr) OPTIONS(t, d string, h Handler) {
	_, exists := c.node.methods[OPTIONS]
	if exists {
		c.report.Error("OPTIONS handler already specified")
		return
	}
	c.node.methods[OPTIONS] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
func (c *NodeCfgr) TRACE(t, d string, h Handler) {
	_, exists := c.node.methods[TRACE]
	if exists {
		c.report.Error("TRACE handler already specified")
		return
	}
	c.node.methods[TRACE] = &Method{
		title:       t,
		description: d,
		handler:     h,
	}
}
