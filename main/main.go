// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"time"
	"log"
	"strings"

	"github.com/CanalClient/canal-go/client"
	protocol "github.com/CanalClient/canal-go/protocol"

	"github.com/golang/protobuf/proto"
	"github.com/Shopify/sarama"

	"github.com/kubernetes/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/json"
)

var (
	logger = log.New(os.Stderr, "[srama]", log.LstdFlags)
	producer sarama.SyncProducer
	err error
)
func sendKafka(producer sarama.SyncProducer,entrys []protocol.Entry) {


	//msg := &sarama.ProducerMessage{}
	//msg.Topic = "example"
	//msg.Partition = int32(-1)
	//msg.Key = sarama.StringEncoder("key")
	//msg.Value = sarama.ByteEncoder("你好, 世界!")

	for _, entry := range entrys {
		if entry.GetEntryType() == protocol.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == protocol.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(protocol.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)
		if rowChange != nil {
			eventType := rowChange.GetEventType()
			header := entry.GetHeader()
			fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))

			for _, rowData := range rowChange.GetRowDatas() {
				msg := &sarama.ProducerMessage{}
				msg.Topic = "example"
				msg.Partition = int32(-1)
				if eventType == protocol.EventType_DELETE {
					js:=columnToJson(rowData.GetBeforeColumns())
					msg.Key = sarama.StringEncoder("delete")
					msg.Value = sarama.ByteEncoder(js)
				} else if eventType == protocol.EventType_INSERT {
					js:=columnToJson(rowData.GetAfterColumns())
					msg.Key = sarama.StringEncoder("insert")
					msg.Value = sarama.ByteEncoder(js)
				} else {
					//fmt.Println("-------> before")
					//js:=columnToJson(rowData.GetBeforeColumns())
					//fmt.Println("-------> after")
					js:=columnToJson(rowData.GetAfterColumns())
					msg.Key = sarama.StringEncoder("other")
					msg.Value = sarama.ByteEncoder(js)
				}
				partition, offset, err := producer.SendMessage(msg)
				if err != nil {
					logger.Println("Failed to produce message: ", err)
				}
				logger.Printf("partition=%d, offset=%d\n", partition, offset)

			}
		}
	}
}

func columnToJson(columns []*protocol.Column) ([]byte){
	var row map[string]string=make(map[string]string)
	for _, col := range columns {
		if(len(col.GetName())>0) {
			row[col.GetName()] = col.GetValue()
		}
	}
	str,err:=json.Marshal(row)
	fmt.Println(string(str),err)
	return str
}

func main() {
	sarama.Logger = logger

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	producer, err = sarama.NewSyncProducer(strings.Split("localhost:9092", ","),nil)// config)
	if err != nil {
		logger.Println("Failed to produce message: %s", err)
		os.Exit(500)
	}
	defer producer.Close()

	connector := client.NewSimpleCanalConnector("127.0.0.1", 11111, "", "", "example", 60000, 60*60*1000)
	//connector := client.NewSimpleCanalConnector("192.168.199.17", 11111, "", "", "example", 60000, 60*60*1000)
	connector.Connect()
	connector.Subscribe(".*\\\\..*")

	for {

		message := connector.Get(100, nil, nil)
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(300 * time.Millisecond)
			fmt.Println("===没有数据了===")
			continue
		}

		printEntry(message.Entries)
		sendKafka(producer,message.Entries)

	}
}

func printEntry(entrys []protocol.Entry) {

	for _, entry := range entrys {
		if entry.GetEntryType() == protocol.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == protocol.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(protocol.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)
		if rowChange != nil {
			eventType := rowChange.GetEventType()
			header := entry.GetHeader()
			fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))

			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == protocol.EventType_DELETE {
					printColumn(rowData.GetBeforeColumns())
				} else if eventType == protocol.EventType_INSERT {
					printColumn(rowData.GetAfterColumns())
				} else {
					fmt.Println("-------> before")
					printColumn(rowData.GetBeforeColumns())
					fmt.Println("-------> after")
					printColumn(rowData.GetAfterColumns())
				}
			}
		}
	}
}

func printColumn(columns []*protocol.Column) {
	for _, col := range columns {
		fmt.Println(fmt.Sprintf("%s : %s  update= %t", col.GetName(), col.GetValue(), col.GetUpdated()))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
