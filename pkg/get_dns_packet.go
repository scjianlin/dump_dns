package pkg

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//定义包过滤器
func filters() string {
	filter := fmt.Sprintf("src port %d", port)
	return filter
}

// 负责从网卡接收网络包, 并把包发送到处理函数
func GetPacket(ch chan gopacket.Packet) {

	// 打开网络接口，抓取在线数据
	handle, err := pcap.OpenLive(deviceName, snapLen, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("pcap open live failed: %v", err)
		return
	}
	// 设置过滤器
	filter := filters()
	if err := handle.SetBPFFilter(filter); err != nil {
		fmt.Printf("set bpf filter failed: %v", err)
		return
	}
	// 关闭设备
	defer handle.Close()

//	fmt.Println("start.....")

	// 使用句柄作为数据包源来处理所有数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true
	for packet := range packetSource.Packets() {
	        if (packet.NetworkLayer() == nil || packet.ApplicationLayer() == nil  || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeUDP || packet.ApplicationLayer().LayerType() != layers.LayerTypeDNS ) {
			//fmt.Println("unexpected packet")
			continue
		}
	//	fmt.Println("send...")
		// 将符合条件的packet发送到处理端
		ch <- packet
	}
}
