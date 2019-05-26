/*
	消息路由抽象层
*/

package tsinterface

//IMsgHandler 消息路由接口
type IMsgHandler interface {
	AddRouter(msgID uint32, router IRouter) //添加路由到map集合中（不同消息ID对应不同路由）
	DoMsgHandler(request IRequest)          //路由调度：根据消息ID选择其对应的路由
	StartWorkerPool()                       //启动Woker工作池
	SendMsgTOTaskQueue(IRequest)            //将消息添加到消息队列
}
