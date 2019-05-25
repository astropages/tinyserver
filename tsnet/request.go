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
	msg  tsinterface.IMessage    //客户发送的消息
}

//NewRequest 初始化请求对象
func NewRequest(conn tsinterface.IConnection, msg tsinterface.IMessage) tsinterface.IRequest {
	req := &Request{
		conn: conn,
		msg:  msg,
	}
	return req
}

//GetConnection 获取连接对象的接口方法
func (r *Request) GetConnection() tsinterface.IConnection {
	return r.conn
}

//GetMsg 获取客户端消息的接口方法
func (r *Request) GetMsg() tsinterface.IMessage {
	return r.msg
}
