/*
	基础路由模块抽象层
*/

package tsinterface

//IRouter 基础路由接口
type IRouter interface {
	PreHandler(request IRequest)  //处理业务之前的方法
	Handler(request IRequest)     //处理业务的方法
	PostHandler(request IRequest) ////处理业务之后的方法
}
