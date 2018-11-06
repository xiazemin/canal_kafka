package main

import (
	"fmt"
	"strings"
	"github.com/Shopify/sarama"
)

func main() {
	/*config := sarama.NewConfig()
	config.Version = sarama.V0_10_0_0
	client, err := sarama.NewClient([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		panic("client create error")
	}
	defer client.Close()
	//获取主题的名称集合
	topics, err := client.Topics()
	if err != nil {
		panic("get topics err")
	}
	for _, e := range topics {
		fmt.Println(e)
	}
	//获取broker集合
	brokers := client.Brokers()
	//输出每个机器的地址
	for _, broker := range brokers {
		fmt.Println(broker.Addr())
	}
	*/
	fmt.Printf("metadata test\n")

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_2

	client, err := sarama.NewClient(strings.Split("localhost:9092", ","),nil)//[]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Printf("metadata_test try create client err :%s\n", err.Error())
	return
	}

	defer client.Close()

	// get topic set
	topics, err := client.Topics()
	if err != nil {
	fmt.Printf("try get topics err %s\n", err.Error())
	return
	}

	fmt.Printf("topics(%d):\n", len(topics))

	for _, topic := range topics {
	fmt.Println(topic)
	}

	// get broker set
	brokers := client.Brokers()
	fmt.Printf("broker set(%d):\n", len(brokers))
	for _, broker := range brokers {
	fmt.Printf("%s\n", broker.Addr())
}
}
