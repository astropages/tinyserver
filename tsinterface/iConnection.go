/*
	连接模块抽象层
*/

package tsinterface

import "net"

//IConnection 连接接口
type IConnection interface {
	Start()                                  //开启连接
	Stop()                                   //关闭连接
	GetConnID() uint32                       //获取连接ID
	GetTCPConn() *net.TCPConn                //获取原生套接字
	GetRemoteAddr() net.Addr                 //获取客户端的IP地址
	Send(msgID uint32, msgData []byte) error //发送数据给客户端

	SetProperty(key string, value interface{})   //设置当前连接模块属性
	GetProperty(key string) (interface{}, error) //获取当前连接模块属性
	RemoveProperty(key string)                   //删除当前连接模块属性
}

//HandleFunc 抽象业务处理方法(参数为请求接口)
type HandleFunc func(IRequest) error
