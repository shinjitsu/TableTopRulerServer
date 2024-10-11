package main

import (
	"TableTopRulerServer/GameData"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	GameData.UnimplementedTableTopRulerServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	Game := Server{}
	GameData.RegisterTableTopRulerServiceServer(server, &Game)
	log.Println("Server listening at", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalln("Server exited with error:", err)
	}
}

//func RunServer(host enet.Host, game TableTopRulerGame) {
//	for { //for ever
//		event := host.Service(1000) // Wait until the next event, 1000 is timeout
//		if event.GetType() == enet.EventNone {
//			continue // No events no event means don't do anything
//		}
//		switch event.GetType() {
//		case enet.EventConnect:
//			if game.GameStarted {
//				log.Printf("Client Tried to connect while game is already started: %s", event.GetPeer().GetAddress().String())
//				continue
//			} else {
//				log.Printf("Client connected: %s", event.GetPeer().GetAddress().String())
//
//			}
//		}
//	}
//}
