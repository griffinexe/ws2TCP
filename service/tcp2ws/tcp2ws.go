package tcp2ws

import (
	"WSTunnel/util"
	"WSTunnel/util/config"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func Start(c *config.Config) {
	forever := make(chan bool)
	log.Println("creating", len(c.Client.Listenmap), "listener(s)")
	for k, v := range c.Client.Listenmap {
		go handleEndpoint(k, v, c.Client.Upstream, c.Client.ACL[k])
		log.Println("listener", k, v[0], v[1], "created")
	}
	<-forever
}

func handleEndpoint(serviceName string, listenNet []string, upstream string, acl string) {
	netListen, netType := getNetListen(listenNet)
	if netType == "tcp" {
		listen := netListen.(*net.TCPListener)
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println(err)
			}
			log.Println("endpoint", serviceName, "hit")
			go handleConn(conn, upstream+"/"+serviceName, acl)
		}
	}
	if netType == "udp" {
		// listen := netListen.(*net.UDPConn)
		// for {
		// 	handleConn(listen, upstream+"/"+serviceName)
		// }
		log.Println("(fusion2)skipping UDP listener")
		return
	}
}

func handleConn(netConn net.Conn, url string, passwd string) {
	header := http.Header{}
	header.Add("SEC-WSTUNNEL-PARAMS", passwd)
	wsConn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Println(err)
	}
	upstream := util.RWC{C: wsConn}
	util.IOCopy(netConn, &upstream)
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
		// laddr, err := net.ResolveUDPAddr("udp", v[0])
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// conn, err := net.ListenUDP("udp", laddr)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// return conn, "udp"
		log.Println("(fusion2)UDP not supported")
		return nil, "udp"
	}
	return nil, "error"
}
