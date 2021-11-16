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
						func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
							w.Write([]byte("welcome"))
							return
						})
					cfg.Static(
						"articles",
						"articles resource.",
						func(cfg *node.NodeCfgr) {
							cfg.GET(
								"all-articles",
								"list of all articles, ordered by publication date",
								func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
									w.Write([]byte("all-articles"))
									return
								})
							cfg.Param(
								"articleID",
								"single article resource",
								func(cfg *node.NodeCfgr) {
									cfg.GET(
										"article",
										"single article full presentation",
										func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
											w.Write([]byte(fmt.Sprintf("article: %s %#v", params["articleID"], params)))
											w.WriteHeader(200)
											return
										})
								})
						})
				})
		})
	fmt.Print(r.String())
	fmt.Println()
	http.ListenAndServe(":8080", d)
}
