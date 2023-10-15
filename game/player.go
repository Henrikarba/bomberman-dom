package game

import "time"

type Movement struct {
	Type string `json:"type"`

	PlayerID int      `json:"-"`
	Keys     []string `json:"keys"`
}

type Player struct {
	ID        int    `json:"id,omitempty"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction,omitempty"`

	LastMoveTime   time.Time `json:"-"`
	Speed          int       `json:"-"`
	AvailableBombs int       `json:"bombs,omitempty"`
	Lives          int       `json:"lives,omitempty"`
	// PowerUps       []PowerUp
}

func NewPlayer(id int, gameboard [][]string) *Player {
	player := &Player{
		ID:             id,
		Speed:          300,
		AvailableBombs: 1,
	}
	rows := len(gameboard)
	cols := len(gameboard[0])

	switch id {
	case 1:
		player.X = 0
		player.Y = 0
		player.Direction = "right"
	case 2:
		player.X = cols - 1
		player.Y = 0
		player.Direction = "left"
	case 3:
		player.X = 0
		player.Y = rows - 1
		player.Direction = "up"
	case 4:
		player.X = cols - 1
		player.Y = rows - 1
		player.Direction = "up"
	}

	return player
}
