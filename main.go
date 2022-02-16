package main

import (
	"WSTunnel/service/tcp2ws"
	"WSTunnel/service/ws2tcp"
	"WSTunnel/util/config"
	"log"
)

/*
=================TODO=================================
02 fix udp, ws convert at fusion3 reference:
	https://github.com/magic000/udp2tcp
05 improve authentication method
-----------------DONE---------------------------------
00 add output, add authentication DONE
01 implent UDP2WS at branch fusion2 INCOMPLETE
	https://varshneyabhi.wordpress.com/2014/12/23/simple-udp-clientserver-in-golang/
03 remove incomplete UDP support in fusion2 and make it release-ready(fusion2)
04 add authentication in fusion2A
=============END TODO=================================
*/

func main() {
	log.Println("WSTunnel -- Expose local services via websocket")
	cfg := config.LoadFile("./config.json")
	log.Println("config file loaded")
	if cfg.IsClient() && cfg.IsServer() {
		log.Println("do not use example config")
	}
	if cfg.IsServer() {
		log.Println("Starting a Tunnel Server")
		ws2tcp.Start(cfg)
	} else if cfg.IsClient() {
		log.Println("Starting a Tunnel Client")
		tcp2ws.Start(cfg)
	} else {
		log.Println("config file error")
	}
}
