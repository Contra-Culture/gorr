package main

import (
	"fmt"
	"net/http"

	"github.com/Contra-Culture/gorr"
)

func main() {
	rr, err := gorr.New(func(r *gorr.RouterProxy) {
		// doc-gen
		r.OnError(gorr.NotFoundError, respondWithNotFoundError)
		r.OnError(gorr.MethodNotAllowedError, respondWithMethodNotAllowed)
		r.OnError(gorr.InternalServerError, respondWithInternalServerError)
		r.Before(func(w http.ResponseWriter, r *http.Request) { fmt.Printf("\n\t-> %s %s", r.Method, r.URL.String()) })
		r.BeforeMethod(func(w http.ResponseWriter, r *http.Request, ps map[string]string) { w.Write([]byte("beforehook\n")) })
		r.AfterMethod(func(w http.ResponseWriter, r *http.Request, ps map[string]string) { fmt.Print("\nafterhook\n") })
		r.After(func(w http.ResponseWriter, r *http.Request) {})
		r.Root("root", "root descr", func(n *gorr.NodeProxy) {
			n.Method("root", "responds with URL", gorr.GET, respondWithURL)
			n.Node(gorr.Static("articles"), "articles resource", func(n *gorr.NodeProxy) {
				n.Method("get-articles", "provides articles", gorr.GET, respondWithArticles)
				n.Method("create-article", "creates article", gorr.POST, createArticle)
				n.Node(gorr.Parameter("article-slug", func(slug string) bool { return slug == "my-article" }), "single article resource by slug", func(n *gorr.NodeProxy) {
					n.Method("get-article", "provides article by its slug", gorr.GET, respondWithArticle)
					n.Method("update-article", "updates article", gorr.PATCH, respondWithArticle)
					n.Node(gorr.Static("edit"), "edits article", func(n *gorr.NodeProxy) {
						n.Method("edit-article", "edits article", gorr.GET, respondWithArticle)
					})
				})
			})
		})
	})
	if err != nil {
		fmt.Printf("\n\nerror: %s\n\n", err.Error())
		panic(err)
	}
	http.ListenAndServe(":8080", rr)
}

func respondWithURL(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	var respbs = []byte(fmt.Sprintf("{\"a\":\"%s\",\"m\":\"%s\",\"p\":\"%#v\"}", r.Method, r.URL.String(), ps))
	w.Write(respbs)
	w.WriteHeader(http.StatusOK)
}
func respondWithInternalServerError(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte("{\"error\":\"internal server error\"}"))
	w.WriteHeader(http.StatusInternalServerError)
}
func respondWithMethodNotAllowed(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte("{\"error\":\"method not allowed\"}"))
	w.WriteHeader(http.StatusMethodNotAllowed)
}
func respondWithNotFoundError(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte("{\"error\":\"not found\"}"))
	w.WriteHeader(http.StatusNotFound)
}
func respondWithArticles(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	var respbs = []byte(fmt.Sprintf("{\"a\":\"%s\",\"m\":\"%s\",\"p\":\"%#v\"}", r.Method, r.URL.String(), ps))
	w.Write(respbs)
	w.WriteHeader(http.StatusOK)
}
func respondWithArticle(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	var respbs = []byte(fmt.Sprintf("{\"a\":\"%s\",\"m\":\"%s\",\"p\":\"%#v\"}", r.Method, r.URL.String(), ps))
	w.Write(respbs)
	w.WriteHeader(http.StatusFound)
}
func createArticle(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	w.Write([]byte(r.URL.String()))
	w.WriteHeader(http.StatusCreated)
}
