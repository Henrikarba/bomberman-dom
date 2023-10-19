package game

import "time"

type Movement struct {
	Type string `json:"type"`

	PlayerID int      `json:"-"`
	Keys     []string `json:"keys"`
}

type Player struct {
	Name      string `json:"name"`
	ID        int    `json:"id,omitempty"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction,omitempty"`

	LastMoveTime   time.Time `json:"-"`
	Speed          int       `json:"-"`
	AvailableBombs int       `json:"bombs,omitempty"`
	FireDistance   int       `json:"-"`
	Lives          int       `json:"lives"`
	Damaged        bool      `json:"damaged"`
	// PowerUps       []PowerUp
}

func NewPlayer(id int, gameboard [][]string, name string) *Player {
	player := &Player{
		Name:           name,
		ID:             id,
		Speed:          300,
		AvailableBombs: 1,
		FireDistance:   2,
		Lives:          3,
		Damaged:        false,
	}

	if gameboard == nil {
		return player
	}

	rows := len(gameboard)
	cols := len(gameboard[0])

	switch id {
	case 1:
		player.X = 0
		player.Y = 0
		player.Direction = "down"
	case 2:
		player.X = cols - 1
		player.Y = 0
		player.Direction = "down"
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
