package routers

import (
	"net/http"
)

const (
	MAX_RESTFUL_METHODS = 4
)

type Router struct {
	methodTrees 	map[string]*MethodTree
	maxParams  		uint16
}

func NewRouter() *Router {
	router := new(Router)
	router.methodTrees = make(map[string]*MethodTree, MAX_RESTFUL_METHODS)
	return router
}

func (r *Router) GET(path string, handle http.HandlerFunc) {
	r.handle(http.MethodGet, path, handle)
}

func (r *Router) POST(path string, handle http.HandlerFunc) {
	r.handle(http.MethodPost, path, handle)
}

func (r *Router) PUT(path string, handle http.HandlerFunc) {
	r.handle(http.MethodGet, path, handle)
}

func (r *Router) DELETE(path string, handle http.HandlerFunc) {
	r.handle(http.MethodPost, path, handle)
}

func (r *Router) handle(method, path string, handle http.HandlerFunc) {
	if r.methodTrees == nil {
		r.methodTrees = make(map[string]*MethodTree)
	}
	root := r.methodTrees[method]
	if root == nil {
		root = NewMethodTree()
		r.methodTrees[method] = root
	}
	handlers := []http.HandlerFunc{handle}
	root.addRuleToMethodTree(path, handlers)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if root := r.methodTrees[req.Method]; root != nil {
		if handlers := root.getHandlersByPath(path); handlers != nil {
			for _, handler := range handlers {
				handler(w, req)
			}
			return
		}
	}
	// Handle 404
	http.NotFound(w, req)
}