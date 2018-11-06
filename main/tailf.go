package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"github.com/axgle/mahonia"
	"github.com/Shopify/sarama"
	"log"
	"strings"
)

type ReadFile struct {
	file    *os.File
	gbkFile *mahonia.Reader
}

func (f *ReadFile) gbkDecode() {
	decoder := mahonia.NewDecoder("gbk")
	f.gbkFile = decoder.NewReader(f.file)
}

func (f *ReadFile) ReadPrint(producer sarama.SyncProducer) {
	var n int
	var err error
	data := make([]byte, 1<<16)
	if runtime.GOOS == "windows" {
		f.gbkDecode()
		n, err = f.gbkFile.Read(data)
	} else {
		n, err = f.file.Read(data)
	}
	switch err {
	case nil:
		var lines int
		out := data
		indexs := make(map[int]int)
		for i, d := range out {
			if d == '\n' {
				lines++
				indexs[lines] = i
			}
		}
		lines += 1

		if lines <= line || line <= 0 {
			fmt.Print(string(data[:n]))
			msg := &sarama.ProducerMessage{}
			msg.Topic = "example"
			msg.Partition = int32(-1)
			msg.Key = sarama.StringEncoder("key")
			msg.Value = sarama.StringEncoder(string(data[: n]))

			if err != nil {
				logger.Println("Failed to produce message: %s", err)
				os.Exit(500)
			}
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				logger.Println("Failed to produce message: ", err)
			}
			logger.Printf("partition=%d, offset=%d\n", partition, offset)

		} else {
			index := indexs[lines-line]
			fmt.Print(string(data[index+1 : n]))
			msg := &sarama.ProducerMessage{}
			msg.Topic = "example"
			msg.Partition = int32(-1)
			msg.Key = sarama.StringEncoder("key")
			msg.Value = sarama.StringEncoder(string(data[index+1 : n]))


			if err != nil {
				logger.Println("Failed to produce message: %s", err)
				os.Exit(500)
			}
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				logger.Println("Failed to produce message: ", err)
			}
			logger.Printf("partition=%d, offset=%d\n", partition, offset)
		}

	case io.EOF:
	default:
		fmt.Println(err)
		return
	}
}

var (
	follow bool
	line   int
)

func init() {
	flag.BoolVar(&follow, "f", false, "即时输出文件变化后追加的数据。")
	flag.IntVar(&line, "n", 10, "output the last K lines, instead of the last 10")
}

var (
	logger = log.New(os.Stderr, "[srama]", log.LstdFlags)
)

func main() {
	flag.Parse()
	var err error
	var readFile ReadFile
	os.Args=[]string{"tail","-f","/Users/didi/1.log"}
	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	readFile.file, err = os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer readFile.file.Close()

	sarama.Logger = logger

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	producer, err := sarama.NewSyncProducer(strings.Split("localhost:9092", ","),nil)// config)
	defer producer.Close()
	for {
		readFile.ReadPrint(producer)
		if !follow {
			//break
		}
	}
}