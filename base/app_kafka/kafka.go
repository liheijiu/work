package app_kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

var (
	client      sarama.SyncProducer //声明一个全局的连接kafka的生产者client
	logDataChan chan *logData
)

//给外部暴露一个函数，做一个配置传递
type logData struct {
	topic string
	data  string
}

//初始化client
func InitKafka(address string, maxSize int) (err error) {

	//生成消息
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	fmt.Println("init kafka，Please wait....")
	client, err = sarama.NewSyncProducer([]string{address}, config)
	if err != nil {
		return
	}

	//loadDataChan初始化
	logDataChan = make(chan *logData, maxSize)
	//开启后台的G从通道中获取数据发往kafka
	go sendTokafka()

	//消费日志

	return
}

func sendTokafka() {
	for {
		select {
		case ld := <-logDataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送消息
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

//给外部暴露一个函数，该函数只把日志数据发送到一个内部的channel中
func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

//根据topic取到所有的分区
//遍历所有的分区,针对每个分区创建一个对应的分区消费者
//异步从每个分区消费信息
