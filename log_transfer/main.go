package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("filepath get err", err)
		panic("load conf err")
	}
	//加载配置项
	var filename = dir + "/conf/log_transfer.conf"
	err = initConfig("ini", filename)
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(logConfig)

	err = initLogger(logConfig.LogPath, logConfig.LogLevel)
	if err != nil {
		panic(err)
		return
	}

	logs.Debug("init logger succ")

	err = initKafka(logConfig.KafkaAddr, logConfig.KafkaTopic)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}

	logs.Debug("init kafka succ")

	err = initES(logConfig.ESAddr)
	if err != nil {
		logs.Error("init es failed, err:%v", err)
		return
	}

	logs.Debug("init es client succ")

	err = run()
	if err != nil {
		logs.Error("run  failed, err:%v", err)
		return
	}

	logs.Warn("warning, log_transfer is exited")

}
