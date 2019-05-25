/*
	数据包模块应用层
	解决TCP的粘包问题
	将客户端传输的数据按照约定的协议格式进行打包发送（|datalen|dataID|data|），让服务器能够根据该协议对数据进行拆解并正确读取
	服务器将依据照协议按字节大小进行多次读写，先按指定字节读取head部分的datalen和dataID，然后再根据datalen的字节长度值读取data
*/

package tsnet

import (
	"bytes"
	"encoding/binary"
	"tinyserver/tsinterface"
)

//DataPack 数据包类
type DataPack struct {
}

//NewDataPack 初始化数据包对象
func NewDataPack() *DataPack {
	return &DataPack{}
}

//GetHeadLen 获取二进制包头部长度的接口方法（返回固定值8）
func (d *DataPack) GetHeadLen() uint32 {
	return 8 //Datalen uint32（4字节) + ID uint32（4字节)
}

//Pack 封包的接口方法（将 Message 打包成 |datalen|dataID|data|）
func (d *DataPack) Pack(msg tsinterface.IMessage) ([]byte, error) {
	//创建一个存放二进制数据的缓冲区（一个空字节切片）
	dataBuffer := bytes.NewBuffer([]byte{})
	//将消息长度Datalen以二进制方式写到buffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将消息ID以二进制方式写到buffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	//将消息内容Data以二进制方式写到buffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}

	//返回一个buffer（二进制数据流）
	return dataBuffer.Bytes(), nil
}

//UnPack 拆包的接口方法（将|datalen|dataID|data| 拆解到 Message 结构体中）
func (d *DataPack) UnPack(binaeyData []byte) (tsinterface.IMessage, error) {
	//解包需分两次（获取head部分（Datalen和ID），获取data部分（data））
	//第一次读取固定的长度8字节，第二次是根据Datalen再次读取data

	msg := &Message{} //msgHead.Datalen, msgHead.dataID

	//创建一个读取二进制数据流的阅读器
	buff := bytes.NewReader(binaeyData)

	//读取二进数据制流：先读取Datalen到msg的DataLen属性中
	if err := binary.Read(buff, binary.LittleEndian, &msg.Datalen); err != nil {
		return nil, err
	}

	//读取二进数据制流：再读取ID到msg的ID属性中
	if err := binary.Read(buff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	//返回一个IMessage
	return msg, nil
}
