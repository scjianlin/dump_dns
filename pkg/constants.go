package pkg

import (
	"github.com/Shopify/sarama"
	"github.com/google/gopacket"
)

const (
	deviceName = "eth0"		//网卡设备
	snapLen    = int32(65535)
	port       = uint16(53)  // 监听数据包端口
	Topic      = "dns-query-log"  // topic
	MaxReqs    = 200 // goroutine 并发数量
)

var (
	RecvMsgChan = make(chan gopacket.Packet)
	SendMsgChan = make(chan []byte)
	//Brokers     = []string{"10.16.244.29:9092"}
	Brokers     = []string{"10.248.224.155:9092"}	// Kafka消息队列地址
	Producer    sarama.AsyncProducer
	Sem         = make(chan int, MaxReqs) // 限制goroutine并发数量
)

type DnsMessage struct {
	Timestamp    string `json:"timestamp" time_format:"2019-01-02 00:00:00"`
	SrcIP        string `json:"src_ip,omitempty"`
	DstIP        string `json:"dst_ip,omitempty"`
	QueryType    string `json:"queryType,omitempty"`
	QueryClass   string `json:"queryClass,omitempty"`
	Data         string `json:"data,omitempty"`
	Name         string `json:"name,omitempty"`
	ResponseCode string `json:"responseCode,omitempty"`
}
