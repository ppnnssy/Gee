package gee

import (
	"net/http"
	"strings"
)

type router struct {
	//handlers 存储每种请求方式的 HandlerFunc
	//key中存储的是请求方式+URL，如handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
	handlers map[string]HandlerFunc


	//使用 roots 来存储每种请求方式的Trie 树根节点
	//key中存储的是请求方式，如POST,GET
	roots map[string]*node
}


func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

//解析请求的URL，拼接成一个去掉了/的字符串切片，工具函数
// Only one * is allowed
func parsePattern(pattern string) []string {
	//返回一个去掉了“/”的string切片
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { //如果检测到第一个就是*，说明所有的都能匹配，没必要再往下设置了
				break
			}
		}
	}
	return parts
}

//主要用来给gee包中的addRoute调用的
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method] //尝试取method请求方式的根节点
	if !ok { //没找到
		r.roots[method] = &node{} //设置一个空节点
	}
	r.roots[method].insert(pattern, parts, 0) //将请求路由插入
	r.handlers[key] = handler
}

/*
getRoute 函数中，还解析了:和*两种匹配符的参数，返回一个 map 。
例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}
/static/css/geektutu.css匹配到/static/*filepath，解析结果为{filepath: "css/geektutu.css"}。
 */
//通过URL获取对应路由
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	//没有method对应的路由
	if !ok {
		return nil, nil
	}

	//调用search，查找path指向的叶子结点
	n := root.search(searchParts, 0)

	if n != nil {
		//设置params的返回信息
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' { //如果路径中出现:，比如/p/:lang/go
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {//本part以*开头并且后面还有内容，比如/p/:lang/*go
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

//调用路由方法
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
