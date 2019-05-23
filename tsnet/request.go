/*
	请求模块应用层
*/

package tsnet

import (
	"tinyserver/tsinterface"
)

//Request 请求类
type Request struct {
	conn tsinterface.IConnection //连接信息
	data []byte                  //数据内容
	len  int                     //数据长度
}

//NewRequest 初始化请求对象
func NewRequest(conn tsinterface.IConnection, data []byte, n int) tsinterface.IRequest {
	req := &Request{
		conn: conn,
		data: data,
		len:  n,
	}
	return req
}

//GetConnection 获取连接对象的接口方法
func (r *Request) GetConnection() tsinterface.IConnection {
	return r.conn
}

//GetData 获取连接请求数据的接口方法
func (r *Request) GetData() []byte {
	return r.data
}

//GetDataLen 获取连接请求数据长度的接口方法
func (r *Request) GetDataLen() int {

	return r.len
}
