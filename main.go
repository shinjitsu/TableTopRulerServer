package main

import (
	"github.com/codecat/go-enet"
	"log"
)

func main() {
	// Initialize enet
	enet.Initialize()

	// Create a host listening on  port 8095
	host, err := enet.NewHost(enet.NewListenAddress(8095), 32, 1, 0, 0)
	if err != nil {
		log.Fatal("Couldn't create host: %s", err.Error())
		return
	}
	//mainGame := SharedData.MPgame{
	//	Players: make([]*SharedData.Player, 0, 20), // no more than 20 concurrent players
	//	Gold:    makeGold(),
	//}
	mainGame := "Hello World"
	RunServer(host, mainGame)

}

func RunServer(host enet.Host, game any) {

}
