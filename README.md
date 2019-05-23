# tinyserver
GO并发服务器框架

Demo

— demoServer.go

— config

​	 — tinyserver.json



demoServe.go

```go
/*
	tinyserver demo
*/

package main

import (
	"fmt"
	"tinyserver/tsinterface"
	"tinyserver/tsnet"
)

//自定义一个路由继承自BaseRouter
type pingRouter struct {
	tsnet.BaseRouter
}

//通过自定义路由进行重写以实现业务方法
func (rp *pingRouter) PreHandler(request tsinterface.IRequest) {
	_, err := request.GetConnection().GetTCPConn().Write([]byte("重写后的PreHandler..."))
	if err != nil {
		fmt.Println(err)
	}
}

func (rp *pingRouter) Handler(request tsinterface.IRequest) {
	_, err := request.GetConnection().GetTCPConn().Write([]byte("重写后的Handler..."))
	if err != nil {
		fmt.Println(err)
	}
}

func (rp *pingRouter) PostHandler(request tsinterface.IRequest) {
	_, err := request.GetConnection().GetTCPConn().Write([]byte("重写后的PostHandler..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//创建一个server
	s := tsnet.NewServer("demo")

	//添加一个自定义路由
	s.AddRouter(&pingRouter{})

	//让server对象启动服务
	s.Serve()
}

```

