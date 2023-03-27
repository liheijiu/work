package app_kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"work/base/app_conf"
)

var (
	client sarama.SyncProducer //声明一个全局的连接kafka的生产者client
)

//初始化client
func InitKafka(cfg *app_conf.AppConf) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	fmt.Println("init kafka，Please wait....")
	client, err = sarama.NewSyncProducer([]string{cfg.KafkaConf.Address}, config)
	if err != nil {
		return
	}
	return
}

func SendTokafka(topic, data string) {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
