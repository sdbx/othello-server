package othello

type Service struct {
	UserStore UserStore
	RoomStore RoomStore
}

func NewService() *Service {
	return &Service{
		UserStore: make(map[string]string),
		RoomStore: []Room{},
	}
}
