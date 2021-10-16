package gee

import (
	"net/http"
)

// HandlerFunc 定义一个请求路由方法
type HandlerFunc func(*Context)

//实现路由分组
type RouterGroup struct {
	//路由组前缀
	prefix      string
	//应用在该分组上的中间件
	middlewares []HandlerFunc
	//当前组的父辈，用来实现路由组的嵌套
	parent      *RouterGroup
	//路由组可以共享一个Engine引擎，简介访问各种接口
	engine      *Engine
}

// Engine 这个引擎用来实现HTTPServer接口
type Engine struct {
	//Engine作为最顶层的分组，也就是说Engine拥有RouterGroup所有的能力。
	*RouterGroup
	//做个切片存储所有路由组
	groups []*RouterGroup
	//路由控制
	router *router
}

// NewEngine 构造一个Engine，用来初始化。里面有一个路由映射表，存储URl和对应的路由。
//目前是静态路由地址，以后改成动态的
func NewEngine() *Engine {
	engine := &Engine{router: newRouter()}
	//初始化的时候暂时只分配一个Engine引擎，就是自己本身
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//创建一个新的Group
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		//这里注意，所有路由组共用了同一个引擎
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//添加一个路由方法 主要是给框架中其他函数调用的
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, engine)
	//ListenAndServe 方法里面会去调用 handler.ServeHTTP() 方法
	return err
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
