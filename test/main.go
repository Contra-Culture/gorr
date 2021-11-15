package main

import (
	"fmt"
	"net/http"

	"github.com/Contra-Culture/gorr"
	"github.com/Contra-Culture/gorr/node"
)

func main() {
	d, r := gorr.New(
		func(cfg *gorr.DispatcherCfgr) {
			cfg.Root(
				"Test root.",
				func(cfg *node.NodeCfgr) {
					cfg.GET(
						"welcome",
						"latest content",
						func(w http.ResponseWriter, r *http.Request, params map[string]string) {
							w.Write([]byte("welcome"))
						})
					cfg.Static(
						"articles",
						"articles resource.",
						func(cfg *node.NodeCfgr) {
							cfg.GET(
								"all-articles",
								"list of all articles, ordered by publication date",
								func(w http.ResponseWriter, r *http.Request, params map[string]string) {
									w.Write([]byte("all-articles"))
								})
						})
				})
		})
	fmt.Print(r.String())
	fmt.Println()
	http.ListenAndServe(":8080", d)
}
