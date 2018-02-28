package othello

type Service struct {
	UserStore *UserStore
	RoomStore *RoomStore
}

func NewService() *Service {
	var userStore UserStore = map[string]*User{}
	return &Service{
		UserStore: &userStore,
		RoomStore: NewRoomStore(&userStore),
	}
}
