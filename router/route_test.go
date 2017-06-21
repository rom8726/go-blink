package router

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute_Add(t *testing.T) {
	r := NewRoute()
	hello := NewRoute()
	helloWorld := NewRoute()
	helloParam := NewRoute()
	helloGoodbyeAll := NewRoute()

	r.Add("/hello", hello)
	r.Add("/hello/world", helloWorld)
	r.Add("/hello/:param", helloParam)
	r.Add("/hello/goodbye/*", helloGoodbyeAll)

	assert.Equal(t, hello, r.Children["hello"])
	assert.Equal(t, helloWorld, r.Children["hello"].Children["world"])
	assert.Equal(t, helloParam, r.Children["hello"].Children[":"])
	assert.Equal(t, helloGoodbyeAll, r.Children["hello"].Children["goodbye"].Children["*"])

	assert.Equal(t, "hello", hello.Name)
	assert.Equal(t, "world", helloWorld.Name)
	assert.Equal(t, ":", helloParam.Name)
	assert.Equal(t, "param", helloParam.Param)
	assert.Equal(t, "*", helloGoodbyeAll.Name)
	assert.Equal(t, catchAllParam, helloGoodbyeAll.Param)
}

func TestRoute_Handler__should_add_handlers(t *testing.T) {
	r := NewRoute()
	r.GET("/", dummyHandler)
	r.POST("/", dummyHandler1)
	r.GET("/hello", dummyHandler)
	r.POST("/hello", dummyHandler1)
	r.GET("/hello/world", dummyHandler)
	r.GET("/hello/world/*", dummyHandler)
	r.ALL("/hello/goodbye", dummyHandler)

	assert.Contains(t, r.Handlers, GET)
	assert.Contains(t, r.Handlers, POST)
	assert.Contains(t, r.Children["hello"].Handlers, GET)
	assert.Contains(t, r.Children["hello"].Handlers, POST)
	assert.Contains(t, r.Children["hello"].Children["world"].Handlers, GET)
	assert.Contains(t, r.Children["hello"].Children["world"].Children["*"].Handlers, GET)
	assert.Contains(t, r.Children["hello"].Children["goodbye"].Handlers, ALL)
}

func TestRoute_Resolve__should_resolve_routes(t *testing.T) {
	r := NewRoute()
	root0 := r.makePath("")
	root1 := r.makePath("/")
	hello := r.makePath("/hello")
	helloIndex := r.makePath("/hello/")
	helloWorld := r.makePath("/hello/world")
	helloParam := r.makePath("/hello/:param")
	helloGoodbyAll := r.makePath("/hello/goodbye/*")

	assert.Equal(t, root0, r)
	assert.Equal(t, root1, r)
	assert.Equal(t, hello, r.Children["hello"])
	assert.Equal(t, helloIndex, r.Children["hello"])
	assert.Equal(t, helloWorld, r.Children["hello"].Children["world"])
	assert.Equal(t, helloParam, r.Children["hello"].Children[":"])
	assert.Equal(t, helloGoodbyAll, r.Children["hello"].Children["goodbye"].Children["*"])
	assert.Equal(t, "*", helloGoodbyAll.Name)
	assert.Equal(t, catchAllParam, helloGoodbyAll.Param)
}

func TestRoute_Resolve(t *testing.T) {
	r := NewRoute()
	root := r.makePath("")
	hello := r.makePath("/hello")
	helloWorld := r.makePath("/hello/world")
	helloParam := r.makePath("/hello/:param0")
	helloGoodbye := r.makePath("/hello/:param0/goodbye")
	helloParam1 := r.makePath("/hello/:param0/goodbye/:param1")
	helloParamAll := r.makePath("/hello/:param0/*")

	cases := []struct {
		Path   string
		Routes []*Route
		Params Params
	}{
		{"", []*Route{root}, Params{}},
		{"/", []*Route{root}, Params{}},
		{"/hello", []*Route{root, hello}, Params{}},
		{"/hello/world", []*Route{root, hello, helloWorld}, Params{}},
		{"/hello/some-param", []*Route{root, hello, helloParam}, Params{"param0": "some-param"}},
		{"/hello/123/goodbye", []*Route{root, hello, helloParam, helloGoodbye}, Params{"param0": "123"}},
		{"/hello/123/goodbye/456", []*Route{root, hello, helloParam, helloGoodbye, helloParam1}, Params{"param0": "123", "param1": "456"}},
		{"/hello/123/all/the/world/", []*Route{root, hello, helloParam, helloParamAll}, Params{"param0": "123", "path": "all/the/world/"}},
	}

	for _, c := range cases {
		path := c.Path
		routes, params, err := r.Resolve(path)
		assert.Nil(t, err)
		assert.Equal(t, c.Routes, routes)
		assert.Equal(t, c.Params, params)
	}
}

func dummyHandler(context.Context, *Req, *Resp) error  { return nil }
func dummyHandler1(context.Context, *Req, *Resp) error { return nil }
