package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logagent/kafka"
	"logagent/tailf"
	"os"
	"path/filepath"
	"time"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("filepath get err", err)
		panic("load conf err")
	}
	//加载配置项
	var filename = dir + "/conf/logagent.conf"
	err = loadConf("ini", filename)
	if err != nil {
		fmt.Println("loadConf err", err)
	}
	//fmt.Println(*appConfig)
	//加载日志
	err = initLogger()
	if err != nil {
		fmt.Println("initLogger err", err)
	}
	etcdClient, err := initEtcd()
	if err != nil {
		fmt.Println("initEtcd err", err)
	}
	collectConf, errCollect := etcdClient.getEtcdConf(appConfig.etcdKey)
	if errCollect != nil {
		logs.Error("getEtcdConf err", err)
	}
	err = etcdClient.initEtcdWatcher()
	if errCollect != nil {
		logs.Error("initEtcdWatcher err", err)
	}
	logs.Debug("init etcd success")
	err = tailf.InitTailf(collectConf, appConfig.chansize)
	if err != nil {
		logs.Error("InitTailf err", err)
	}
	logs.Debug("init tailf success")
	err, kafkaClient := kafka.InitKafka(appConfig.kafkaAddr)
	if err != nil {
		fmt.Printf("init kafa error%+v\n", err)
		return
	}
	logs.Debug("init kafka success")
	go func() {
		var count int
		for {
			count++
			logs.Debug("test for logger %d", count)
			time.Sleep(time.Second * 1)
		}
	}()
	err = serverRun(kafkaClient)
	if err != nil {
		logs.Error("serverRun err", err)
	}

	logs.Info("program finished")
}
