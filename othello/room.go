package othello

type Room struct {
	Name  string   `json:"name"`
	King  string   `json:"king"`
	Users []string `json:"users"`
}

type RoomStore []Room

func (rs *RoomStore) AddRoom(username string, name string) error {
	*rs = append(*rs, Room{
		Name:  name,
		King:  username,
		Users: []string{username},
	})
	return nil
}
