/*
	连接管理模块应用层
*/

package tsnet

import (
	"errors"
	"fmt"
	"sync"
	"tinyserver/tsinterface"
)

//ConnManager 连接管理类
type ConnManager struct {
	connections map[uint32]tsinterface.IConnection //全部连接
	connLock    sync.RWMutex                       //互斥锁
}

//NewConnManager 初始化连接管理对象
func NewConnManager() tsinterface.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]tsinterface.IConnection),
	}
}

//Add 添加连接的接口方法
func (connMgr *ConnManager) Add(conn tsinterface.IConnection) {
	//加锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//添加到map
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Printf("连接%d记录已添加\n", conn.GetConnID())
}

//Remove 删除连接的接口方法
func (connMgr *ConnManager) Remove(connID uint32) {
	//加锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除
	delete(connMgr.connections, connID)

	fmt.Printf("连接%d记录已删除\n", connID)
}

//Get 根据连接ID获取连接的接口方法
func (connMgr *ConnManager) Get(connID uint32) (tsinterface.IConnection, error) {
	//加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	//获取连接
	if conn, ok := connMgr.connections[connID]; !ok {
		return conn, nil

	}
	return nil, errors.New("Connection not found")
}

//Len 获取全部连接数的接口方法
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//ClearConn 清空全部连接的接口方法
func (connMgr *ConnManager) ClearConn() {
	//加锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//遍历清空
	for connID, conn := range connMgr.connections {
		//执行连接关闭方法
		conn.Stop()
		//从map中删除连接
		delete(connMgr.connections, connID)
	}
}
