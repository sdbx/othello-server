package game

import (
	"errors"
	"regexp"
	"strconv"
)

type (
	History    []Move
	Board      [][]Tile
	Tile       uint
	Turn       string
	Move       string
	Coordinate struct {
		X int
		Y int
	}
)

const (
	GameTileBlack Tile = iota
	GameTileWhite
	GameTileNone
	GameTileInvalid
)

const (
	GameTurnBlack = "black"
	GameTurnWhite = "white"
)

const MoveNone = "--"

var CordNone = Coordinate{-1, -1}

func CordFromMove(move Move) (Coordinate, error) {
	if move == MoveNone {
		return CordNone, nil
	}
	r := regexp.MustCompile(`^([a-z]+)([0-9]{1})$`)
	arr := r.FindStringSubmatch(string(move))
	if len(arr) != 3 {
		return Coordinate{}, errors.New("invalid move")
	}
	x := arr[1][0] - 'a'
	y, _ := strconv.Atoi(arr[2])
	return Coordinate{int(x), y - 1}, nil
}

func (c Coordinate) ToMove() Move {
	if c == CordNone {
		return MoveNone
	}
	y := strconv.Itoa(c.Y + 1)
	return Move(string(c.X+'a') + y)
}

func (c *Coordinate) Add(cor Coordinate) {
	c.X += cor.X
	c.Y += cor.Y
}

func (c Coordinate) IsValid(sizeCord Coordinate) bool {
	return c.X >= 0 && c.Y >= 0 && c.X < sizeCord.X && c.Y < sizeCord.Y
}

func (t Tile) GetFlip() Tile {
	switch t {
	case GameTileBlack:
		return GameTileWhite
	case GameTileWhite:
		return GameTileBlack
	case GameTileNone:
		return GameTileNone
	default:
		return GameTileInvalid
	}
}

func TurnFromTile(t Tile) Turn {
	if t == GameTileBlack {
		return GameTurnBlack
	}
	if t == GameTileWhite {
		return GameTurnWhite
	}
	return ""
}

func (t Turn) GetTile() Tile {
	if t == GameTurnBlack {
		return GameTileBlack
	}
	if t == GameTurnWhite {
		return GameTileWhite
	}
	return GameTileInvalid
}

func (t Turn) GetOpp() Turn {
	switch t {
	case GameTurnBlack:
		return GameTurnWhite
	case GameTurnWhite:
		return GameTurnBlack
	default:
		return ""
	}
}
