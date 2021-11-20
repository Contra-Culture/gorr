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
				func(root *node.NodeCfgr) {
					root.HandleInternalServerErrorWith(
						func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
							w.Write([]byte("not found"))
							w.WriteHeader(404)
							return
						})
					root.HandleMethodNotAllowedErrorWith(
						func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
							w.Write([]byte("method not allowed"))
							w.WriteHeader(404)
							return
						})
					root.HandleInternalServerErrorWith(
						func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
							w.Write([]byte("internal server error"))
							w.WriteHeader(404)
							return
						})
					root.GET(
						func(cfg *node.MethodCfgr) {
							cfg.Title("welcome")
							cfg.Description("latest content")
							cfg.Handler(
								func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
									w.Write([]byte("welcome"))
									return
								})
						})
					root.Static(
						"articles",
						func(articles *node.NodeCfgr) {
							articles.Title("articles")
							articles.Description("articles resource.")
							articles.GET(
								func(cfg *node.MethodCfgr) {
									cfg.Title("all-articles")
									cfg.Description("list of all articles, ordered by publication date")
									cfg.Handler(
										func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
											w.Write([]byte("all-articles"))
											return
										})
								})
							articles.Param(
								"articleID",
								func(article *node.NodeCfgr) {
									article.Title("articleID")
									article.Description("single article resource")
									article.GET(
										func(cfg *node.MethodCfgr) {
											cfg.Title("article")
											cfg.Description("single article full presentation")
											cfg.Handler(
												func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
													w.Write([]byte(fmt.Sprintf("article: %s %#v", params["articleID"], params)))
													w.WriteHeader(200)
													return
												})
										})
								})
						})
				})
		})
	fmt.Print(r.String())
	fmt.Println()
	http.ListenAndServe(":8080", d)
}
