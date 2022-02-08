package ws2tcp

import (
	"WSTunnel/util"
	"WSTunnel/util/config"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{}

func Start(c *config.Config) {
	mux := http.NewServeMux()
	log.Println("loading", len(c.Server.Servicemap), "endpoint(s)")
	for k, v := range c.Server.Servicemap {
		mux.HandleFunc(c.Server.Path+"/"+k, getHandler(v))
		log.Println("service", k, v[0], v[1], "loaded")
	}
	if c.Server.TLS.Enabled {
		log.Println("Starting HTTPS server", c.Server.Listen, "with TLS")
		http.ListenAndServeTLS(c.Server.Listen, c.Server.TLS.Certfile, c.Server.TLS.Keyfile, mux)
	} else {
		log.Println("Starting HTTP server", c.Server.Listen, "without TLS")
		http.ListenAndServe(c.Server.Listen, mux)
	}
}

func getHandler(v []string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("endpoint", r.URL.Path, "hit")
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		wsRWC := util.RWC{C: conn}
		netConn, netType := getNetConn(v)
		if netType == "tcp" {
			tcpConn := netConn.(*net.TCPConn)
			biDirectionalCopy(tcpConn, &wsRWC)
		}
		if netType == "udp" {
			udpConn := netConn.(*net.UDPConn)
			biDirectionalCopy(udpConn, &wsRWC)
		}
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
