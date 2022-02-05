package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	// "golang.org/x/net/websocket"
)

// TODO: fix flags

var tcpAddr string

var wsUp = websocket.Upgrader{}

type cfg struct {
	Listen  string `json:"listen"`
	Path    string `json:"path"`
	TcpAddr string `json:"tcpaddr"`
}

func main() {
	// // parse flags
	// var addr string
	// var path string
	// flag.StringVar(&addr, "addr", "127.0.0.1:8090", "define address to listen")
	// flag.StringVar(&path, "path", "/ws", "define path of websocket")
	// flag.Usage = usage
	// flag.Parse()
	// tcpAddr = flag.Arg(0)
	// if tcpAddr == "" {
	// 	log.Fatal("no TCP address provided")
	// }

	//read config
	var ws2tcpCfg cfg
	f, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	b, _ := ioutil.ReadAll(f)
	err = json.Unmarshal(b, &ws2tcpCfg)
	if err != nil {
		log.Fatal(err)
	}
	tcpAddr = ws2tcpCfg.TcpAddr

	// start web server
	mux := http.NewServeMux()
	mux.HandleFunc(ws2tcpCfg.Path, wsHandler)
	http.ListenAndServe(ws2tcpCfg.Listen, mux)
}

func copyWorker(src io.Reader, dst io.Writer, doneCh chan bool) {
	io.Copy(dst, src)
	doneCh <- true
}

type rwc struct {
	r io.Reader
	c *websocket.Conn
}

func (c *rwc) Write(p []byte) (int, error) {
	err := c.c.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *rwc) Read(p []byte) (int, error) {
	for {
		if c.r == nil {
			// Advance to next message.
			var err error
			_, c.r, err = c.c.NextReader()
			if err != nil {
				return 0, err
			}
		}
		n, err := c.r.Read(p)
		if err == io.EOF {
			// At end of message.
			c.r = nil
			if n > 0 {
				return n, nil
			} else {
				// No data read, continue to next message.
				continue
			}
		}
		return n, err
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUp.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	ioConn := rwc{c: conn}
	netConn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		log.Println(err)
	}
	ch := make(chan bool)
	go copyWorker(netConn, &ioConn, ch)
	go copyWorker(&ioConn, netConn, ch)
	<-ch
	netConn.Close()
	ioConn.c.Close()
	<-ch
}
