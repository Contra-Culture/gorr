package main

import (
	"errors"
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
					root.BeforeDo(
						func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
							w.Write([]byte(fmt.Sprintf("<pre>root before do: %#v</pre>", params)))
							return
						})
					root.InheritableBeforeDo(
						func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
							w.Write([]byte(fmt.Sprintf("<pre>root inheritable before do: %#v</pre>", params)))
							return
						})
					root.AfterDo(
						func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
							fmt.Printf("\nroot after do: %#v", params)
							return
						})
					root.InheritableAfterDo(
						func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
							fmt.Printf("\nroot inheritable after do: %#v", params)
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
							articles.BeforeDo(
								func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
									w.Write([]byte(fmt.Sprintf("<pre>articles before do: %#v</pre>", params)))
									return
								})
							articles.InheritableBeforeDo(
								func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
									w.Write([]byte(fmt.Sprintf("<pre>articles inheritable before do: %#v</pre>", params)))
									return
								})
							articles.AfterDo(
								func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
									fmt.Printf("\narticles after do: %#v", params)
									return
								})
							articles.InheritableAfterDo(
								func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
									fmt.Printf("\narticles inheritable after do: %#v", params)
									return
								})
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
							articles.IDParam(
								"article",
								func(article *node.IDParamNodeCfgr) {
									article.Title("article")
									article.Description("single article resource")
									article.BeforeDo(
										func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
											w.Write([]byte(fmt.Sprintf("<pre>article before do: %#v</pre>", params)))
											return
										})
									article.InheritableBeforeDo(
										func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
											w.Write([]byte(fmt.Sprintf("<pre>article inheritable before do: %#v</pre>", params)))
											return
										})
									article.AfterDo(
										func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
											fmt.Printf("\narticle after do: %#v", params)
											return
										})
									article.InheritableAfterDo(
										func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
											fmt.Printf("\narticle inheritable after do: %#v", params)
											return
										})
									article.Query(
										func(params node.Params) (obj interface{}, err error) {
											id, _ := params.Get("articleID")
											idString, ok := id.(string)
											if !ok {
												err = errors.New("no articleID given")
												fmt.Printf("\n\nNOT FOUND: %s\n\n", err.Error())
												return
											}
											obj = map[string]string{
												"id": idString,
											}
											return
										})
									article.GET(
										func(cfg *node.MethodCfgr) {
											cfg.Title("article")
											cfg.Description("single article full presentation")
											cfg.Handler(
												func(w http.ResponseWriter, r *http.Request, params node.Params) (err error) {
													articleID, ok := params.Get("articleID")
													if !ok {
														w.Write([]byte(fmt.Sprintf("article: <no articleID> %#v", params)))
														return
													}
													article, ok := params.Get("article")
													if ok {
														w.Write([]byte(fmt.Sprintf("article: %s %#v | %#v", articleID, article, params)))
														w.WriteHeader(200)
														return
													}
													w.Write([]byte(fmt.Sprintf("article: %s [no article!] | %#v", articleID, params)))
													err = errors.New("not found")
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
