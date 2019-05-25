/*
	连接模块应用层
*/

package tsnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"tinyserver/tsinterface"
)

//Connection 连接类
type Connection struct {
	Conn       *net.TCPConn            //当前连接的原生套接字
	ConnID     uint32                  //连接ID
	isClosed   bool                    //当前连接状态
	MsgHandler tsinterface.IMsgHandler //当前连接绑定消息控制器模块（多路由）
}

//NewConnection 初始化连接对象
func NewConnection(conn *net.TCPConn, connID uint32, handler tsinterface.IMsgHandler) tsinterface.IConnection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
	}
	return c
}

//StartReader 针对连接的读方法
func (c *Connection) StartReader() {
	fmt.Printf("连接%d准备就绪（读）\n", c.ConnID)
	defer fmt.Printf("连接%d已关闭\n", c.ConnID)
	defer c.Stop() //关闭连接

	//接收客户端数据
	for {

		//创建数据包对象
		dp := NewDataPack()

		//读取客户端消息的head部分
		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, head); err != nil { //从Conn读取数据
			fmt.Println("io.ReadFull error", err)
			break
		}
		//根据数据长度读取数据
		msg, err := dp.UnPack(head) //将head部分填满
		if err != nil {
			fmt.Println("dp.UnPack error", err)
			break
		}
		//再次读取
		var data []byte
		if msg.GetMsgLen() > 0 { //如果有数据才读取
			data = make([]byte, msg.GetMsgLen())
			if io.ReadFull(c.Conn, data); err != nil { //从Conn读取数据
				fmt.Println("io.ReadFull error", err)
				break
			}
		}
		msg.SetData(data) //填充数据

		//将消息封装成一个Request类的对象
		req := NewRequest(c, msg)

		//使用自定义路由提供的业务方法
		go c.MsgHandler.DoMsgHandler(req)
	}
}

//Start 开启连接的接口方法
func (c *Connection) Start() {

	//读业务
	go c.StartReader()

	//写业务

}

//Stop 关闭连接的接口方法
func (c *Connection) Stop() {

	//回收工作
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭原生套接字
	c.Conn.Close()
}

//GetConnID 获取连接ID的接口方法
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//GetTCPConn 获取原生套接字的接口方法
func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

//GetRemoteAddr 获取客户端的IP地址的接口方法
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//Send 发送数据给客户端的接口方法
func (c *Connection) Send(msgID uint32, msgData []byte) error {
	if c.isClosed == true {
		return errors.New("Conn is Closed")
	}
	//封装成Message
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, msgData))
	if err != nil {
		fmt.Println("dp.Pack error", err)
		return nil
	}

	//将Message发给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write data error", err)
		return err
	}
	return nil
}