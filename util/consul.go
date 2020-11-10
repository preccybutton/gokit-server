package util

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"time"

	"log"
)

func init(){
	config  := consulapi.DefaultConfig()
	config.Address = "192.168.29.128:8500"
	client, err := consulapi.NewClient(config)
	if err != nil{
		log.Fatal(err)
	}
	ConsulClient = client
	ServiceID = "userservice" + time.Now().String()
}

var (
	ConsulClient *consulapi.Client
	ServiceID string
	ServiceName string
	ServicePort int
)

func SetServiceNameAndPort(name string, port int){
	ServiceName = name
	ServicePort = port
}

func RegService() {
	var err error
	config := consulapi.DefaultConfig()
	config.Address = "10.61.1.232:8500"

	//创建服务注册struct
	reg := consulapi.AgentServiceRegistration{}
	reg.ID = ServiceID	//服务的id
	reg.Name = ServiceName	//服务的name
	reg.Address = "10.61.2.202"		//服务的ip地址
	reg.Port = ServicePort		//服务端口
	reg.Tags = []string{"primary"}

	//设置健康检查
	check := consulapi.AgentServiceCheck{}
	check.Interval = "5s"
	check.HTTP = "http://10.61.2.202:8080/health"
	check.HTTP = fmt.Sprintf("http://%s:%d/health", reg.Address, ServicePort)

	reg.Check = &check
	ConsulClient, err = consulapi.NewClient(config)
	if err != nil{
		log.Fatal(err)
	}
	err = ConsulClient.Agent().ServiceRegister(&reg)
	if err != nil{
		log.Fatal(err)
	}

}

func Unregservice() {
	ConsulClient.Agent().ServiceDeregister(ServiceID)

}