package game

type Movement struct {
	Type     string `json:"type"`
	Movement struct {
		Direction string `json:"direction"`
	} `json:"movement"`
}

type Player struct {
	ID        int    `json:"id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	Speed     int    `json:"speed,omitempty"`
}

type PositionUpdate struct {
	Type string   `json:"type"`
	Data []Player `json:"players"`
}
