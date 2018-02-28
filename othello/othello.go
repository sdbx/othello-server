package othello

type h map[string]interface{}

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
