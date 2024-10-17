package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/shinjitsu/TableTopRulerServer/GameData"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"sync"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

type Game struct {
	client      GameData.TableTopRulerServiceClient
	Name        string
	Code        string
	lock        sync.RWMutex
	stream      GameData.TableTopRulerService_ReceiveGameEventsClient
	gameState   GameData.GameState
	GameMessage string
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		resp, err := g.client.PlayTurn(context.Background(), &GameData.TempTurn{
			Name:                  g.Name,
			Code:                  g.Code,
			TempturnDisplayString: fmt.Sprintf("%s Playingg Turn %d", g.Name, g.gameState.TurnNumber),
		})
		if err != nil {
			g.GameMessage = "Not your Turn"
		} else {
			g.GameMessage = resp.TempResponse
		}
	}
	if g.stream == nil {
		return nil //not ready yet
	}
	//ToDo: fix this Replace with go routine
	if g.gameState.Player1 == nil {
		return nil
	}
	g.lock.RLock()
	g.GameMessage = fmt.Sprintf("Turn %d, Player1: %s, Player2: %s",
		g.gameState.TurnNumber, g.gameState.Player1.Name, g.gameState.Player2.Name)
	g.lock.RUnlock()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//TODO implement me
	//screen.Fill(color.White)
	//drawOps := ebiten.DrawImageOptions{}
	textOpts := &text.DrawOptions{}
	textOpts.GeoM.Translate(200, 200)

	text.Draw(screen, g.GameMessage, mplusNormalFace, textOpts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

const windowWidth = 1000
const windowHeight = 1080

func main() {

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Tabletop Emperor")
	name := getPlayerName()
	gameClient := Game{Name: name}
	request := GameData.GetPlayersRequest{Name: name}
	serverConnection, err := grpc.NewClient(fmt.Sprintf("localhost:%d", 9090),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	gameClient.client = GameData.NewTableTopRulerServiceClient(serverConnection)
	response, err := gameClient.client.Connect(context.Background(), &request)
	if err != nil {
		log.Fatal("DEBUG!!!!!", err)
	}
	gameClient.Code = response.Code
	metaD := metadata.New(map[string]string{"playerCode": gameClient.Code})
	ctx := metadata.NewOutgoingContext(context.Background(), metaD)
	gameClient.stream, err = gameClient.client.ReceiveGameEvents(ctx, &GameData.Empty{})
	if err != nil {
		log.Fatal("Error Preparing to get Game events:", err)
	}

	receiveGameEvents := func(gameClient *Game) {
		for {
			gameState, err := gameClient.stream.Recv()
			if err != nil {
				log.Fatal("Error getting Game events:", err)
			}
			if gameState == nil {
				return
			}
			gameClient.lock.Lock()
			gameClient.gameState = *gameState
			gameClient.lock.Unlock()
		}
	}
	go receiveGameEvents(&gameClient)
	gameClient.GameMessage = fmt.Sprintf("Start of Game %s Connected", name)
	//dataStream, err := gameClient.client.ReceiveGameEvents(context.Background(), &GameData.Empty{})
	//gameClient.stream = dataStream
	if err := ebiten.RunGame(&gameClient); err != nil {
		log.Fatal("Error running game:", err)
	}
}

func getPlayerName() string {
	var name string
	log.Print("Enter your name: ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

// swiped from the ebitengine text demo
var (
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
	mplusBigFace    *text.GoTextFace
)

func init() { //font loading needed to be in init as of 2024
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	mplusNormalFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}
	mplusBigFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   32,
	}
}
