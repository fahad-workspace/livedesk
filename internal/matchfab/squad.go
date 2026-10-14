package matchfab

import (
	"fmt"
	"math/rand"
)

// Player is a squad member. Match state fields (X, Y, DistM) are internal and
// never serialised into events.
type Player struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	TeamID   string  `json:"team_id"`
	Number   int     `json:"number"`
	Position string  `json:"position"`
	X        float64 `json:"-"`
	Y        float64 `json:"-"`
	DistM    float64 `json:"-"`
}

// Team is an eleven player squad in a 4-4-2 shape.
type Team struct {
	ID      string
	Name    string
	Players []*Player
}

var firstNames = []string{
	"James", "Jack", "Harry", "Adam", "Daniel", "Mohamed", "Kevin",
	"Bruno", "Marcus", "Bukayo", "Declan", "Phil", "John", "Kyle",
	"Luke", "Aaron", "Jordan", "Ben", "Kieran", "Cole",
}

var lastNames = []string{
	"Carter", "Walker", "Stone", "Reid", "Bishop", "Hughes", "Morgan",
	"Bell", "Ward", "Brooks", "Hayes", "Price", "Wood", "Fox", "Lane",
	"Shaw", "Grant", "Hart", "Boyd", "Nash",
}

// newTeam builds a squad. Home teams set up on the left of the pitch and
// attack towards x = 105; away teams mirror them.
func newTeam(id, name string, home bool, rng *rand.Rand) *Team {
	type slot struct {
		pos string
		x   float64
	}
	slots := []slot{
		{"GK", 5},
		{"DF", 25}, {"DF", 25}, {"DF", 25}, {"DF", 25},
		{"MF", 50}, {"MF", 50}, {"MF", 50}, {"MF", 50},
		{"FW", 72}, {"FW", 72},
	}
	ys := []float64{34, 12, 28, 40, 56, 20, 30, 38, 48, 26, 42}
	used := map[string]bool{}
	players := make([]*Player, 0, len(slots))
	for i, s := range slots {
		x := s.x
		if !home {
			x = PitchLength - x
		}
		players = append(players, &Player{
			ID:       fmt.Sprintf("%s-p%02d", id, i+1),
			Name:     pickName(rng, used),
			TeamID:   id,
			Number:   i + 1,
			Position: s.pos,
			X:        x,
			Y:        ys[i],
		})
	}
	return &Team{ID: id, Name: name, Players: players}
}

func pickName(rng *rand.Rand, used map[string]bool) string {
	for {
		n := firstNames[rng.Intn(len(firstNames))] + " " + lastNames[rng.Intn(len(lastNames))]
		if !used[n] {
			used[n] = true
			return n
		}
	}
}
