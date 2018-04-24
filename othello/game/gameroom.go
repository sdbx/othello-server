package game

import (
	"time"

	"github.com/sdbx/othello-server/othello/ws"
)

type gameRoom struct {
	*ws.WSRoom
	game     *Game
	ticker1  *time.Ticker
	ticker10 *time.Ticker
}

func (g *gameRoom) runGame() {
	end := g.game.Emitter.On("end")
	turn := g.game.Emitter.On("turn")
	undo := g.game.Emitter.On("undo")
	undoAns := g.game.Emitter.On("undo_answer")
	for {
		select {
		case event := <-end:
			g.Emit("end", event.Args[0].(ws.H))
			g.Close()
			return
		case event := <-turn:
			g.Emit("turn", event.Args[0].(ws.H))
		case event := <-undo:
			g.Emit("undo", event.Args[0].(ws.H))
		case event := <-undoAns:
			g.Emit("undo_answer", event.Args[0].(ws.H))
		case <-g.ticker1.C:
			g.game.TimeCount()
		case <-g.ticker10.C:
			g.Emit("tick", ws.H{
				"black": g.game.BlackTime,
				"white": g.game.WhiteTime,
			})
		}
	}
}
