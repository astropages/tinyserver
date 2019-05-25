/*
	消息路由应用层
*/

package tsnet

import (
	"fmt"
	"tinyserver/tsinterface"
)

//MsgHandler 消息路由类
type MsgHandler struct {
	Apis map[uint32]tsinterface.IRouter //存放消息ID和其对应路由的map集合
}

//NewMsgHandler 初始化消息路由对象
func NewMsgHandler() tsinterface.IMsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]tsinterface.IRouter), //给map开辟空间
	}
}

//AddRouter 添加路由到map集合中的接口方法
func (mh *MsgHandler) AddRouter(msgID uint32, router tsinterface.IRouter) {
	//判断要添加的消息ID是否存在
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Printf("消息%d已存在\n", msgID)
		return
	}
	//将自定义路由关联到消息ID并添加到map集合中
	mh.Apis[msgID] = router
	fmt.Printf("消息%d路由添加成功\n", msgID)
}

//DoMsgHandler 路由调度的接口方法（根据消息ID选择其对应的路由）
func (mh *MsgHandler) DoMsgHandler(request tsinterface.IRequest) {
	//request获取消息ID，然后获得该ID对应的路由
	router, ok := mh.Apis[request.GetMsg().GetMsgID()]
	if !ok {
		fmt.Printf("消息%d路由不存在", request.GetMsg().GetMsgID())
	}
	//根据ID调用其对应路由的方法
	router.PreHandler(request)
	router.Handler(request)
	router.PostHandler(request)
}
