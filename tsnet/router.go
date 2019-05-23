/*
	路由模块应用层
*/

package tsnet

import (
	"tinyserver/tsinterface"
)

//BaseRouter 路由类
type BaseRouter struct {
}

//PreHandler 处理业务之前的接口方法
func (r *BaseRouter) PreHandler(request tsinterface.IRequest) {
	//框架使用者须通过添加自定义路由进行方法重写
}

//Handler 处理业务的接口方法
func (r *BaseRouter) Handler(request tsinterface.IRequest) {
	//框架使用者须通过添加自定义路由进行方法重写
}

//PostHandler 处理业务之后的接口方法
func (r *BaseRouter) PostHandler(request tsinterface.IRequest) {
	//框架使用者须通过添加自定义路由进行方法重写
}
