package gee

import (
	"fmt"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/a/c", nil)
	r.addRoute("GET", "/hello/a/b", nil)
	r.addRoute("GET", "/hello/b/a", nil)
	r.addRoute("GET", "/hello/b/b", nil)
	r.addRoute("GET", "/hello/c/a", nil)
	r.addRoute("GET", "/hello/c/b", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestGetRoute(t *testing.T) {
	r:=newTestRouter()
	a:=parsePattern("/hello/b/a")
	fmt.Println(a)
	b:=r.roots["GET"].search(a,0)
	fmt.Println(b)
}
