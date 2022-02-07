package tcp2ws

import (
	"log"
	"miniFTL/util"
	"miniFTL/util/config"
	"net"

	"github.com/gorilla/websocket"
)

func Start(c *config.Config) {
	forever := make(chan bool)
	for k, v := range c.Client.Listenmap {
		go handleEndpoint(k, v, c.Client.Upstream)
	}
	<-forever
}

func handleEndpoint(serviceName, localListen, upstream string) {
	tcpListen, err := net.Listen("tcp", localListen)
	if err != nil {
		log.Println(err)
	}
	for {
		tcpConn, err := tcpListen.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConn(tcpConn, upstream+"/"+serviceName)
	}
}

func handleConn(netConn net.Conn, url string) {
	wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println(err)
	}
	upstream := util.RWC{C: wsConn}
	ch := make(chan bool)
	go util.CopyWorker(netConn, &upstream, ch)
	go util.CopyWorker(&upstream, netConn, ch)
	<-ch
	netConn.Close()
	upstream.C.Close()
	<-ch
}
