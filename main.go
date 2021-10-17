package main

import (
	"Gee/gee"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

//一个返回年月日的函数
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.NewEngine()
	r.Use(gee.Logger())
	//设置FuncMap，key是string，value是一个函数
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	//解析templates下的模板文件
	r.LoadHTMLGlob("templates/*")
	//将URL中的/assets替换成文件实际存储路径./static，同时注册了对应的路由方法
	r.Static("/assets", "./static")


	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
