package main

import (
	"WSTunnel/service/tcp2ws"
	"WSTunnel/service/ws2tcp"
	"WSTunnel/util/config"
	"log"
)

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
