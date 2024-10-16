package main

import (
	"TableTopRulerServer/GameData"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

// const maxPlayers = 4 the real player count
const maxPlayers = 2 //for testing only

type Server struct {
	GameData.UnimplementedTableTopRulerServiceServer
	Players       []Player
	PlayerStreams map[string]GameData.TableTopRulerService_ReceiveGameEventsServer
	CurrentPlayer Player
	GameStarted   bool
	TurnNumber    int32
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
		StandingArmy:   make([]*Unit, 0, 4),
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

func (s *Server) ReceiveGameEvents(_ *GameData.Empty, stream GameData.TableTopRulerService_ReceiveGameEventsServer) error {
	metaData, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("No metadata")
	}
	playerCode := metaData.Get("playerCode")[0]
	if playerCode == "" {
		return errors.New("No playerCode")
	}
	//might need a mutex here
	s.PlayerStreams[playerCode] = stream
	<-stream.Context().Done() //keep stream open
	return nil
}

func (s *Server) PlayTurn(ctx context.Context, request *GameData.TempTurn) (*GameData.TempResponse, error) {
	if request.Name != s.CurrentPlayer.Name {
		return nil, errors.New("Not your turn")
	}
	s.TurnNumber++
	s.BroadcastTurn()
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

func (s *Server) BroadcastTurn() {
	for playerCode, stream := range s.PlayerStreams {
		if playerCode == s.CurrentPlayer.Code {
			stream.Send(&GameData.GameState{
				TurnNumber: s.TurnNumber,
				Player1: &GameData.Player{
					Name:           s.Players[0].Name,
					PrestigePoints: s.Players[0].PrestigePoints,
					StandingArmy:   nil, //fix - but lets see communication first
					Domain:         nil, //fix - but lets see communication first
				},

				Player2: &GameData.Player{
					Name:           s.Players[1].Name,
					PrestigePoints: s.Players[1].PrestigePoints,
					StandingArmy:   nil, //fix - but lets see communication first
					Domain:         nil, //fix - but lets see communication first
				},
				Player3: nil,
				Player4: nil,
				Winner:  -1,
			})
		}
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
