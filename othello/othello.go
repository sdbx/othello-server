package othello

import (
	"github.com/sdbx/othello-server/othello/db"
	"github.com/sdbx/othello-server/othello/game"
	"github.com/sdbx/othello-server/othello/models"
)

type Service struct {
	UserStore models.UserStore
	GameStore *game.GameStore
}

func NewService() *Service {
	var userStore db.DBUserStore = map[string]*models.User{}
	return &Service{
		UserStore: &userStore,
		GameStore: game.NewGameStore(&userStore),
	}
}
