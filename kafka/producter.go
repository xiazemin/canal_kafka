package main

/*import (
	"github.com/Shopify/sarama"
	"fmt"
)

func main() {
	//设置配置
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机的分区类型
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_11_0_0

	//使用配置,新建一个异步生产者//,"IP:9092","IP:9092"
	producer, e := sarama.NewAsyncProducer([]string{"127.0.0.1:9092"}, config)
	if e != nil {
	    panic(e)
	}
	defer producer.AsyncClose()

	//发送的消息,主题,key
	msg := &sarama.ProducerMessage{
	Topic: "example",
	Key:   sarama.StringEncoder("test"),
	}

	var value string
	for {
		value = "this is a message"
		//设置发送的真正内容
		fmt.Scanln(&value)
		//将字符串转化为字节数组
		msg.Value = sarama.ByteEncoder(value)
		fmt.Println(value)

		//使用通道发送
		producer.Input() <- msg

		//循环判断哪个通道发送过来数据.
		select {
		case suc := <-producer.Successes():
		fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
		case fail := <-producer.Errors():
		fmt.Println("err: ", fail.Err)
	}
	}
}
*/

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strings"
)

var (
	logger = log.New(os.Stderr, "[srama]", log.LstdFlags)
)

func main() {
	sarama.Logger = logger

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	msg := &sarama.ProducerMessage{}
	msg.Topic = "example"
	msg.Partition = int32(-1)
	msg.Key = sarama.StringEncoder("key")
	msg.Value = sarama.ByteEncoder("你好, 世界!")

	producer, err := sarama.NewSyncProducer(strings.Split("localhost:9092", ","),nil)// config)
	if err != nil {
		logger.Println("Failed to produce message: %s", err)
		os.Exit(500)
	}
	defer producer.Close()

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		logger.Println("Failed to produce message: ", err)
	}
	logger.Printf("partition=%d, offset=%d\n", partition, offset)
}