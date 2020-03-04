package consul

import (
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"

	"fmt"
	"strconv"
)

var count int64

// consul 服务端会自己发送请求，来进行健康检查
//func consulCheck(w http.ResponseWriter, r *http.Request) {
//
//	s := "consulCheck" + fmt.Sprint(count) + "remote:" + r.RemoteAddr + " " + r.URL.String()
//	fmt.Println(s)
//	fmt.Fprintln(w, s)
//	count++
//}

func registerServer() {
	config := consulapi.DefaultConfig()
	config.Address = beego.AppConfig.String("consul::consul_ip") + ":" + beego.AppConfig.String("consul::consul_port")
	client, err := consulapi.NewClient(config)
	if err != nil {
		//log.Log.Fatal("consul client error : ", err)
	}
	port, _ := strconv.Atoi(beego.AppConfig.String("consul::httpport"))
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = beego.AppConfig.String("consul::service_id")     // 服务节点的名称
	registration.Name = beego.AppConfig.String("consul::service_name") // 服务名称
	registration.Port = port                                           // 服务端口
	registration.Address = beego.AppConfig.String("consul::httpip")    // 服务
	registration.Check = &consulapi.AgentServiceCheck{                 // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务，注销时间，相当于过期时间
		// GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		//log.Log.Fatal("register server error : ", err)
	}
}

var MicroClient *consulapi.Client
var GetUserApi string
var GetWalletApi string
var GetAuthApi string

// 从consul中发现服务
func init() {

	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = beego.AppConfig.String("consul::consul_ip") + ":" + beego.AppConfig.String("consul::consul_port")
	var err error
	MicroClient, err = consulapi.NewClient(config)
	if err != nil {
		//log.Log.Fatal("consul client error : ", err)
	}
	GetUserApi = GetService(beego.AppConfig.String("api::user_tfor"), "http://192.168.8.126")
	GetAuthApi = GetService(beego.AppConfig.String("api::auth_tfor"), "http://192.168.8.126")
	GetWalletApi = GetService(beego.AppConfig.String("api::wallet_tfor"), "http://192.168.8.126")
	//consulDeRegister()
}

func GetService(serviceid, defau string) string {
	service, _, err := MicroClient.Agent().Service(serviceid, nil)
	if err == nil {
		return "http://" + service.Address + ":" + strconv.Itoa(service.Port)
	}
	return defau
}

// 取消consul注册的服务
func consulDeRegister() {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = beego.AppConfig.String("consul::consul_ip") + ":" + beego.AppConfig.String("consul::consul_port")
	client, err := consulapi.NewClient(config)
	if err != nil {
		//log.Fatal("consul client error : ", err)
	}

	client.Agent().ServiceDeregister("111")
}
