package othello

type h map[string]interface{}

type Service struct {
	UserStore *UserStore
	GameStore *GameStore
}

func NewService() *Service {
	var userStore UserStore = map[string]*User{}
	return &Service{
		UserStore: &userStore,
		GameStore: NewGameStore(&userStore),
	}
}
