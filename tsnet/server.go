/*
	服务器模块应用层
*/

package tsnet

import (
	"fmt"
	"net"
	"tinyserver/tsinterface"
	"tinyserver/utils"
)

//Server 服务器类
type Server struct {
	IPVersion string              //连接类型
	IP        string              //IP地址
	Port      int                 //服务器端口
	Name      string              //服务器名称
	Router    tsinterface.IRouter //路由属性
}

//NewServer 初始化服务器对象
func NewServer(name string) tsinterface.Iserver {
	s := &Server{
		IPVersion: "tcp4",
		IP:        utils.GloalObject.Host,
		Port:      utils.GloalObject.Port,
		Name:      utils.GloalObject.Name,
		Router:    nil,
	}
	return s
}

//Start 启动服务器的接口方法
func (s *Server) Start() {
	fmt.Printf("服务器%s启动中\n", s.Name)

	//创建服务器socket（TCP address）
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	//监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("服务器%s已启动并开始监听%s:%d\n", s.Name, s.IP, s.Port)

	//生成连接ID
	var cid uint32
	cid = 1

	//开启goroutine处理listener服务(持续监听直到进程结束)
	go func() {
		//启动server连接服务
		for {
			//循环阻塞等待客户端发送连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println(err)
				continue
			}

			//创建一个Connection对象（将原生Conn和路由绑定）
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//连接建立
			go dealConn.Start()
		}
	}()

}

//Stop 停止服务器的接口方法
func (s *Server) Stop() {

}

//Serve 运行服务器的接口方法
func (s *Server) Serve() {
	//启动server监听功能
	s.Start()

	//阻塞保证main函数不退出（CPU不再处理，节省cpu资源）
	select {}
}

//AddRouter 添加路由的接口方法
func (s *Server) AddRouter(router tsinterface.IRouter) {
	s.Router = router
}
