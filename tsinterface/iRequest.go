/*
	请求模块抽象层
*/

package tsinterface

//IRequest 请求接口
type IRequest interface {
	GetConnection() IConnection //获取连接模块对象
	GetMsg() IMessage           //获取客户端消息
}
