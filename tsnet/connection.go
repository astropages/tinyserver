/*
	连接模块应用层
*/

package tsnet

import (
	"fmt"
	"io"
	"net"
	"tinyserver/tsinterface"
	"tinyserver/utils"
)

//Connection 连接类
type Connection struct {
	Conn     *net.TCPConn        //当前连接的原生套接字
	ConnID   uint32              //连接ID
	isClosed bool                //当前连接状态
	Router   tsinterface.IRouter //当前连接绑定的路由
}

//NewConnection 初始化连接对象
func NewConnection(conn *net.TCPConn, connID uint32, router tsinterface.IRouter) tsinterface.IConnection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
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
		//读取数据
		buf := make([]byte, utils.GloalObject.MaxPackageSize)
		n, err := c.Conn.Read(buf)
		if n == 0 {
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("Read buf error", err)
			return
		}

		//将数据封装成一个Request类的对象
		req := NewRequest(c, buf, n)

		//使用自定义路由提供的业务方法
		go func() {
			c.Router.PreHandler(req)
			c.Router.Handler(req)
			c.Router.PostHandler(req)
		}()
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
func (c *Connection) Send(data []byte, n int) error {
	if _, err := c.Conn.Write(data[:n]); err != nil {
		fmt.Println("Write data error", err)
		return err
	}
	return nil
}
