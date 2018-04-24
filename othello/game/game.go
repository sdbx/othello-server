package game

import (
	"errors"
	"sync"

	"github.com/olebedev/emitter"
	"github.com/sdbx/othello-server/othello/ws"
)

type (
	Game struct {
		sync.Mutex
		Black     string
		White     string
		BlackTime uint
		WhiteTime uint
		Board     Board
		History   History
		GameType  GameType
		Emitter   *emitter.Emitter
		gameRoom  *gameRoom
		undoing   bool
		undoIndex int
	}
)

func newGame(gameRoom *gameRoom, black string, white string, gameType GameType) *Game {
	return &Game{
		Black:     black,
		White:     white,
		BlackTime: gameType.Time(),
		WhiteTime: gameType.Time(),
		Board:     gameType.Initial(),
		History:   []Move{},
		Emitter:   &emitter.Emitter{},
		GameType:  gameType,
		gameRoom:  gameRoom,
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

// black / white
func (g *Game) Total() (int, int) {
	black := 0
	white := 0
	size := g.GameType.Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			if tile := g.GetTile(Coordinate{x, y}); tile == GameTileBlack {
				black++
			} else if tile == GameTileWhite {
				white++
			}
		}
	}
	return black, white
}

func (g *Game) TimeCount() {
	if g.Turn() == GameTurnBlack {
		g.BlackTime--
		if g.BlackTime == 0 {
			<-g.Emitter.Emit("end", ws.H{
				"winner": "white",
				"cause":  "timeout",
			})
		}
	} else {
		g.WhiteTime--
		if g.WhiteTime == 0 {
			<-g.Emitter.Emit("end", ws.H{
				"winner": "black",
				"cause":  "timeout",
			})
		}
	}
}

func (g *Game) UndoReq(turn Turn) {
	if !g.undoing {
		if g.Turn() == turn {
			g.undoIndex = len(g.History) - 2
		} else {
			g.undoIndex = len(g.History) - 1
		}
		g.undoing = true
		<-g.Emitter.Emit("undo", ws.H{
			"index": g.undoIndex,
			"color": turn,
		})
	}
}

func (g *Game) GoBackTo(index int) {
	if index > len(g.History) {
		return
	}
	g.History = g.History[:index]
	g.Board = g.GameType.Initial()
	tile := GameTileBlack
	history := g.History

	for _, mv := range history {
		cord, _ := CordFromMove(mv)
		g.put(cord, tile)
		tile = tile.GetFlip()
	}
}

func (g *Game) UndoAnswer(ok bool) {
	if ok {
		g.GoBackTo(g.undoIndex)
	}
	g.undoing = false
	<-g.Emitter.Emit("undo_answer", ws.H{
		"ok":    ok,
		"index": g.undoIndex,
	})

}

func (g *Game) CheckEnd() bool {
	if g.NumberOfPossibleMoves(GameTurnBlack) == 0 &&
		g.NumberOfPossibleMoves(GameTurnWhite) == 0 {
		black, white := g.Total()
		winner := ""
		if black > white {
			winner = "black"
		} else if black < white {
			winner = "white"
		} else {
			winner = "drew"
		}
		<-g.Emitter.Emit("end", ws.H{
			"winner": winner,
			"cause":  "normally",
		})
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
	g.History = append(g.History, move)
	<-g.Emitter.Emit("turn", ws.H{
		"color": turn,
		"move":  move,
	})
	if g.CheckEnd() {
		return nil
	}
	if g.NumberOfPossibleMoves(turn.GetOpp()) == 0 {
		g.History = append(g.History, MoveNone)
		<-g.Emitter.Emit("turn", ws.H{
			"color": turn.GetOpp(),
			"move":  MoveNone,
		})
	}
	return nil
}

func (g *Game) put(cord Coordinate, tile Tile) {
	if cord == CordNone {
		return
	}
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

// not thread safe!
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
