package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type Kafka struct {
	client sarama.SyncProducer
}

func InitKafka(addr string) (err error, kafka *Kafka) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	cli, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	kafka = &Kafka{
		client: cli,
	}
	return
}

func (kafka *Kafka) SendMsgToKafka(data string, topic string) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	_, _, err = kafka.client.SendMessage(msg)
	if err != nil {
		logs.Error("send message failed,", err)
		return
	}
	return nil
}
