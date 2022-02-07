// package main

// import (
// 	"encoding/json"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"

// 	"github.com/gorilla/websocket"
// 	// "golang.org/x/net/websocket"
// )

// // TODO: fix flags

// // var tcpAddr string
// var addrMap map[string]string

// var wsUp = websocket.Upgrader{}

// type cfg struct {
// 	Listen     string            `json:"listen"`
// 	Path       string            `json:"path"`
// 	ServiceMap map[string]string `json:"servicemap"`
// 	Tls        struct {
// 		Enabled  bool   `json:"enabled"`
// 		KeyFile  string `json:"keyfile"`
// 		CertFile string `json:"certfile"`
// 	} `json:"tls"`
// }

// func main() {
// 	// // parse flags
// 	// var addr string
// 	// var path string
// 	// flag.StringVar(&addr, "addr", "127.0.0.1:8090", "define address to listen")
// 	// flag.StringVar(&path, "path", "/ws", "define path of websocket")
// 	// flag.Usage = usage
// 	// flag.Parse()
// 	// tcpAddr = flag.Arg(0)
// 	// if tcpAddr == "" {
// 	// 	log.Fatal("no TCP address provided")
// 	// }

// 	//read config
// 	var ws2tcpCfg cfg
// 	f, err := os.Open("./config.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	b, _ := ioutil.ReadAll(f)
// 	err = json.Unmarshal(b, &ws2tcpCfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	addrMap = ws2tcpCfg.ServiceMap

// 	// start web server
// 	mux := http.NewServeMux()
// 	// mux.HandleFunc(ws2tcpCfg.Path, wsHandler)
// 	print("loading ", len(addrMap), " handleres\n")
// 	println(addrMap)
// 	for k, v := range addrMap {
// 		mux.HandleFunc(ws2tcpCfg.Path+"/"+k, getHandler(k, v))
// 	}
// 	// mux.HandleFunc(ws2tcpCfg.Path+"/remotecfg", func(w http.ResponseWriter, r *http.Request) {

// 	// })
// 	if ws2tcpCfg.Tls.Enabled {
// 		http.ListenAndServeTLS(ws2tcpCfg.Listen, ws2tcpCfg.Tls.CertFile, ws2tcpCfg.Tls.KeyFile, mux)
// 	} else {
// 		http.ListenAndServe(ws2tcpCfg.Listen, mux)
// 	}
// }

// func copyWorker(src io.Reader, dst io.Writer, doneCh chan bool) {
// 	io.Copy(dst, src)
// 	doneCh <- true
// }

// func getHandler(k, v string) func(w http.ResponseWriter, r *http.Request) {
// 	println("k: ", "/"+k, " v: ", v)
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		println("hit ", k, "lcoal ", v)
// 		println(r.URL.Path)
// 		conn, err := wsUp.Upgrade(w, r, nil)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		ioConn := rwc{c: conn}
// 		netConn, err := net.Dial("tcp", v)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		ch := make(chan bool)
// 		go copyWorker(netConn, &ioConn, ch)
// 		go copyWorker(&ioConn, netConn, ch)
// 		<-ch
// 		netConn.Close()
// 		ioConn.c.Close()
// 		<-ch
// 	}
// }

// type rwc struct {
// 	r io.Reader
// 	c *websocket.Conn
// }

// func (c *rwc) Write(p []byte) (int, error) {
// 	err := c.c.WriteMessage(websocket.BinaryMessage, p)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return len(p), nil
// }

// func (c *rwc) Read(p []byte) (int, error) {
// 	for {
// 		if c.r == nil {
// 			// Advance to next message.
// 			var err error
// 			_, c.r, err = c.c.NextReader()
// 			if err != nil {
// 				return 0, err
// 			}
// 		}
// 		n, err := c.r.Read(p)
// 		if err == io.EOF {
// 			// At end of message.
// 			c.r = nil
// 			if n > 0 {
// 				return n, nil
// 			} else {
// 				// No data read, continue to next message.
// 				continue
// 			}
// 		}
// 		return n, err
// 	}
// }

// // func wsHandler(w http.ResponseWriter, r *http.Request) {
// // 	conn, err := wsUp.Upgrade(w, r, nil)
// // 	if err != nil {
// // 		log.Println(err)
// // 	}
// // 	ioConn := rwc{c: conn}
// // 	netConn, err := net.Dial("tcp", tcpAddr)
// // 	if err != nil {
// // 		log.Println(err)
// // 	}
// // 	ch := make(chan bool)
// // 	go copyWorker(netConn, &ioConn, ch)
// // 	go copyWorker(&ioConn, netConn, ch)
// // 	<-ch
// // 	netConn.Close()
// // 	ioConn.c.Close()
// // 	<-ch
// // }
package main

import (
	"log"
	"miniFTL/service/tcp2ws"
	"miniFTL/service/ws2tcp"
	"miniFTL/util/config"
)

func main() {
	cfg := config.LoadFile("./config.json")
	if cfg.IsClient() && cfg.IsServer() {
		log.Println("do not use example config")
	}
	if cfg.IsServer() {
		ws2tcp.Start(cfg)
	} else if cfg.IsClient() {
		tcp2ws.Start(cfg)
	} else {
		log.Println("config file error")
	}
}
