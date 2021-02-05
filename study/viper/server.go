package main

import (
	"context"
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"log"
)

type Company struct {
	Name                 string
	MarketCapitalization int64
	EmployeeNum          int64
	Department           []interface{}
	IsOpen               bool
}

type TotalYml struct {
	TimeStamp   string
	Address     string
	Postcode    int64
	CompanyInfo Company
}

func main() {
	v := viper.GetViper()

	//设置配置文件的名字
	v.SetConfigName("bootstrap")

	//添加配置文件所在的路径, 可以绝对路径， 也可以相对路径
	v.AddConfigPath("./viper/")

	//设置配置文件类型
	v.SetConfigType("yaml")

	//读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Fatal("read config error!", err.Error())
	}

	//将配置解析为Struct对象
	var cf TotalYml
	if err := v.Unmarshal(&cf); err != nil {
		log.Fatal("read config error!")
	} else {
		log.Println("yamlConfig: ", cf)
	}

	//将配置文件打印为yaml
	ff := v.AllSettings()
	//将配置文件打印为yaml
	bs, err := yaml.Marshal(ff)
	if err != nil {
		log.Fatal("unable to marshal config to YAML")
	} else {
		log.Println("yamlconfig: ", string(bs))
	}

	//将配置文件打印为json
	bs, err = json.Marshal(ff)
	if err != nil {
		log.Fatal("unable to marshal config to json")
	} else {
		log.Println("jsonconfig: ", string(bs))
	}

	//监听配置文件的修改和变动
	ctx, cancel := context.WithCancel(context.Background())
	v.WatchConfig()
	watch := func(e fsnotify.Event) {
		log.Printf("Config file is changed: %s \n", e.String())
		cancel()
	}
	v.OnConfigChange(watch)
	<-ctx.Done()

	log.Println("xx end xx")
}
