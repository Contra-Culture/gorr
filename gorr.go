package gorr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Contra-Culture/report"
	"github.com/google/uuid"
)

type (
	Dispatcher struct {
		root *Node
	}
	DispatcherCfgr struct {
		dispatcher *Dispatcher
		report     report.Node
	}
	PathHelper  func(map[string]string) string
	PathHelpers map[string]PathHelper
)

func (c *DispatcherCfgr) Root(d string, cfg func(*StaticNodeCfgr)) {
	if c.dispatcher.root != nil {
		c.report.Error("root node already specified")
		return
	}
	r := report.New("root")
	root := new(nil, STATIC)
	nc := &StaticNodeCfgr{
		NodeCfgr{
			node:   root,
			report: r,
		},
	}
	cfg(nc)
	c.dispatcher.root = root
	c.report = r
}
func New(cfg func(*DispatcherCfgr)) (d *Dispatcher, r report.Node) {
	d = &Dispatcher{}
	r = report.New("dispatcher")
	cfg(
		&DispatcherCfgr{
			dispatcher: d,
			report:     r,
		})
	return
}
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rep := report.New("request <%s>", requestUUID())
	rep.Info("URL: %s", r.URL.String())
	d.root.Handle(rep, w, r)
	fmt.Println(report.ToString(rep))
}
func requestUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
