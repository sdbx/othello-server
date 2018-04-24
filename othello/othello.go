package othello

import (
	"github.com/sdbx/othello-server/othello/game"
	"github.com/sdbx/othello-server/othello/room"
)

type Service struct {
	GameStore *game.GameStore
	RoomStore *room.RoomStore
}

func NewService() *Service {
	gameStore := game.NewGameStore()
	return &Service{
		GameStore: gameStore,
		RoomStore: room.NewRoomStore(gameStore),
	}
}
