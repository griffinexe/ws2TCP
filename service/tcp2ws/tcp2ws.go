package tcp2ws

import (
	"WSTunnel/util"
	"WSTunnel/util/config"
	"log"
	"net"

	"github.com/gorilla/websocket"
)

func Start(c *config.Config) {
	forever := make(chan bool)
	log.Println("creating", len(c.Client.Listenmap), "listener(s)")
	for k, v := range c.Client.Listenmap {
		go handleEndpoint(k, v, c.Client.Upstream)
		log.Println("listener", k, v[0], v[1], "created")
	}
	<-forever
}

func handleEndpoint(serviceName string, listenNet []string, upstream string) {
	log.Println("endpoint", serviceName, "hit")
	netListen, netType := getNetListen(listenNet)
	if netType == "tcp" {
		listen := netListen.(*net.TCPListener)
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println(err)
			}
			go handleConn(conn, upstream+"/"+serviceName)
		}
	}
	if netType == "udp" {
		listen := netListen.(*net.UDPConn)
		for {
			handleConn(listen, upstream+"/"+serviceName)
		}
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

func getNetListen(v []string) (interface{}, string) {
	if v[1] == "tcp" {
		laddr, err := net.ResolveTCPAddr("tcp", v[0])
		if err != nil {
			log.Fatal(err)
		}
		listen, err := net.ListenTCP("tcp", laddr)
		if err != nil {
			log.Fatal(err)
		}
		return listen, "tcp"
	}
	if v[1] == "udp" {
		laddr, err := net.ResolveUDPAddr("udp", v[0])
		if err != nil {
			log.Fatal(err)
		}
		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			log.Fatal(err)
		}
		return conn, "udp"
	}
	return nil, "error"
}
