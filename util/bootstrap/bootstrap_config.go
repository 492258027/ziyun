package bootstrap

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func init() {
	//先预加载匹配的环境变量
	viper.AutomaticEnv()
	//设置读取的配置文件
	viper.SetConfigName("bootstrap")
	//添加读取的配置文件路径
	viper.AddConfigPath("./bootstrap")
	//windows环境下为%GOPATH，linux环境下为$GOPATH
	viper.AddConfigPath("$GOPATH/src/")
	//设置配置文件类型
	viper.SetConfigType("yaml")

	//读取内容
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	if err := subParse("http", &HttpConfig); err != nil {
		log.Println("Failure to parse Http config", err)
	}

	if err := subParse("rpc", &RpcConfig); err != nil {
		log.Println("Failure to parse rpc config", err)
	}

	if err := subParse("consul", &ConsulConfig); err != nil {
		log.Println("Failure to parse Discover config", err)
	}

	if err := subParse("hystrix", &HystrixConfig); err != nil {
		log.Println("Failure to parse hystrix config", err)
	}

	if err := subParse("gateway", &GatewayConfig); err != nil {
		log.Println("Failure to parse gateway config", err)
	}

	if err := subParse("jwt", &JwtConfig); err != nil {
		log.Println("Failure to parse jwt config", err)
	}

	if err := subParse("websocket", &Wsconfig); err != nil {
		log.Println("Failure to parse websocket config", err)
	}

	if err := subParse("manager", &Mgrconfig); err != nil {
		log.Println("Failure to parse manager config", err)
	}

	if err := subParse("rabbitMq", &Mqconfig); err != nil {
		log.Println("Failure to parse rabbitMq config", err)
	}

	if err := subParse("redis", &Redisconfig); err != nil {
		log.Println("Failure to parse redis config", err)
	}

	if err := subParse("mysql", &Mysqlconfig); err != nil {
		log.Println("Failure to parse mysql config", err)
	}

}

func subParse(key string, value interface{}) error {
	log.Printf("prefix：%v", key)
	sub := viper.Sub(key)
	if sub == nil {
		return errors.New("not found!")
	}
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
