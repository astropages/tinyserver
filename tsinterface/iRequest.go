/*
	请求模块抽象层
*/

package tsinterface

//IRequest 请求接口
type IRequest interface {
	GetConnection() IConnection //获取连接模块对象
	GetData() []byte            //获取连接请求的数据
	GetDataLen() int            //获取连接请求的数据长度
}
