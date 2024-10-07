package main

type TableTopRulerGame struct {
	Players                    []Player
	LandDeck                   []LandTile
	AvailalbeSpecialCharacters []SpecialCharacter
	PlayingDeck                []any
}

type Player struct {
}

type LandTile struct {
	TileType string
	TileName string
}

type SpecialCharacter struct {
}
