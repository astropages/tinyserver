/*
	服务器模块抽象层
*/

package tsinterface

//IServer 服务器接口
type IServer interface {
	Start()                                         //启动服务器的方法
	Stop()                                          //停止服务器的方法
	Serve()                                         //运行服务器的方法
	AddRouter(msgID uint32, router IRouter)         //添加路由的方法
	GetConnMgr() IConnManager                       //获取连接管理模块
	AddOnConnStart(hookFunc func(conn IConnection)) //注册建立连接之后的自动Hook函数
	AddOnConnStop(hookFunc func(conn IConnection))  //注册断开连接之前的自动Hook函数
	CallOnConnStar(conn IConnection)                //调用建立连接之后的自动Hook函数
	CallOnConnStop(conn IConnection)                //调用断开连接之前的自动Hook函数
}
