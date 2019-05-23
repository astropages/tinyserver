package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Global 全局类
type Global struct {
	Host           string //当前监听的IP
	Port           int    //当前监听的端口
	Name           string //当前服务器的名称
	Version        string //当前框架的版本号
	MaxPackageSize uint32 //每次Read的最大长度
}

//GloalObject 定义全局对象
var GloalObject *Global

//LoadConfig 加载文件的方法
func (g *Global) LoadConfig() {
	data, err := ioutil.ReadFile("config/tinyserver.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(data, &GloalObject) //解析json内容到全局变量
	if err != nil {
		log.Fatalln(err)
	}
}

//加载全局变量
func init() {
	GloalObject = &Global{
		//默认值
		Host:           "127.0.0.1",
		Port:           8999,
		Name:           "TinyServer",
		Version:        "0.1",
		MaxPackageSize: 512,
	}

	//加载自定义配置
	GloalObject.LoadConfig()

}
