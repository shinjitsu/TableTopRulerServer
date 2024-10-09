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
	mainGame := TableTopRulerGame{
		Players:                    make([]Player, 0, 4), // 2-4 concurrent players
		LandDeck:                   make([]LandTile, 0, 54),
		AvailableSpecialCharacters: make([]SpecialCharacter, 0, 11),
		PlayingDeck:                make([]any, 0, 124),
	}
	RunServer(host, mainGame)

}

func RunServer(host enet.Host, game TableTopRulerGame) {
	for { //for ever
		event := host.Service(1000) // Wait until the next event, 1000 is timeout
		if event.GetType() == enet.EventNone {
			continue // No events no event means don't do anything
		}
		switch event.GetType() {
		case enet.EventConnect:
			if game.GameStarted {
				log.Printf("Client Tried to connect while game is already started: %s", event.GetPeer().GetAddress().String())
				continue
			} else {
				log.Printf("Client connected: %s", event.GetPeer().GetAddress().String())

			}
		}
	}
}
