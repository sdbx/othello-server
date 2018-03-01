package othello

import (
	"errors"
	"regexp"
	"strconv"
)

type gameClientType uint

type Tile uint

const (
	GameTileBlack Tile = iota
	GameTileWhite
	GameTileNone
	GameTileInvalid
)

type Turn string

const (
	GameTurnBlack = "black"
	GameTurnWhite = "white"
)

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

func TurnFromTile(t Tile) Turn {
	if t == GameTileBlack {
		return GameTurnBlack
	}
	if t == GameTileWhite {
		return GameTurnWhite
	}
	return ""
}

type Move string

const MoveNone = "none"

type (
	Coordinate struct {
		X int
		Y int
	}
	History []Move
	Board   [][]Tile

	Game struct {
		Black    string
		White    string
		Board    Board
		History  History
		GameType GameType
		gameRoom *gameRoom
	}

	GameType interface {
		Name() string
		Initial() Board
		Size() Coordinate
	}

	DefaultOthello struct {
	}
)

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

func newGame(gameRoom *gameRoom, black string, white string, gameType GameType) *Game {
	return &Game{
		Black:    black,
		White:    white,
		Board:    gameType.Initial(),
		History:  []Move{},
		GameType: gameType,
		gameRoom: gameRoom,
	}
}

func (c Coordinate) ToMove() Move {
	y := strconv.Itoa(c.Y + 1)
	return Move(string(c.X+'a') + y)
}

func CordFromMove(move Move) (Coordinate, error) {
	r := regexp.MustCompile(`^([a-z]+)([0-9]{1})$`)
	arr := r.FindStringSubmatch(string(move))
	if len(arr) != 3 {
		return Coordinate{}, errors.New("invalid move")
	}
	x := arr[1][0] - 'a'
	y, _ := strconv.Atoi(arr[2])
	return Coordinate{int(x), y - 1}, nil
}

func (c *Coordinate) Add(cor Coordinate) {
	c.X += cor.X
	c.Y += cor.Y
}

func (c Coordinate) IsValid(sizeCord Coordinate) bool {
	return c.X >= 0 && c.Y >= 0 && c.X < sizeCord.X && c.Y < sizeCord.Y
}

var dirs = []Coordinate{{0, 1}, {0, -1}, {1, 0}, {-1, 0}, {1, 1}, {1, -1}, {-1, -1}, {-1, 1}}

func (g *Game) NumberOfPossibleMoves(turn Turn) uint {
	if turn != GameTurnBlack && turn != GameTurnWhite {
		return 0
	}
	num := uint(0)
	size := g.GameType.Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			if g.PossibleAt(Coordinate{x, y}, turn.GetTile()) {
				num++
			}
		}
	}
	return num
}

func (g *Game) CheckEnd() bool {
	if g.NumberOfPossibleMoves(GameTurnBlack) == 0 &&
		g.NumberOfPossibleMoves(GameTurnWhite) == 0 {
		g.gameRoom.emit("end", h{})
		g.gameRoom.close <- true
		return true
	}
	return false
}

func (g *Game) Put(cord Coordinate, tile Tile) error {
	turn := g.Turn()
	if turn != TurnFromTile(tile) {
		return errors.New("it's opponent's turn")
	}
	if !g.PossibleAt(cord, tile) {
		return errors.New("invalid move")
	}
	g.put(cord, tile)
	move := cord.ToMove()
	g.gameRoom.emit("turn", h{
		"color": turn,
		"move":  move,
	})
	g.History = append(g.History, move)
	if g.CheckEnd() {
		return nil
	}
	if g.NumberOfPossibleMoves(turn.GetOpp()) == 0 {
		g.History = append(g.History, MoveNone)
		g.gameRoom.emit("turn", h{
			"color": turn.GetOpp(),
			"move":  "none",
		})
	}
	return nil
}

func (g *Game) put(cord Coordinate, tile Tile) {
	opp := tile.GetFlip()
	for _, dir := range dirs {
		if g.PossibleInDir(cord, tile, dir) {
			temp := cord
			temp.Add(dir)
			for g.GetTile(temp) == opp {
				g.set(temp, tile)
				temp.Add(dir)
			}
		}
	}
	g.set(cord, tile)
}

func (g *Game) Turn() Turn {
	if len(g.History)%2 == 0 {
		return GameTurnBlack
	}
	return GameTurnWhite
}

func (g *Game) GetTile(cord Coordinate) Tile {
	if cord.IsValid(g.GameType.Size()) {
		return g.Board[cord.Y][cord.X]
	}
	return GameTileInvalid
}

func (g *Game) set(cord Coordinate, tile Tile) {
	g.Board[cord.Y][cord.X] = tile
}

func (g *Game) PossibleAt(cord Coordinate, tile Tile) bool {
	for _, dir := range dirs {
		if g.PossibleInDir(cord, tile, dir) {
			return true
		}
	}
	return false
}

func (g *Game) PossibleInDir(cord Coordinate, tile Tile, dir Coordinate) bool {
	if g.GetTile(cord) != GameTileNone {
		return false
	}
	if tile != GameTileBlack && tile != GameTileWhite {
		return false
	}
	opp := tile.GetFlip()
	size := g.GameType.Size()
	temp := cord
	start := false
	for temp.IsValid(size) {
		temp.Add(dir)
		ttile := g.GetTile(temp)
		if !start {
			if ttile != opp {
				return false
			}
			start = true
			continue
		}
		if ttile == tile {
			return true
		}
		if ttile != opp {
			return false
		}
	}
	return false
}

func (d DefaultOthello) Name() string {
	return "default"
}

func (d DefaultOthello) Initial() Board {
	ma := make([][]Tile, 8)
	for i := 0; i < 8; i++ {
		ma[i] = make([]Tile, 8)
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			ma[i][j] = GameTileNone
		}
	}
	ma[3][3] = GameTileWhite
	ma[3][4] = GameTileBlack
	ma[4][3] = GameTileBlack
	ma[4][4] = GameTileWhite
	return ma
}

func (d DefaultOthello) Size() Coordinate {
	return Coordinate{8, 8}
}
