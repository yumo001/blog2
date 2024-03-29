package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

// 配置文件的映射结构
type SeverConfig struct {
	ServiceName   string        `mapstructure:"serverName"`
	RpcPort       int           `yaml:"rpc_port"`
	Nacos         NacosConfig   `mapstructure:"nacosConfig"`
	Mysql         MysqlConfig   `yaml:"mysql"`
	Elastic       ElasticConfig `yaml:"elastic"`
	JwtSigningKey string        `yaml:"JwtSigningKey"`
	Consul        Consul        `yaml:"consul"`
	Redis         Redis         `yaml:"redis"`
}

type NacosConfig struct {
	ServerConfig struct {
		IpAddr string `mapstructure:"IpAddr"`
		Port   uint64 `mapstructure:"Port"`
	}
	ClientConfig struct {
		NamespaceId string `mapstructure:"NamespaceId"`
	}
	ConfigParam struct {
		DataId string `mapstructure:"DataId"`
		Group  string `mapstructure:"Group"`
	}
}

type MysqlConfig struct {
	Root     string `yaml:"root"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type ElasticConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Consul struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

var (
	// 配置文件结构体实例
	SevConf       = SeverConfig{}
	ElasticClient *elastic.Client
	MysqlDB       = &gorm.DB{}
	RedisClient   = &redis.Client{}
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
