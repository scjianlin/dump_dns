package pkg

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/google/gopacket"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	snapLen    = int32(65535)
	port       = uint16(53)  // 监听数据包端口
	MaxReqs    = 200 // goroutine 并发数量
)

var (
	RecvMsgChan = make(chan gopacket.Packet)
	SendMsgChan = make(chan []byte)
	Producer    sarama.AsyncProducer
	Conf		Config	// 配置文件
	Sem         = make(chan int,  MaxReqs) // 限制goroutine并发数量
)

// DNS报文
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

// 配置文件
type Config struct {
	Base   Base 	`yaml:"base"`
}

type Base struct {
	DeviceName string	`yaml:"deviceName"`
	Brokers    string	`yaml:"brokers"`
	Topic      string	`yaml:"topic"`
}

func initConf(env string)  {
	content, _ := ioutil.ReadFile(fmt.Sprintf("conf/%s.yaml",env))
	err := yaml.Unmarshal(content,&Conf)
	if err != nil {
		fmt.Println("get server config error,", err)
		return
	}
	Producer = asyncProducer([]string{Conf.Base.Brokers})
}