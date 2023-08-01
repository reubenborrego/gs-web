package web

import (
	"net/http"
)

type DataHandler func(path string, w http.ResponseWriter, r *http.Request) (interface{}, error)
type DataWriter func(path string, data interface{}, w http.ResponseWriter, r *http.Request) error

type Router struct {
	parent   *Router
	matcher  Matcher
	resolver Resolver
	handler  DataHandler
	writer   DataWriter
	hops     []*Router
}

func NewRouter(matcher Matcher, resolver Resolver, handler DataHandler, writer DataWriter) *Router {
	return &Router{
		matcher:  matcher,
		resolver: resolver,
		handler:  handler,
		writer:   writer,
	}
}

func (router *Router) AddHop(hop *Router) *Router {
	router.hops = append(router.hops, hop)
	return hop
}

func (router *Router) AddRelativeHop(hop *Router) *Router {
	//hop.path = filepath.Join(filepath.Dir(router.path), hop.path)
	hop.resolver.pathPrefix(router.resolver.directory())
	hop.parent = router

	return router.AddHop(hop)
}
