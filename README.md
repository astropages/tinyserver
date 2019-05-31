# tinyserver
GO并发服务器框架

Demo：

- demoServer.go

- conf
  - tinyserver.json





*demoServer.go*

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

//客户端消息ID为1时对应的路由:

//自定义一个路由继承自BaseRouter
type hiRouter struct {
	tsnet.BaseRouter
}

//通过自定义路由进行重写以实现业务方法
func (rp *hiRouter) Handler(request tsinterface.IRequest) {
	err := request.GetConnection().Send(100, []byte("重写后的Handler... Hi~"))
	if err != nil {
		fmt.Println(err)
	}
}

//客户端消息ID为2时对应的路由:

//自定义一个路由继承自BaseRouter
type helloRouter struct {
	tsnet.BaseRouter
}

//通过自定义路由进行重写以实现业务方法
func (rp *helloRouter) Handler(request tsinterface.IRequest) {
	err := request.GetConnection().Send(200, []byte("重写后的Handler... Hello~"))
	if err != nil {
		fmt.Println(err)
	}
}

//Hook函数：

//welcome ...
func welcome(conn tsinterface.IConnection) {
	fmt.Printf("连接%d即将连接\n", conn.GetConnID())
	if err := conn.Send(201, []byte(fmt.Sprintf("欢迎连接%d", conn.GetConnID()))); err != nil {
		fmt.Println(err)
	}

	conn.SetProperty("Name", "tsc1") //给当前连接模块设置属性
}

//bye ...
func bye(conn tsinterface.IConnection) {
	fmt.Printf("连接%d即将断开\n", conn.GetConnID())

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Printf("%s退出了\n", name)
	}

}

func main() {
	//创建一个server
	s := tsnet.NewServer("demo")

	//注册建立连接之后的Hook函数
	s.AddOnConnStart(welcome)
	//注册断开连接之前的Hook函数
	s.AddOnConnStop(bye)

	//注册自定义路由：客户端发送不同的消息，服务器根据消息处理不同的业务
	s.AddRouter(1, &hiRouter{})    //发送消息ID为1时对应的路由
	s.AddRouter(2, &helloRouter{}) //发送消息ID为2时对应的路由

	//让server对象启动服务
	s.Serve()
}

```



*conf/tinyserver.json*

```json
{
  "Host": "127.0.0.1",
	"Port": 9999,
	"Name": "tinyserver",
	"Version": "0.1",
	"MaxPackageSize": 512,
	"WorkerPoolSize": 10,
	"MaxWorkerTaskLen": 1000,
	"MaxConn": 2
}
```



------



*demoClient.go*

```go
package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"tinyserver/tsnet"
)

func main() {

	fmt.Println("客户端启动中")
	//3秒之后发起请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("Dial error", err)
		return
	}

	//读写测试
	for {

		dp := tsnet.NewDataPack()

		binaryMsg, err := dp.Pack(tsnet.NewMsgPackage(2, []byte("test")))
		if err != nil {
			fmt.Println("dp.Pack error", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("conn.Write error", err)
			return
		}

		//解析服务器回发的数据包
		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, head); err != nil {
			fmt.Println("io.ReadFull( error", err)
			return
		}

		//通过拆包方法把head部分数据填充到Datalen和ID属性中
		msgHead, err := dp.UnPack(head)

		//判断数据长度后读取数据
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*tsnet.Message) //通过接口类型断言向下转换为Message对象
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("io.ReadFull( error", err)
				return
			}
			//打印服务器回发的数据
			fmt.Printf("服务器数据%d: %s (长度为%d)\n", msg.ID, string(msg.Data), msg.Datalen)
		}

		time.Sleep(1 * time.Second)
	}
}

```

