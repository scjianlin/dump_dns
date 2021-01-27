package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
)

func ParserMsg(ch chan gopacket.Packet) {
	// 接受packet
	after := time.After(5 * time.Second)
	for {
		select {
		case <-ch:
			Sem <- 1
			go recvMsg(<-ch) // 异步发送到kafka
		case <-after:
			// timeout
			continue
		}
	}
}

func recvMsg(packet gopacket.Packet) {
	// 转换报文协议
        if (packet == nil) {
           fmt.Println("nil...")
           return
        }
	dns := packet.ApplicationLayer().(*layers.DNS)
	net := packet.NetworkLayer().(*layers.IPv4)

	newTime := time.Now()
	timestamp := newTime.Format(time.RFC3339)

	packetObj := &DnsMessage{}

	// 判断包的方向
	if (int(dns.ANCount) == 0 && int(dns.ResponseCode) > 0) || (int(dns.ANCount) > 0) {
		if int(dns.ANCount) > 0 {
			// 正常回复
			packetObj.ResponseCode = dns.ResponseCode.String()
			packetObj.Name = string(dns.Answers[0].Name)
			packetObj.QueryType = dns.Answers[0].Type.String()
			packetObj.QueryClass = dns.Answers[0].Class.String()
			packetObj.Data = dns.Answers[0].String()
			packetObj.SrcIP = net.SrcIP.String()
			packetObj.DstIP = net.DstIP.String()
			packetObj.Timestamp = timestamp
		} else {
			// 异常回复
			packetObj.ResponseCode = dns.ResponseCode.String()
			packetObj.Name = string(dns.Questions[0].Name)
			packetObj.QueryType = dns.Questions[0].Type.String()
			packetObj.QueryClass = dns.Questions[0].Class.String()
			packetObj.SrcIP = net.SrcIP.String()
			packetObj.DstIP = net.DstIP.String()
			packetObj.Timestamp = timestamp
		}
	}
	msg, err := json.Marshal(packetObj)
	if err != nil {
		fmt.Println("json Marshal dns package error,", err)
		return
	}
        //fmt.Println("发送数据包:",string(msg))
	// 将处理的byte信息发送到另一端.
	SendMsgChan <- msg
}
