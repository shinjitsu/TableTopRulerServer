package main

import "math/rand"

//import "github.com/shinjitsu/TableTopRulerServer/GameData"

func drawFourLands() []LandTile {
	firstLands := make([]LandTile, 4)
	var LandKind LandType
	var LandName string
	for i := 0; i < 4; i++ {
		kind := rand.Intn(int(LAND_MAX)) + 1
		switch kind {
		case 1:
			LandKind = Forest
			LandName = "Forest"
		case 2:
			LandKind = Mountain
			LandName = "Mountain"
		case 3:
			LandKind = Desert
			LandName = "Desert"
		case 4:
			LandKind = Swamp
			LandName = "Swamp"
		case 5:
			LandKind = Plains
			LandName = "Plains"
		}
		land := LandTile{
			TileType: LandKind,
			TileName: LandName,
			Pict:     nil,
		}

		firstLands = append(firstLands, land)
	}
	return firstLands
}

func initializeDomain() []DomainSpot {
	playerDomain := make([]DomainSpot, 4)
	initialLands := drawFourLands()
	for _, Land := range initialLands {
		domainSpot := DomainSpot{
			Land: Land,
		}
		playerDomain = append(playerDomain, domainSpot)
	}
	return playerDomain
}
