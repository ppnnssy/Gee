package gee

import (
	"fmt"
	"strings"
)

//路由树的节点
type node struct {
	pattern  string // 待匹配路由，例如 /p/:lang
	part     string // 路由中的一部分，一般是第一个匹配到的子节点，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	//为了实现动态路由，加上判断是够精准匹配的布尔值
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
}

// 用于匹配并返回子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}
// 用于匹配子节点，返回的子节点放进切片中
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		//因为插入的时候不会在一个子节点的children切片中放相同的part，所以一定只有一个匹配到的
		//也就是说nodes的长度一定是1
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}


//插入子节点
func (n *node) insert(pattern string, parts []string, height int) {
	//递归结束条件，在最末子节点上才会设置pattern。前面的节点的pattern都是空的
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	//找找看n的子节点有没有part
	child := n.matchChild(part)
	if child == nil { //没找到
		//把part作为子节点插入进去，如果出现：和*，就把isWild设置为ture
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}

		//通过这个操作修改n的子节点，因为传入的是指针，所以实参n被修改了
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)

}


//查找路由方法，返回最末叶子结点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //如果查询到叶子节点，或者子节点是*开头，结束递归
		if n.pattern == "" {//子节点没有设置pattern
			return nil
		}
		return n
	}
	part := parts[height]
	//查询所有匹配的子节点
	children := n.matchChildren(part)
	for k,v:=range children{
		fmt.Printf("children[%d]:%v\n",k,v)
	}

	//对每一个匹配的子节点都执行一个查询
	//返回第一个匹配成功的
	for _, child := range children {
		//fmt.Println("child:",child)
		//如果匹配正确，最终的返回值是最末的叶子节点
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
