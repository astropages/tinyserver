/*
	消息模块应用层
*/

package tsnet

//Message 消息类
type Message struct {
	ID      uint32 //消息ID
	Datalen uint32 //数据长度
	Data    []byte //数据
}

//NewMsgPackage 封装消息对象
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		Datalen: uint32(len(data)),
		Data:    data,
	}
}

//GetMsgID 获取消息ID的接口方法
func (m *Message) GetMsgID() uint32 {
	return m.ID
}

//GetMsgLen 获取消息长度的接口方法
func (m *Message) GetMsgLen() uint32 {
	return m.Datalen
}

//GetMsgData 获取消息的接口方法
func (m *Message) GetMsgData() []byte {
	return m.Data
}

//SetMsgID 设置消息ID的接口方法
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

//SetDataLen 设置消息长度的接口方法
func (m *Message) SetDataLen(len uint32) {
	m.Datalen = len
}

//SetData 设置消息的接口方法
func (m *Message) SetData(data []byte) {
	m.Data = data
}
