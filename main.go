package main

import (
	"log"
	"miniFTL/service/tcp2ws"
	"miniFTL/service/ws2tcp"
	"miniFTL/util/config"
)

//TODO:add output, add authentication
//TODO:implent UDP2WS at branch fusion2
// https://varshneyabhi.wordpress.com/2014/12/23/simple-udp-clientserver-in-golang/
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
