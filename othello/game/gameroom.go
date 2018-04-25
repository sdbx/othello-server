package game

import (
	"time"

	"github.com/sdbx/othello-server/othello/ws"
)

type gameRoom struct {
	*ws.WSRoom
	game  *Game
	timer *time.Timer
}

func retimer(t *time.Timer, d time.Duration) {
	if !t.Stop() {
		<-t.C
	}
	t.Reset(d)
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
			g.game.Emitter.Off("*")
			g.Close()
			return
		case event := <-turn:
			arg := event.Args[0].(ws.H)
			g.Emit("turn", arg)
			go func() {
				g.game.RLock()
				defer g.game.RUnlock()
				if arg["color"].(Turn) == "black" {
					retimer(g.timer, g.game.WhiteTime)
				}
				if arg["color"].(Turn) == "white" {
					retimer(g.timer, g.game.BlackTime)
				}
			}()
		case event := <-undo:
			g.Emit("undo", event.Args[0].(ws.H))
		case event := <-undoAns:
			g.Emit("undo_answer", event.Args[0].(ws.H))
		case <-g.timer.C:
			// prevent deadlock
			go func() {
				g.game.Lock()
				defer g.game.Unlock()
				g.game.end()
			}()
		}
	}
}
