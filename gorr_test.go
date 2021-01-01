package gorr_test

import (
	"net/http"
	"net/url"

	. "github.com/Contra-Culture/gorr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dumbHandler(w http.ResponseWriter, r *http.Request, ps map[string]string) {}

var _ = Describe("gorr", func() {
	Describe("router DSL", func() {
		var (
			router *Router
			err    error
		)
		Context("when configured not properly", func() {
			It("fails", func() {
				// root node is required
				router, err = New(func(r *RouterProxy) {})
				Expect(err).To(MatchError("root node not specified"))
				Expect(router).To(BeNil())

				// ... and can be specified only one
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
				})
				Expect(err).To(MatchError("root node already specified"))
				Expect(router).To(BeNil())

				// "not found" route handler is required!
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
				})
				Expect(err).To(MatchError("`Not Found` handler not specified"))
				Expect(router).To(BeNil())

				// ... and can be specified only once
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(NotFoundError, dumbHandler)
				})
				Expect(err).To(MatchError("`Not Found` handler already specified"))
				Expect(router).To(BeNil())

				// "method not allowed" handler is required!
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
					r.OnError(NotFoundError, dumbHandler)
				})
				Expect(err).To(MatchError("`Method Not Allowed` handler not specified"))
				Expect(router).To(BeNil())

				// ... and can be specified only once
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
				})
				Expect(err).To(MatchError("`Method Not Allowed` handler already specified"))
				Expect(router).To(BeNil())

				// "internal server error" handler is required!
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
				})
				Expect(err).To(MatchError("`Internal Server Error` handler not specified"))
				Expect(router).To(BeNil())

				// ... and can be specified only once
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("`Internal Server Error` handler already specified"))
				Expect(router).To(BeNil())
			})
		})
		Context("when configured properly", func() {
			It("returns router", func() {
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
							w.Write([]byte("ok"))
							w.WriteHeader(http.StatusFound)
						})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(router).NotTo(BeNil())
			})
		})
	})

	Describe("Node", func() {
		Context("when configured not properly", func() {
			var (
				router *Router
				err    error
			)
			It("fails", func() {
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("node should have at least one method/handler or a child node"))
				Expect(router).To(BeNil())

				// only one handler for GET HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", GET, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("GET handler already specified"))
				Expect(router).To(BeNil())

				// only one handler for HEAD HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", HEAD, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", HEAD, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("HEAD handler already specified"))
				Expect(router).To(BeNil())

				// only one handler for POST HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", POST, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", POST, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("POST handler already specified"))
				Expect(router).To(BeNil())

				// only one handler for PUT HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returs sitemap", PUT, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", PUT, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("PUT handler already specified"))
				Expect(router).To(BeNil())

				// only one handler for DELETE HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", DELETE, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", DELETE, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("DELETE handler already specified"))
				Expect(router).To(BeNil())

				// only one handler for CONNECT HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", CONNECT, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", CONNECT, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("CONNECT handler already specified"))
				Expect(router).To(BeNil())

				// only one handler for OPTIONS HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", OPTIONS, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", OPTIONS, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("OPTIONS handler already specified"))
				Expect(router).To(BeNil())

				// only one handler for TRACE HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", TRACE, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", TRACE, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("TRACE handler already specified"))
				Expect(router).To(BeNil())

				// only one handler for PATCH HTTP method can be specified
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Method("root", "returns sitemap", PATCH, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
						n.Method("root", "returns sitemap", PATCH, func(w http.ResponseWriter, r *http.Request, ps map[string]string) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("PATCH handler already specified"))
				Expect(router).To(BeNil())

				// child nodes should have at least one method/andler or child
				router, err = New(func(r *RouterProxy) {
					r.Root("root", "root node", func(n *NodeProxy) {
						n.Node(Static("articles"), "articles resource", func(n *NodeProxy) {})
					})
					r.OnError(NotFoundError, dumbHandler)
					r.OnError(MethodNotAllowedError, dumbHandler)
					r.OnError(InternalServerError, dumbHandler)
				})
				Expect(err).To(MatchError("node should have at least one method/handler or a child node"))
				Expect(router).To(BeNil())
			})
		})

		Describe("chunker", func() {
			var (
				chunker *Chunker
				addr    *url.URL
			)
			Context("creation", func() {
				Describe("NewChunker()", func() {
					It("returns chunker", func() {
						addr, _ = url.ParseRequestURI("/")
						chunker = NewChunker(addr)
						Expect(chunker).NotTo(BeNil())
						Expect(chunker.Chunk()).To(BeEquivalentTo(""))
					})
				})
			})
			Context("usage", func() {
				Describe(".Chunk()", func() {
					It("returns chunk", func() {
						addr, _ = url.ParseRequestURI("/")
						chunker = NewChunker(addr)
						Expect(chunker).NotTo(BeNil())
						Expect(chunker.Chunk()).To(BeEquivalentTo(""))

						addr, _ = url.ParseRequestURI("/articles")
						chunker = NewChunker(addr)
						Expect(chunker).NotTo(BeNil())
						Expect(chunker.Chunk()).To(BeEquivalentTo(""))

						addr, _ = url.ParseRequestURI("/articles/")
						chunker = NewChunker(addr)
						Expect(chunker).NotTo(BeNil())
						Expect(chunker.Chunk()).To(BeEquivalentTo(""))
					})
				})
				Describe(".Next()", func() {
					Context("when has next", func() {
						It("returns true and increments index", func() {
							addr, _ = url.ParseRequestURI("/articles/some-article")
							chunker = NewChunker(addr)
							Expect(chunker).NotTo(BeNil())
							Expect(chunker.Chunk()).To(BeEquivalentTo(""))
							Expect(chunker.Next()).To(BeTrue())
							Expect(chunker.Chunk()).To(BeEquivalentTo("articles"))
							Expect(chunker.Next()).To(BeTrue())
							Expect(chunker.Chunk()).To(BeEquivalentTo("some-article"))
							Expect(chunker.Next()).To(BeFalse())
							Expect(chunker.Chunk()).To(BeEquivalentTo("some-article"))
						})
					})
					Context("when has no next", func() {
						It("returns false", func() {
							addr, _ = url.ParseRequestURI("/")
							chunker = NewChunker(addr)
							Expect(chunker).NotTo(BeNil())
							Expect(chunker.Chunk()).To(BeEquivalentTo(""))
							Expect(chunker.Next()).To(BeFalse())
							Expect(chunker.Chunk()).To(BeEquivalentTo(""))
						})
					})
				})
				Describe(".Set()", func() {
					Context("when new param", func() {
						It("sets param", func() {
							addr, _ = url.ParseRequestURI("/")
							chunker = NewChunker(addr)
							Expect(chunker.Params()).To(HaveLen(2))
							Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
							Expect(chunker.Set("param", "value")).NotTo(HaveOccurred())
							Expect(chunker.Params()).To(HaveLen(3))
							Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["param"]).To(BeEquivalentTo("value"))
							Expect(chunker.Set("param2", "value2")).NotTo(HaveOccurred())
							Expect(chunker.Params()).To(HaveLen(4))
							Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["param"]).To(BeEquivalentTo("value"))
							Expect(chunker.Params()["param2"]).To(BeEquivalentTo("value2"))
						})
					})
					Context("when param already exists", func() {
						It("fails", func() {
							addr, _ = url.ParseRequestURI("/")
							chunker = NewChunker(addr)
							Expect(chunker.Params()).To(HaveLen(2))
							Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
							Expect(chunker.Set("param", "value")).NotTo(HaveOccurred())
							Expect(chunker.Params()).To(HaveLen(3))
							Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["param"]).To(BeEquivalentTo("value"))
							Expect(chunker.Set("param", "value2")).To(MatchError("param already set"))
							Expect(chunker.Params()).To(HaveLen(3))
							Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
							Expect(chunker.Params()["param"]).To(BeEquivalentTo("value"))
						})
					})
				})
				Describe(".Params()", func() {
					It("returns params", func() {
						addr, _ = url.ParseRequestURI("/")
						chunker = NewChunker(addr)
						params := chunker.Params()
						Expect(params).To(HaveLen(2))
						Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
						Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
						Expect(chunker.Set("param", "value")).NotTo(HaveOccurred())
						Expect(params).To(HaveLen(3))
						Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
						Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
						Expect(chunker.Params()["param"]).To(BeEquivalentTo("value"))
						Expect(chunker.Set("param2", "value2")).NotTo(HaveOccurred())
						Expect(params).To(HaveLen(4))
						Expect(chunker.Params()["$url"]).To(BeEquivalentTo("/"))
						Expect(chunker.Params()["$path"]).To(BeEquivalentTo("/"))
						Expect(chunker.Params()["param"]).To(BeEquivalentTo("value"))
						Expect(chunker.Params()["param2"]).To(BeEquivalentTo("value2"))
					})
				})
			})
		})
	})
})
