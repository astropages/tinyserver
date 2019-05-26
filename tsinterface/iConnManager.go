/*
	连接管理模块抽象层
*/

package tsinterface

//IConnManager 连接管理接口
type IConnManager interface {
	Add(conn IConnection)                   //添加连接
	Remove(connID uint32)                   //删除连接
	Get(connID uint32) (IConnection, error) //根据连接ID获取连接
	Len() int                               //获取全部连接数
	ClearConn()                             //清空全部连接
}
