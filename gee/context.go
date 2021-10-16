package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 给常用的map取个别名，用起来更简洁
type H map[string]interface{}

//把w和r封装起来。实际执行程序时所有的handler都会集中到c中，然后执行
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	//这两个是常用的属性，拿出来方便直接访问，其实在req中都有
	Path   string
	Method string
	// 状态码
	StatusCode int

	//用于输出动态路由查询时的确切URL
	Params map[string]string

	//增加两个中间件使用的参数
	//把所有方法（包括中间件函数和待执行的路由方法）保存到handlers中。
	//其中中间件是由RouterGroup.middlewares []HandlerFunc中传过来的
	handlers []HandlerFunc
	//index是记录当前执行到第几个中间件
	index    int
}


//初始化
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index: -1,

	}
}

//执行c.handlers中的所有函数。此时切片中存储了中间件和路由方法，依次执行。
// ServeHTTP（）执行流程中先向handlers添加中间件，所以会先执行中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		//依次执行中间件中的函数
		c.handlers[c.index](c)
	}
}


func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置响应头状态码和context的状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头参数
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//发送一个字符串给客户端
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//发送
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 发送[]byte切片类型的数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//发送html类型数据
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

//失败的方法
func (c *Context) Fail(code int, err string) {
	//如果失败，执行到这里，index的值直接变成最大，终止后续程序执行
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

//一个备用函数，输入动态路由，返回本次请求的确切URL（如果有的话）
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
