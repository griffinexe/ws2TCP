package ws2tcp

import (
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

func getHandler(v string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		ioConn := util.RWC{C: conn}
		netConn, err := net.Dial("tcp", v)
		if err != nil {
			log.Println(err)
		}
		ch := make(chan bool)
		go util.CopyWorker(netConn, &ioConn, ch)
		go util.CopyWorker(&ioConn, netConn, ch)
		<-ch
		netConn.Close()
		ioConn.C.Close()
		<-ch
	}
}
