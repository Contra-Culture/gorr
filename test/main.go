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
						"welcome",
						"latest content",
						func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
							w.Write([]byte("welcome"))
							return
						})
					root.Static(
						"articles",
						func(articles *node.NodeCfgr) {
							articles.Title("articles")
							articles.Description("articles resource.")
							articles.GET(
								"all-articles",
								"list of all articles, ordered by publication date",
								func(w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
									w.Write([]byte("all-articles"))
									return
								})
							articles.Param(
								"articleID",
								func(article *node.NodeCfgr) {
									article.Title("articleID")
									article.Description("single article resource")
									article.GET(
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
