package main

import (
	"github.com/shinjitsu/TableTopRulerServer/GameData"
	"math/rand"
)

//import "github.com/shinjitsu/TableTopRulerServer/GameData"

func drawFourLands() []*GameData.LandTile {
	firstLands := make([]*GameData.LandTile, 0)
	var LandKind GameData.LandType
	var LandName string
	var FileName string
	for i := 0; i < 4; i++ {
		kind := rand.Intn(int(LAND_MAX-1)) + 1 //adjust for skipping instructions
		switch kind {
		case 1:
			LandKind = GameData.LandType_FOREST
			LandName = "Forest"
			FileName = "Forest.png"
		case 2:
			LandKind = GameData.LandType_MOUNTAIN
			LandName = "Mountain"
			FileName = "Mountains.png"
		case 3:
			LandKind = GameData.LandType_DESERT
			LandName = "Desert"
			FileName = "Desert.png"
		case 4:
			LandKind = GameData.LandType_SWAMP
			LandName = "Swamp"
			FileName = "Swamp.png"
		case 5:
			LandKind = GameData.LandType_PLAINS
			LandName = "Plains"
			FileName = "Plains.png"
		}
		land := GameData.LandTile{
			TileType: LandKind,
			TileName: LandName,
			Pict:     FileName,
		}

		firstLands = append(firstLands, &land)
	}
	return firstLands
}

func initializeDomain() []*GameData.DomainSpot {
	playerDomain := make([]*GameData.DomainSpot, 0)
	initialLands := drawFourLands()
	for _, Land := range initialLands {
		domainSpot := GameData.DomainSpot{
			Land: Land,
		}
		playerDomain = append(playerDomain, &domainSpot)
	}
	return playerDomain
}
