package pkg

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)


// 提交数据到MQ
func SubmitMsgToMQ() {
	// 收数据发送到MQ
	for {
		resultMsg := <-SendMsgChan
		//msg := producerMessage(Topic, resultMsg)
		msg := producerMessage(Conf.Base.Topic, resultMsg)
		Producer.Input() <- msg

		select {
		case <-Producer.Successes():
			<-Sem // 处理结束, 释放线程
		case err := <-Producer.Errors():
			fmt.Println("Failed to send msg: %s", err)
		}
	}
}

// 实例化MQ对象
func asyncProducer(address []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Net.MaxOpenRequests = 5
	config.Producer.MaxMessageBytes = 1000000
	config.Producer.RequiredAcks = sarama.RequiredAcks(1)
	config.Producer.Timeout = 10 * time.Second
	config.Producer.Flush.Frequency = 0
	config.Producer.Flush.Bytes = 0
	config.Producer.Flush.Messages = 0
	config.Producer.Flush.MaxMessages = 0
	config.ChannelBufferSize = 256
	config.Version, _ = sarama.ParseKafkaVersion("0.8.2.0")
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		fmt.Println("Failed to create producer: %s", err)
	}
	return producer

}

//
func producerMessage(topic string, msg []byte) *sarama.ProducerMessage {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}
	return message
}
