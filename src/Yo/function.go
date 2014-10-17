package yo

import (
	"restcontrol"
)

var s *HttpServer

/*
函数名：初始化

返回值：error 错误码

success

complaxinit 重复初始化
*/
func Init() string {
	if s == nil {
		s = new(HttpServer)
		return s.Init()
	}

	return "complaxinit"
}

/*
函数名：增加一个 control

返回值：error 错误码

success

uninit 未初始化

isexist 重复绑定
*/
func AddControl(szPath string, control restcontrol.RESTControl) string {
	if s == nil {
		return "uninit"
	}

	return s.AddControl(szPath, control)
}

/*
函数名：启动服务

返回值：error 错误码

success

uninit 未初始化

httperr golang http 服务内部错误
*/
func StartServer() string {
	if s == nil {
		return "uninit"
	}

	return s.Start(":8080")
}