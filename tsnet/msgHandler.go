/*
	消息路由应用层
*/

package tsnet

import (
	"fmt"
	"tinyserver/tsinterface"
	"tinyserver/utils"
)

//MsgHandler 消息路由类
type MsgHandler struct {
	Apis             map[uint32]tsinterface.IRouter //存放消息ID和其对应路由的map集合
	TaskQueue        []chan tsinterface.IRequest    //消息队列：一个消息队列对应一个Worker
	WorkerPoolSize   uint32                         //Woker工作池的Worker数量
	MaxWorkerTaskLen uint32                         //每个Worker处理的最大任务数量
}

//NewMsgHandler 初始化消息路由对象
func NewMsgHandler() tsinterface.IMsgHandler {
	return &MsgHandler{
		Apis:             make(map[uint32]tsinterface.IRouter), //给map开辟空间
		TaskQueue:        make([]chan tsinterface.IRequest, utils.GloalObject.WorkerPoolSize),
		WorkerPoolSize:   utils.GloalObject.WorkerPoolSize,
		MaxWorkerTaskLen: utils.GloalObject.MaxWorkerTaskLen,
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
		fmt.Printf("消息%d路由不存在\n", request.GetMsg().GetMsgID())
		return
	}
	//根据ID调用其对应路由的方法
	router.PreHandler(request)
	router.Handler(request)
	router.PostHandler(request)
}

//StarOneWorker 3: Worker工作业务
func (mh *MsgHandler) StarOneWorker(workerID int, taskQueue chan tsinterface.IRequest) {
	fmt.Printf("工作池：Worker %d 已开始工作\n", workerID)

	//循环阻塞：等待接收来自消息队列管道的IRequest消息
	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandler(req) //使用路由调度方法处理业务
		}
	}

}

//StartWorkerPool 1: 启动工作池的接口方法（服务器启动时开启）
func (mh *MsgHandler) StartWorkerPool() {
	fmt.Println("工作池启动中")

	//根据配置的工作池的Worker数量分别开辟管道空间（给每个worker对应的消息队列都开辟一个管道）
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan tsinterface.IRequest)
		//为每一个Worker开启一个goroutine进行工作
		go mh.StarOneWorker(i, mh.TaskQueue[i])
	}
}

//SendMsgToTaskQueue 2: 将消息添加到消息队列的接口方法
func (mh *MsgHandler) SendMsgToTaskQueue(request tsinterface.IRequest) {
	//将Worker和request一一对应绑定来进行任务平均分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//将request发送到绑定了Worker的消息队列（ConnID为10,11,12的消息队列分别对应Worker0,Worker1,Worker2...）
	mh.TaskQueue[workerID] <- request
}
