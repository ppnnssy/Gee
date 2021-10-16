package gee

import (
	"log"
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
	//定义函数的时候使用RouterGroup，调用的时候既可以使用RouterGroup，也可以使用Engine
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
	//初始化的时候engine.RouterGroup里暂时只分配一个Engine引擎，就是engine自己本身
	engine.RouterGroup = &RouterGroup{engine: engine}
	//同样的groups中也只有初始的engine.RouterGroup
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//创建一个新的Group。一般使用engine调用
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	//这里是指针传递地址，所以可以在接下来的操作中直接修改
	//因为NewEngine()的时候group.engine赋值的是engine本身，所以使用engine调用时，下面的这个engine就是调用者本身
	engine := group.engine
	newGroup := &RouterGroup{
		//如果是最外层的Group，prefix开始应该是空的
		prefix: group.prefix + prefix,
		//父路由组就是调用者的group
		parent: group,
		//这里注意，所有路由组共用了同一个引擎
		engine: engine,
	}

	//给engine中的路由组列表加上新建的路由组
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//添加一个路由方法 主要是给框架中其他函数调用的
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	//如果没有路由组，直接engine调用的话，prefix是空的
	//有路由组，就在前面加上prefix
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	/*
	调用了group.engine.router.addRoute来实现了路由的映射。
	由于Engine从某种意义上继承了RouterGroup的所有属性和方法，因为 (*Engine).engine 是指向自己的。
	这样实现，我们既可以像原来一样添加路由，也可以通过分组添加路由。
	 */
	group.engine.router.addRoute(method, pattern, handler)
}

//添加一个GET方法
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//添加一个POST方法
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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
