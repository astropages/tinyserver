/*
	服务器模块抽象层
*/

package tsinterface

//Iserver 服务器接口
type Iserver interface {
	Start()                                 //启动服务器的方法
	Stop()                                  //停止服务器的方法
	Serve()                                 //运行服务器的方法
	AddRouter(msgID uint32, router IRouter) //添加路由的方法

}
