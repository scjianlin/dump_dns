package main

import (
	"fmt"
	"dns/pkg"
//	"net/http"
//	_ "net/http/pprof"
)

func main() {
//	go func() {
//		fmt.Println(http.ListenAndServe("0.0.0.0:8080", nil))
//	}()
	go pkg.GetPacket(pkg.RecvMsgChan)
	go pkg.ParserMsg(pkg.RecvMsgChan)
	go pkg.SubmitMsgToMQ()

	defer func() {
		if err := pkg.Producer.Close(); err != nil {
			fmt.Println("Failed to close producer: %s", err)
		}
	}()
	select {}
}
