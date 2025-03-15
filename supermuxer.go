package supermuxer

import (
	"fmt"
	"net/http"
	"slices"
)

type (
	MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

	router struct {
		mux         *http.ServeMux
		basePath    string
		middlewares []MiddlewareFunc
	}

	Router interface {
		Get(string, http.HandlerFunc) *router
		Post(path string, handler http.HandlerFunc) *router
		Put(path string, handler http.HandlerFunc) *router
		Delete(path string, handler http.HandlerFunc) *router

		AddMiddlewares(middleware ...MiddlewareFunc) *router

		// Group creates a group of routes for a base path without middlewares.
		// The original router is not modified, as Group uses a copy.
		//
		// Returns:
		//   - A reference to the Group router.
		//
		// Example:
		//
		//	superRouter := supermuxer.New(serveMux)
		//	superRouter.AddMiddlewares(middleware1, middleware2)
		//	usersRouter := superRouter.Group("/users")
		//	usersRouter.Get("", handler).Post("/{id}", handler)
		//
		//	# Result: supermuxer configuration to handle the request for the endpoints 'GET /users' and 'POST /users/{id}' with 0 middlewares
		Group(basePath string) *router

		// SubGroup creates a subgroup of routes for a base path that REUSES the middlewares defined in the original router.
		// The original router is not modified, as SubGroup uses a copy.
		//
		// Returns:
		//   - A reference to the SubGroup router.
		//
		// Example:
		//
		//	superRouter := supermuxer.New(serveMux)
		//	superRouter.AddMiddlewares(middleware1, middleware2)
		//	superRouter.SubGroup("/users").Get("", handler).Post("/{id}", handler)
		//
		//	# Result: supermuxer configuration to handle the request for the endpoints 'GET /users' and 'POST /users/{id}'
		//		each wrapped in middleware1 and middleware2
		SubGroup(basePath string) *router
	}
)

func getFullPath(method string, basePath string, endpoint string) string {
	fullPath := fmt.Sprintf("%s %s%s", method, basePath, endpoint)
	return fullPath
}

func handlerWithMiddlewares(handler http.HandlerFunc, middlewares []MiddlewareFunc) http.HandlerFunc {
	if len(middlewares) <= 0 {
		return handler
	}

	next := handler

	for _, middleware := range slices.Backward(middlewares) {
		next = middleware(next)
	}

	return next
}

func setRoute(r *router, method string, path string, handler http.HandlerFunc) *router {
	fullPath := getFullPath(method, r.basePath, path)
	wrappedHandler := handlerWithMiddlewares(handler, r.middlewares)

	r.mux.HandleFunc(fullPath, wrappedHandler)
	return r
}

func (r *router) Group(basePath string) *router {
	rCopy := *r

	rCopy.middlewares = []MiddlewareFunc{}
	rCopy.basePath = basePath

	return &rCopy
}

func (r *router) SubGroup(basePath string) *router {
	rCopy := *r

	rCopy.basePath = fmt.Sprintf("%s%s", rCopy.basePath, basePath)

	return &rCopy
}

func (r *router) Get(path string, handler http.HandlerFunc) *router {
	return setRoute(r, http.MethodGet, path, handler)
}

func (r *router) Post(path string, handler http.HandlerFunc) *router {
	return setRoute(r, http.MethodPost, path, handler)
}

func (r *router) Patch(path string, handler http.HandlerFunc) *router {
	return setRoute(r, http.MethodPatch, path, handler)
}

func (r *router) Put(path string, handler http.HandlerFunc) *router {
	return setRoute(r, http.MethodPut, path, handler)
}

func (r *router) Delete(path string, handler http.HandlerFunc) *router {
	return setRoute(r, http.MethodDelete, path, handler)
}

func (r *router) AddMiddlewares(middlewares ...MiddlewareFunc) *router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func New(mux *http.ServeMux) Router {
	return &router{
		mux:         mux,
		middlewares: []MiddlewareFunc{},
	}
}
