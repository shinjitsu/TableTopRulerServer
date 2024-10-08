package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type TableTopRulerGame struct {
	Players                    []Player
	LandDeck                   []LandTile
	AvailalbeSpecialCharacters []SpecialCharacter
	PlayingDeck                []any
}

type Player struct {
	Name           string
	Code           string
	PrestigePoints int
}

type LandTile struct {
	TileType string
	TileName string
}

type SpecialCharacter struct {
	Name           string
	Pict           *ebiten.Image
	CombatValue    int
	SpecialAbility string
}
