package othello

type Service struct {
	UserStore UserStore
}

func NewService() *Service {
	return &Service{
		UserStore: make(map[string]*User),
	}
}
