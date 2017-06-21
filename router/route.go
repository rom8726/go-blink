package router

import (
	"fmt"
	"strings"
)

const (
	ALL    = ""
	HEAD   = "HEAD"
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

const (
	paramSegment    = ":"
	catchAllSegment = "*"
	catchAllParam   = "path"
)

// Route is a linked tree of routes with handlers and middleware.
type Route struct {
	Name     string
	Param    string
	Middle   []Middleware
	Handlers map[string]Handler // map[method]Handler
	Children map[string]*Route  // map[pattern]*Route
}

// NewRoute creates a root route.
func NewRoute() *Route {
	return newRoute("", "")
}

func newRoute(name string, param string) *Route {
	return &Route{
		Name:     name,
		Param:    param,
		Handlers: make(map[string]Handler),
		Children: make(map[string]*Route),
	}
}

func (r *Route) ALL(pattern string, h Handler)    { r.Handler(ALL, pattern, h) }
func (r *Route) HEAD(pattern string, h Handler)   { r.Handler(HEAD, pattern, h) }
func (r *Route) GET(pattern string, h Handler)    { r.Handler(GET, pattern, h) }
func (r *Route) POST(pattern string, h Handler)   { r.Handler(POST, pattern, h) }
func (r *Route) PUT(pattern string, h Handler)    { r.Handler(PUT, pattern, h) }
func (r *Route) DELETE(pattern string, h Handler) { r.Handler(DELETE, pattern, h) }

func (r *Route) Add(pattern string, child *Route) {
	if pattern == "" || pattern == "/" {
		panic("router: Cannot add a child with a root pattern, the pattern must contain a segment, i.e. /child")
	}
	if child.Name != "" {
		panic("router: The child has already been added to another route")
	}

	// Get the last segment in path.
	// It is the child name or param.
	i := strings.LastIndex(pattern, "/")
	name, param := parseNameParam(pattern[i+1:])
	pattern = pattern[:i]

	// Resolve the child parent.
	route := r.makePath(pattern)

	// Prevent duplicate children.
	if _, ok := route.Children[name]; ok {
		panic("router: Duplicate child")
	}

	// Add the child and set its name and param.
	route.Children[name] = child
	child.Name = name
	child.Param = param
}

func (r *Route) Handler(method string, p string, h Handler) {
	if h == nil {
		panic("router: Nil handler")
	}

	route := r.makePath(p)
	method = strings.ToUpper(method)

	switch method {
	case ALL:
		if len(route.Handlers) > 0 {
			panic("router: Duplicate handler")
		}

		route.Handlers[method] = h

	default:
		if _, ok := route.Handlers[method]; ok {
			panic("router: Duplicate handler")
		}
		if _, ok := route.Handlers[ALL]; ok {
			panic("router: Duplicate handler")
		}

		route.Handlers[method] = h
	}
}

func (r *Route) Middleware(pattern string, m Middleware) {
	if m == nil {
		panic("router: Nil middleware")
	}

	route := r.makePath(pattern)
	route.Middle = append(r.Middle, m)
}

func (r *Route) Match(method string, path string) ([]Middleware, Handler, Params, error) {
	routes, params, err := r.Resolve(path)
	if err != nil {
		return nil, nil, nil, err
	}

	last := routes[len(routes)-1]
	handler := last.Handlers[method]
	if handler == nil {
		handler = last.Handlers[ALL]
		if handler == nil {
			return nil, nil, nil, ErrMethodNotAllowed
		}
	}

	middleware := []Middleware{}
	for _, route := range routes {
		middleware = append(middleware, route.Middle...)
	}

	return middleware, handler, params, nil
}

func (r *Route) Resolve(path string) ([]*Route, Params, error) {
	if path == "" || path == "/" {
		return []*Route{r}, Params{}, nil
	}

	// Add the first slash if absent.
	// The path must always start with a slash.
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	route := r
	segments := strings.Split(path, "/")[1:]

	routes := []*Route{route}
	params := Params{}

	// Traverse the tree.
	for len(segments) > 0 {
		segment := segments[0]

		// A static segment, i.e. "hello" in /hello.
		child, ok := route.Children[segment]
		if !ok {
			// A param segment, i.e. ":param".
			child, ok = route.Children[paramSegment]
			if !ok {
				// A catch all segment, i.e. "*".
				child, ok = route.Children[catchAllSegment]

				if !ok {
					return nil, nil, ErrNotFound
				}

				routes = append(routes, child)
				params[child.Param] = strings.Join(segments, "/")
				break
			}

			params[child.Param] = segment
		}

		routes = append(routes, child)
		route = child
		segments = segments[1:]
	}

	return routes, params, nil
}

// makePath gets or creates a path and returns its last segment.
func (r *Route) makePath(p string) *Route {
	if p == "" || p == "/" {
		return r
	}

	// Require the starting slash.
	if !strings.HasPrefix(p, "/") {
		panic(fmt.Sprint("router: Pattern must start with a slash"))
	}

	// Trim the ending slash.
	if strings.HasSuffix(p, "/") {
		p = p[:len(p)-1]
	}

	// Traverse the route tree, create absent nodes.
	route := r
	segments := strings.Split(p, "/")[1:]

	for len(segments) > 0 {
		name, param := parseNameParam(segments[0])

		// Create a child when absent.
		child, ok := route.Children[name]
		if !ok {
			child = newRoute(name, param)
			route.Children[name] = child
		}

		// Check that the param name matches the child param name.
		if child.Param != param {
			panic(fmt.Sprintf("router: Positional params must have the same name, previous=%v, current=%v",
				child.Param, param))
		}

		route = child
		segments = segments[1:]
	}
	return route
}

func parseNameParam(segment string) (name string, param string) {
	name = segment
	param = ""

	switch {
	case strings.HasPrefix(name, ":"):
		param = name[1:]
		name = ":"
	case name == catchAllSegment:
		param = catchAllParam
	}

	return
}
