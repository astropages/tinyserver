/*
	连接模块应用层
*/

package tsnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"tinyserver/tsinterface"
	"tinyserver/utils"
)

//Connection 连接类
type Connection struct {
	Server     tsinterface.IServer     //当前连接所属服务器
	Conn       *net.TCPConn            //当前连接的原生套接字
	ConnID     uint32                  //连接ID
	isClosed   bool                    //当前连接状态
	MsgHandler tsinterface.IMsgHandler //当前连接绑定消息控制器模块（多路由）
	msgChan    chan []byte             //消息通信管道
	exitChan   chan bool               //退出通信管道

	property     map[string]interface{} //当前连接模块的属性集合
	propertyLock sync.RWMutex           //属性锁
}

//NewConnection 初始化连接对象
func NewConnection(server tsinterface.IServer, conn *net.TCPConn, connID uint32, handler tsinterface.IMsgHandler) tsinterface.IConnection {
	c := &Connection{
		Server:     server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
		msgChan:    make(chan []byte),
		exitChan:   make(chan bool),
		property:   make(map[string]interface{}),
	}
	c.Server.GetConnMgr().Add(c) //将连接添加到该服务器的连接管理中

	return c
}

//StartReader 针对连接的读方法
func (c *Connection) StartReader() {
	fmt.Printf("连接%d准备就绪（读）\n", c.ConnID)
	defer fmt.Printf("连接%d已关闭（读）\n", c.ConnID)
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
			if _, err := io.ReadFull(c.Conn, data); err != nil { //从Conn读取数据
				fmt.Println("io.ReadFull error", err)
				break
			}
		}
		msg.SetData(data) //填充数据

		//将消息封装成一个Request类的对象
		req := NewRequest(c, msg)

		//如果开启了工作池：将Request对象交给工作池来处理
		if utils.GloalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req) //发送到消息队列
		} else {
			//如果没开启工作池：
			go c.MsgHandler.DoMsgHandler(req) //使用自定义路由提供的业务方法
		}

	}
}

//StartWriter 给客户端写消息的方法
func (c *Connection) StartWriter() {
	fmt.Printf("连接%d准备就绪（写）\n", c.ConnID)
	defer fmt.Printf("连接%d已关闭（写）\n", c.ConnID)

	//IO多路复用
	for {
		select {
		case data := <-c.msgChan: //判断是否有数据传入
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("c.Conn.Write error", err)
				return
			}

		case <-c.exitChan: //监听退出通信管道
			return
		}
	}
}

//Start 开启连接的接口方法
func (c *Connection) Start() {

	//读业务
	go c.StartReader()

	//写业务
	go c.StartWriter()

	//
	c.Server.CallOnConnStar(c)
}

//Stop 关闭连接的接口方法
func (c *Connection) Stop() {

	//调用销毁断开前用户自定义的Hook函数
	c.Server.CallOnConnStop(c)

	//回收工作
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//向Writer发出退出信息
	c.exitChan <- true

	//关闭原生套接字
	c.Conn.Close()

	//将连接从服务器连接管理模块删除
	c.Server.GetConnMgr().Remove(c.ConnID)

	//释放channel资源
	close(c.msgChan)
	close(c.exitChan)
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

	//将Message通过管道发送给Writer
	c.msgChan <- binaryMsg

	return nil
}

//SetProperty 设置当前连接模块属性的接口方法
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//GetProperty 获取当前连接模块属性的接口方法
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("No this property" + key)
}

//RemoveProperty 删除当前连接模块属性的接口方法
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	c.propertyLock.Unlock()

	delete(c.property, key)
}
