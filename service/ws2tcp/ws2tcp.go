package ws2tcp

import (
	"io"
	"log"
	"miniFTL/util"
	"miniFTL/util/config"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{}

func Start(c *config.Config) {
	mux := http.NewServeMux()
	for k, v := range c.Server.Servicemap {
		mux.HandleFunc(c.Server.Path+"/"+k, getHandler(v))
	}
	if c.Server.TLS.Enabled {
		http.ListenAndServeTLS(c.Server.Listen, c.Server.TLS.Certfile, c.Server.TLS.Keyfile, mux)
	} else {
		http.ListenAndServe(c.Server.Listen, mux)
	}
}

func getHandler(v []string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		wsRWC := util.RWC{C: conn}
		// netConn, err := net.Dial(v[0], v[1])
		// if err != nil {
		// 	log.Println(err)
		// }
		netConn, netType := getNetConn(v)
		if netType == "tcp" {
			tcpConn := netConn.(*net.TCPConn)
			biDirectionalCopy(tcpConn, &wsRWC)
		}
		if netType == "udp" {
			udpConn := netConn.(*net.UDPConn)
			biDirectionalCopy(udpConn, &wsRWC)
		}
		// ch := make(chan bool)
		// go util.CopyWorker(netConn, &ioConn, ch)
		// go util.CopyWorker(&ioConn, netConn, ch)
		// <-ch
		// netConn.Close()
		// ioConn.C.Close()
		// <-ch
	}
}

func biDirectionalCopy(io1, io2 io.ReadWriteCloser) {
	ch := make(chan bool)
	go util.CopyWorker(io1, io2, ch)
	go util.CopyWorker(io2, io1, ch)
	<-ch
	io1.Close()
	io2.Close()
	<-ch
}

func getNetConn(v []string) (interface{}, string) {
	if v[1] == "tcp" {
		raddr, err := net.ResolveTCPAddr("tcp", v[0])
		if err != nil {
			log.Fatal(err)
		}
		conn, err := net.DialTCP("tcp", nil, raddr)
		if err != nil {
			log.Fatal(err)
		}
		return conn, "tcp"
	}
	if v[1] == "udp" {
		raddr, err := net.ResolveUDPAddr("udp", v[0])
		if err != nil {
			log.Fatal(err)
		}
		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			log.Fatal(err)
		}
		return conn, "udp"
	}
	return nil, "error"
}
