/*
	消息模块抽象层
*/

package tsinterface

//IMessage 消息接口
type IMessage interface {
	GetMsgID() uint32   //获取消息ID
	GetMsgLen() uint32  //获取数据长度
	GetMsgData() []byte //获取数据
	SetMsgID(uint32)    //设置消息ID
	SetDataLen(uint32)  //设置数据长度
	SetData([]byte)     //设置数据
}
