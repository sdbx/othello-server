package game

type (
	GameType interface {
		Name() string
		Initial() Board
		Size() Coordinate
		Time() uint
	}
	DefaultOthello struct {
	}
)

func (d DefaultOthello) Name() string {
	return "default"
}

func (d DefaultOthello) Initial() Board {
	ma := make([][]Tile, 8)
	for i := 0; i < 8; i++ {
		ma[i] = make([]Tile, 8)
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			ma[i][j] = GameTileNone
		}
	}
	ma[3][3] = GameTileWhite
	ma[3][4] = GameTileBlack
	ma[4][3] = GameTileBlack
	ma[4][4] = GameTileWhite
	return ma
}

func (d DefaultOthello) Size() Coordinate {
	return Coordinate{8, 8}
}

func (d DefaultOthello) Time() uint {
	return 10
}
