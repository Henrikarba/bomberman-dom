package game

import (
	"time"
)

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

func NewPlayer(name string, id int) *Player {
	player := &Player{
		ID:             id,
		Name:           name,
		Speed:          500,
		AvailableBombs: 1,
		FireDistance:   1,
		Lives:          3,
		Damaged:        false,
	}

	return player
}

func StartingPositions(player *Player) *Player {
	rows := 13
	cols := 19
	player.Lives = 3
	switch player.ID {
	case 1:
		player.X = 0
		player.Y = 0
		player.Direction = "down"
	case 2:
		player.X = cols - 1
		player.Y = rows - 1
		player.Direction = "up"
	case 3:
		player.X = cols - 1
		player.Y = 0
		player.Direction = "down"
	case 4:
		player.X = 0
		player.Y = rows - 1
		player.Direction = "up"
	}
	return player
}
