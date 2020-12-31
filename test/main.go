package main

import (
	"net/http"

	"github.com/Contra-Culture/gorr"
)

func main() {
	rr, err := gorr.New(func(r *gorr.RouterProxy) {
		// doc-gen
		r.OnError(gorr.NotFoundError, respondWithNotFoundError)
		r.OnError(gorr.MethodNotAllowedError, respondWithMethodNotAllowed)
		r.OnError(gorr.InternalServerError, respondWithInternalServerError)
		r.Before(func(w http.ResponseWriter, r *http.Request) {})
		r.BeforeMethod(func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
		r.AfterMethod(func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
		r.After(func(w http.ResponseWriter, r *http.Request) {})
		r.Root("root", "root descr", func(n *gorr.NodeProxy) {
			n.Method("root", "responds with URL", gorr.GET, respondWithURL)
			n.Node(gorr.Static("articles"), "articles resource", func(n *gorr.NodeProxy) {
				n.Method("get-articles", "provides articles", gorr.GET, respondWithArticles)
				n.Method("create-article", "creates article", gorr.POST, createArticle)
				n.Node(gorr.Parameter("article-slug", func(slug string) bool { return slug == "my-article" }), "single article resource by slug", func(n *gorr.NodeProxy) {
					n.Method("get-article", "provides article by its slug", gorr.GET, respondWithArticle)
				})
			})
		})
	})
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", rr)
}

func respondWithURL(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte(r.URL.String()))
	w.WriteHeader(http.StatusFound)
}
func respondWithInternalServerError(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte(r.URL.String()))
	w.WriteHeader(http.StatusInternalServerError)
}
func respondWithMethodNotAllowed(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte(r.URL.String()))
	w.WriteHeader(http.StatusMethodNotAllowed)
}
func respondWithNotFoundError(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte(r.URL.String()))
	w.WriteHeader(http.StatusNotFound)
}
func respondWithArticles(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte(r.URL.String()))
	w.WriteHeader(http.StatusNotFound)
}
func respondWithArticle(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte(r.URL.String()))
	w.WriteHeader(http.StatusNotFound)
}
func createArticle(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte(r.URL.String()))
	w.WriteHeader(http.StatusNotFound)
}
