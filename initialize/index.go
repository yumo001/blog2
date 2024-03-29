package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"strconv"

	"github.com/yumo001/blog2/global"
)

func Viper() {
	v := viper.New()
	v.SetConfigFile("./conf/config.yaml")

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("读取yaml配置文件失败", err)
		return
	}
	err = v.Unmarshal(&global.SevConf)
	if err != nil {
		log.Fatal("解析配置失败", err)
		return
	}

	v.WatchConfig()

	v.OnConfigChange(func(c fsnotify.Event) {
		err = v.Unmarshal(&global.SevConf)
		if err != nil {
			log.Fatal("解析配置失败", err)
			return
		}
		log.Println("本地配置文件发生变动")
		Nacos()
	})
	log.Println("viper初始化完成")
}

func Nacos() {

	sc := []constant.ServerConfig{
		{
			IpAddr: global.SevConf.Nacos.ServerConfig.IpAddr,
			Port:   global.SevConf.Nacos.ServerConfig.Port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId: global.SevConf.Nacos.ClientConfig.NamespaceId,
		LogDir:      "tmp/log",
		CacheDir:    "tmp/cache",
		LogLevel:    "debug",
	}

	cci, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig":  cc,
		"serverConfigs": sc,
	})
	if err != nil {
		log.Fatal("nacos配置失败", err)
	}

	configStr, err := cci.GetConfig(vo.ConfigParam{
		DataId: global.SevConf.Nacos.ConfigParam.DataId,
		Group:  global.SevConf.Nacos.ConfigParam.Group,
	})
	if err != nil {
		log.Fatal("nacos获取配置失败", err)
		return
	}
	err = yaml.Unmarshal([]byte(configStr), &global.SevConf)
	if err != nil {
		log.Fatal("nacos解析配置失败", err)
		return
	}
	log.Println(global.SevConf)
	log.Println("nacos初始化完成")
	Mysql()
}

func Consul() {
	conf := api.DefaultConfig()
	conf.Address = global.SevConf.Consul.Host + ":" + strconv.Itoa(global.SevConf.Consul.Port)

	client, err := api.NewClient(conf)
	if err != nil {
		log.Fatal("consul创建实例失败", err)
		return
	}

	asr := api.AgentServiceRegistration{
		ID:      uuid.New().String(),
		Name:    "bolg_grpc",
		Address: GetIp()[2],
		Port:    global.SevConf.RpcPort,
		Check: &api.AgentServiceCheck{
			GRPC:                           GetIp()[2] + ":" + strconv.Itoa(global.SevConf.RpcPort),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	err = client.Agent().ServiceRegister(&asr)
	if err != nil {
		log.Fatal("consul注册服务失败", err)
		return
	}
	log.Println("consul初始化完成")
}

// 封装docker镜像时内部ip
func GetIp() (ip []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("获取当前服务公网ip失败")
		return nil
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = append(ip, ipNet.IP.String())
			}
		}

	}
	log.Println(ip)
	return ip
}

func Mysql() {
	var err error
	global.MysqlDB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", global.SevConf.Mysql.Root, global.SevConf.Mysql.Password, global.SevConf.Mysql.Host, global.SevConf.Mysql.Port, global.SevConf.Mysql.Database)), &gorm.Config{})
	if err != nil {
		zap.S().Panic("数据库连接失败", err)
		return
	}
	log.Println("Mysql初始化完成")
}
