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
				func(root *node.StaticNodeCfgr) {
					root.HandleNotFoundErrorWith(
						func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
							w.WriteHeader(404)
							w.Write([]byte("not found"))
							return
						})
					root.HandleMethodNotAllowedErrorWith(
						func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
							w.WriteHeader(404)
							w.Write([]byte("method not allowed"))
							return
						})
					root.HandleInternalServerErrorWith(
						func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
							w.WriteHeader(404)
							w.Write([]byte("internal server error"))
							return
						})
					root.GET(
						func(cfg *node.MethodCfgr) {
							cfg.Title("welcome")
							cfg.Description("latest content")
							cfg.Handler(
								func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
									w.Write([]byte("welcome"))
									return
								})
						})
					root.Static(
						"articles",
						func(articles *node.StaticNodeCfgr) {
							articles.Title("articles")
							articles.Description("articles resource.")
							articles.GET(
								func(cfg *node.MethodCfgr) {
									cfg.Title("all-articles")
									cfg.Description("list of all articles, ordered by publication date")
									cfg.Handler(
										func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
											w.Write([]byte("all-articles"))
											return
										})
								})
							articles.StringParam(
								"articleID",
								func(article *node.StringParamNodeCfgr) {
									article.Title("articleID")
									article.Description("single article resource")
									// article.Query(
									// 	func(params node.Params) (obj interface{}, err error) {
									// 		id, _ := params.Get("articleID")
									// 		idString, ok := id.(string)
									// 		if !ok {
									// 			err = errors.New("no articleID given")
									// 			return
									// 		}
									// 		obj = map[string]string{
									// 			"id": idString,
									// 		}
									// 		return
									// 	})
									article.GET(
										func(cfg *node.MethodCfgr) {
											cfg.Title("article")
											cfg.Description("single article full presentation")
											cfg.Handler(
												func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
													w.WriteHeader(200)
													articleID, ok := params.Get("articleID")
													if ok {
														w.Write([]byte(fmt.Sprintf("article: %s %#v", articleID, params)))
													} else {
														w.Write([]byte(fmt.Sprintf("article: <no articleID> %#v", params)))
													}
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
