package main

import (
	"TableTopRulerServer/GameData"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
)

// const maxPlayers = 4 the real player count
const maxPlayers = 2 //for testing only

type Server struct {
	GameData.UnimplementedTableTopRulerServiceServer
	Players       []Player
	CurrentPlayer Player
	GameStarted   bool
}

func main() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	Game := Server{
		Players: make([]Player, 0, 4), // 2-4 concurrent players
	}
	GameData.RegisterTableTopRulerServiceServer(server, &Game)
	log.Println("Server listening at", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalln("Server exited with error:", err)
	}
}

func (s *Server) Connect(ctx context.Context, request *GameData.GetPlayersRequest) (*GameData.GetPlayersResponse, error) {
	if s.GameStarted {
		return nil, errors.New("Game already started")
	} else if len(s.Players) >= maxPlayers {
		return nil, errors.New("Max players reached")
	}
	newPlayer := Player{
		Name:           request.Name,
		Code:           uuid.New().String(),
		PrestigePoints: 0,
		StandingArmy:   make([]Unit, 0, 4),
		Domain:         make([]DomainSpot, 0, 4),
		Gold:           0,
		Hand:           make([]any, 0, 4),
	}
	s.Players = append(s.Players, newPlayer)
	if len(s.Players) == maxPlayers {
		s.GameStarted = true
	}
	return &GameData.GetPlayersResponse{
		Name: newPlayer.Name,
		Code: newPlayer.Code,
	}, nil
}

func (s *Server) PlayTurn(ctx context.Context, request *GameData.TempTurn) (*GameData.TempResponse, error) {
	if request.Name != s.CurrentPlayer.Name {
		return nil, errors.New("Not your turn")
	}
	//now broadcast to all players
	return &GameData.TempResponse{
		TempResponse: fmt.Sprintf("Turn %s played", request.Name),
	}, nil
}

func (s *Server) Defend(ctx context.Context, request *GameData.TempDefend) (*GameData.TempDefendResponse, error) {
	if request.Name == s.CurrentPlayer.Name {
		return nil, errors.New("You are not Defending")
	}
	//now broadcast to all players
	return &GameData.TempDefendResponse{
		TempResponse: fmt.Sprintf("Turn %s Defended", request.Name),
	}, nil
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
