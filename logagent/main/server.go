package main

import (
	"github.com/astaxie/beego/logs"
	"logagent/kafka"
	"logagent/tailf"
	"time"
)

func serverRun(kafka *kafka.Kafka) (err error) {

	for {
		msg := tailf.GetOneLine()
		err = kafka.SendMsgToKafka(msg.Msg, msg.Topic)
		if err != nil {
			logs.Error("send to kafka failed err:%+v", err)
			time.Sleep(100 * time.Second)
		}
		//time.Sleep(time.Second)
	}
	return
}
