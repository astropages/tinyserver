/*
	连接模块抽象层
*/

package tsinterface

import "net"

//IConnection 连接接口
type IConnection interface {
	Start()                        //开启连接
	Stop()                         //关闭连接
	GetConnID() uint32             //获取连接ID
	GetTCPConn() *net.TCPConn      //获取原生套接字
	GetRemoteAddr() net.Addr       //获取客户端的IP地址
	Send(data []byte, n int) error //发送数据给客户端
}

//HandleFunc 抽象业务处理方法(参数为请求接口)
type HandleFunc func(IRequest) error
