package othello

import (
	"github.com/sdbx/othello-server/othello/db"
	"github.com/sdbx/othello-server/othello/game"
	"github.com/sdbx/othello-server/othello/models"
	"github.com/sdbx/othello-server/othello/room"
)

type Service struct {
	GameStore *game.GameStore
	RoomStore *room.RoomStore
}

func NewService() *Service {
	var userStore db.DBUserStore = map[string]*models.User{}
	gameStore := game.NewGameStore(&userStore)
	return &Service{
		UserStore: &userStore,
		GameStore: gameStore,
		RoomStore: room.NewRoomStore(&userStore, gameStore),
	}
}
