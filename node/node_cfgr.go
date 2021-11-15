package node

import "github.com/Contra-Culture/report"

type (
	NodeCfgr struct {
		node   *Node
		report *report.RContext
	}
)

func (c *NodeCfgr) Wildcard() {

}
func (c *NodeCfgr) Static() {

}
func (c *NodeCfgr) Parameter() {

}
func (c *NodeCfgr) HandleNotFoundError() {

}
func (c *NodeCfgr) HandleMethodNotAllowedError() {

}
func (c *NodeCfgr) HandleInternalServerError() {

}
func (c *NodeCfgr) DoAfter() {

}
func (c *NodeCfgr) DoBefore() {

}
func (c *NodeCfgr) GET() {

}
func (c *NodeCfgr) POST() {

}
func (c *NodeCfgr) PUT() {

}
func (c *NodeCfgr) PATCH() {

}
func (c *NodeCfgr) DELETE() {

}
func (c *NodeCfgr) HEAD() {

}
func (c *NodeCfgr) CONNECT() {

}
func (c *NodeCfgr) OPTIONS() {

}
func (c *NodeCfgr) TRACE() {

}
