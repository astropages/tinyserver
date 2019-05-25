/*
	数据包模块抽象层
	解决TCP粘包问题的
	将客户端传输的二进制流按照约定的协议格式进行打包发送（|datalen|dataID|data|），让服务器能够根据该协议对数据进行拆解并正确读取
	服务器将依据照协议按字节大小进行多次读写，先按指定字节读取head部分的datalen和dataID，然后再根据datalen的字节长度值读取data
*/

package tsinterface

//IDataPack 数据包接口
type IDataPack interface {
	GetHeadLen() uint32                //获取二进制包的头部长度（返回固定值8） //Datalen uint32（4字节) + ID uint32（4字节)
	Pack(msg IMessage) ([]byte, error) //封包方法（将 Message 打包成 |datalen|dataID|data|）
	UnPack([]byte) (IMessage, error)   //拆包方法（这里主要将Conn读取的head部分数据进行拆解填充到Message中，返回包含head属性的IMessage后再通过数据长度值进行二次读取Conn剩余的数据）
}
