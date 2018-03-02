package game

import (
	"errors"
)

type (
	h    map[string]interface{}
	Game struct {
		Black    string
		White    string
		Board    Board
		History  History
		GameType GameType
		gameRoom *gameRoom
	}
)

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

func (g *Game) CheckEnd() (ended bool) {
	defer func() {
		if ended {
			g.gameRoom.emit("end", h{})
			g.gameRoom.close <- true
		}
	}()

	if g.NumberOfPossibleMoves(GameTurnBlack) == 0 &&
		g.NumberOfPossibleMoves(GameTurnWhite) == 0 {
		return true
	}
	if len(g.History) == 60 {
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
				break
			}
			start = true
			continue
		}
		if ttile == tile {
			return true
		}
		if ttile != opp {
			break
		}
	}
	return false
}
