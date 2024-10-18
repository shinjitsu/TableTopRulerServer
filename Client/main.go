package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/shinjitsu/TableTopRulerServer/GameData"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"path"
	"sync"
)

//go:embed Assets/*
var embeddedAssets embed.FS

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

const windowWidth = 1500
const windowHeight = 1000
const PlayerNameHeight = 50
const BufferSpace = 10

var images map[string]*ebiten.Image

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

	textOpts := &text.DrawOptions{}
	textOpts.GeoM.Translate(200, 200)
	g.lock.RLock()
	localstate := g.gameState
	g.lock.RUnlock()
	drawOps := &ebiten.DrawImageOptions{}
	DrawPlayer1(screen, localstate.Player1, drawOps)
	DrawPlayer2(screen, localstate.Player2, drawOps)
	text.Draw(screen, g.GameMessage, mplusNormalFace, textOpts)
}

func DrawPlayer2(screen *ebiten.Image, player2 *GameData.Player, ops *ebiten.DrawImageOptions) {
	//player2 wll be drawn at the top
	if player2 == nil { //if we haven't connected player2 yet
		return
	}
	textOpts := &text.DrawOptions{}
	textOpts.GeoM.Translate(windowWidth/4, 0)
	player1Banner := fmt.Sprintf("Player: %s Prestige Points: %d Gold: %d", player2.Name, player2.PrestigePoints, player2.Gold)
	text.Draw(screen, player1Banner, mplusNormalFace, textOpts) //DrawPlayer Name at top
	//Draw Land
	for spotNum, domainSpot := range player2.Domain {
		textOpts.GeoM.Reset()
		land := domainSpot.Land
		image := GetImage(land.Pict)
		if image == nil {
			fmt.Println("Image not found!!!!!!!!:", land.TileName, land.Pict)
		}
		ops.GeoM.Translate(float64((spotNum)*(image.Bounds().Dx()+BufferSpace)), float64(PlayerNameHeight+BufferSpace))
		screen.DrawImage(image, ops)
		textOpts.GeoM.Translate(float64((spotNum)*(image.Bounds().Dx()+BufferSpace)+BufferSpace), float64(PlayerNameHeight+BufferSpace))
		text.Draw(screen, land.TileName, mplusNormalFace, textOpts)
		ops.GeoM.Reset()
	}
}

func DrawPlayer1(screen *ebiten.Image, player1 *GameData.Player, ops *ebiten.DrawImageOptions) {
	//Player1 wll be drawn at the
	if player1 == nil { //if we haven't connected player1 yet
		return
	}
	textOpts := &text.DrawOptions{}
	textOpts.GeoM.Translate(windowWidth/4, windowHeight-PlayerNameHeight)
	player1Banner := fmt.Sprintf("Player: %s Prestige Points: %d Gold: %d", player1.Name, player1.PrestigePoints, player1.Gold)
	text.Draw(screen, player1Banner, mplusNormalFace, textOpts) //DrawPlayer Name at top

	//Draw Land
	for spotNum, domainSpot := range player1.Domain {
		textOpts.GeoM.Reset()
		land := domainSpot.Land
		image := GetImage(land.Pict)
		//fmt.Println(land.TileName)
		ops.GeoM.Translate(float64((spotNum)*(image.Bounds().Dx()+BufferSpace)), float64(windowHeight-(PlayerNameHeight+image.Bounds().Dy()+BufferSpace)))
		screen.DrawImage(image, ops)
		textOpts.GeoM.Translate(float64((spotNum)*(image.Bounds().Dx()+BufferSpace)+BufferSpace), float64(windowHeight-(PlayerNameHeight+image.Bounds().Dy()+BufferSpace)))
		text.Draw(screen, land.TileName, mplusNormalFace, textOpts)
		ops.GeoM.Reset()
	}
}

func GetImage(imageName string) *ebiten.Image {
	if image, ok := images[imageName]; ok {
		return image
	}
	image := LoadEmbeddedImage("Images", imageName)
	images[imageName] = image
	return image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	images = make(map[string]*ebiten.Image)
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
		log.Fatal("Error connecting to server:", err)
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

func LoadEmbeddedImage(folderName string, imageName string) *ebiten.Image {
	embeddedFile, err := embeddedAssets.Open(path.Join("Assets", folderName, imageName))
	if err != nil {
		log.Fatal("failed to load embedded image ", imageName, err)
	}
	ebitenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("Error loading tile image:", imageName, err)
	}
	return ebitenImage
}
