package pkg

import "fmt"

func Start(env string) {
	initConf(env)	// init config
	go GetPacket(RecvMsgChan)
	go ParserMsg(RecvMsgChan)
	go SubmitMsgToMQ()
	defer func() {
		if err := Producer.Close(); err != nil {
			fmt.Println("Failed to close producer: %s", err)
		}
	}()
	select {}
}
