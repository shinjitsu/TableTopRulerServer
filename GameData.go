package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SpecialAbility int

const (
	None SpecialAbility = iota
	Flying
	Magic
	ChargeBonus
	RangedAttack
)

type LandType int

const (
	Instruction LandType = iota
	Forest
	Mountain
	Desert
	Swamp
	Plains
)

type TableTopRulerGame struct {
	//	GameData.UnimplementedTableTopRulerServiceServer
	GameStarted bool
	Players     []Player
	//	PlayerStreams map[string]GameData.TableTopRulerService_RecieveGameEventsServer
	LandDeck                   []LandTile
	AvailableSpecialCharacters []SpecialCharacter
	PlayingDeck                []any
}

type DomainSpot struct {
	Land          LandTile
	Upgrade       Improvement
	Fortification Fortification
}

type Player struct {
	Name           string
	Code           string
	PrestigePoints int32
	StandingArmy   []*Unit
	Domain         []DomainSpot
	Gold           int32
	Hand           []any
}

type LandTile struct {
	TileType LandType
	TileName string
	Pict     *ebiten.Image
}

type Unit struct {
	Name           string
	Pict           *ebiten.Image
	CombatValue    int
	SpecialAbility SpecialAbility
}

type Improvement struct {
	Name         string
	Pict         *ebiten.Image
	DefenseValue int
	GoldValue    int
}

type Fortification struct {
	Name                 string
	Pict                 *ebiten.Image
	DefensePrestigeValue int
}

type SpecialCharacter struct {
	Name           string
	Pict           *ebiten.Image
	CombatValue    int
	SpecialAbility SpecialAbility
}
