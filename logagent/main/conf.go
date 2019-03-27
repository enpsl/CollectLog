package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

var appConfig *Configs

type Configs struct {
	logLevel  string
	logPath   string
	chansize  int
	kafkaAddr string
	etcdAddr  string
	etcdKey   string
}

func loadConf(logType string, fileName string) (err error) {
	conf, err := config.NewConfig(logType, fileName)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	appConfig = &Configs{}
	appConfig.logLevel = conf.String("logs::log_level")
	if len(appConfig.logLevel) == 0 {
		appConfig.logLevel = "debug"
	}
	appConfig.logPath = conf.String("logs::log_path")
	if len(appConfig.logPath) == 0 {
		appConfig.logPath = "./logs/logagent.log"
	}
	appConfig.kafkaAddr = conf.String("kafka::addr")
	if len(appConfig.kafkaAddr) == 0 {
		fmt.Errorf("kafka address err")
		return
	}
	appConfig.etcdAddr = conf.String("etcd::addr")
	if len(appConfig.etcdAddr) == 0 {
		fmt.Errorf("etcd address err")
		return
	}
	appConfig.etcdKey = conf.String("etcd::configKey")
	if len(appConfig.etcdKey) == 0 {
		appConfig.etcdKey = "/logagent"
	}
	appConfig.chansize, err = conf.Int("logs::chan_size")
	if err != nil {
		appConfig.chansize = 100
	}
	return
}
