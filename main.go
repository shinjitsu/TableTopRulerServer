package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinjitsu/TableTopRulerServer/GameData"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
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
		Players:       make([]Player, 0, 4),
		PlayerStreams: make(map[string]GameData.TableTopRulerService_ReceiveGameEventsServer),
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
		PrestigePoints: 4, //start with 4 prestige points one for each Land
		StandingArmy:   make([]*Unit, 0, 4),
		Domain:         initializeDomain(),
		Gold:           rand.Int31n(6) + 1 + rand.Int31n(6) + 1, //roll 2D6 for starting gold
		Hand:           make([]any, 0, 4),
	}
	s.Players = append(s.Players, newPlayer)
	if len(s.Players) == maxPlayers {
		s.GameStarted = true
		s.CurrentPlayer = s.Players[0]
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
	//move to next player
	currentPlayerNum := 0
	for loc, player := range s.Players {
		if player.Name == s.CurrentPlayer.Name {
			currentPlayerNum = loc
		}
	}
	currentPlayerNum++
	if currentPlayerNum >= maxPlayers {
		currentPlayerNum = 0
	}
	s.CurrentPlayer = s.Players[currentPlayerNum]
	//end move to next player
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
	for _, stream := range s.PlayerStreams {
		//if playerCode == s.CurrentPlayer.Code {
		stream.Send(&GameData.GameState{
			TurnNumber: s.TurnNumber,
			Player1: &GameData.Player{
				Name:           s.Players[0].Name,
				PrestigePoints: s.Players[0].PrestigePoints,
				StandingArmy:   nil,                 //fix - but lets see communication first
				Domain:         s.Players[0].Domain, //fix - but lets see communication first
				Gold:           s.Players[0].Gold,
			},

			Player2: &GameData.Player{
				Name:           s.Players[1].Name,
				PrestigePoints: s.Players[1].PrestigePoints,
				StandingArmy:   nil, //fix - but lets see communication first
				Domain:         s.Players[1].Domain,
				Gold:           s.Players[1].Gold,
			},
			Player3: nil,
			Player4: nil,
			Winner:  -1,
		})
		//	}
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
