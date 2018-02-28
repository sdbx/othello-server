package othello

type RoomStore struct {
	*Hub
}

func NewRoomStore(userStore *UserStore) *RoomStore {
	return &RoomStore{
		Hub: NewHub(userStore),
	}
}
