package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 给常用的map取个别名，用起来更简洁
type H map[string]interface{}

//把w和r封装起来。
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
}

//初始化
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
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


//一个备用函数，输入动态路由，返回本次请求的确切URL（如果有的话）
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
