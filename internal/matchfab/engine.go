package matchfab

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	PitchLength = 105.0
	PitchWidth  = 68.0
	halfMs      = int64(45 * 60 * 1000)
	fullMs      = int64(90 * 60 * 1000)
)

// Event types emitted by the engine.
const (
	EventMatchStart       = "match_start"
	EventHalfEnd          = "half_end"
	EventHalfStart        = "half_start"
	EventMatchEnd         = "match_end"
	EventPass             = "pass"
	EventShot             = "shot"
	EventGoal             = "goal"
	EventTackle           = "tackle"
	EventPossessionChange = "possession_change"
	EventPressure         = "pressure"
	EventFoul             = "foul"
	EventCard             = "card"
	EventPlayerSpeed      = "player_speed"
	EventPlayerDistance   = "player_distance"
)

// Engine simulates one match. The same seed always produces the same match,
// which keeps demos and tests reproducible.
type Engine struct {
	matchID string
	rng     *rand.Rand
	home    *Team
	away    *Team

	seq     int
	clockMs int64
	period  int
	poss    *Player

	scoreHome int
	scoreAway int

	nextDistMark map[*Player]float64
}

// NewEngine builds a match between two synthetic clubs.
func NewEngine(matchID string, rng *rand.Rand) *Engine {
	return &Engine{
		matchID:      matchID,
		rng:          rng,
		home:         newTeam("home", "Northbridge", true, rng),
		away:         newTeam("away", "Southport", false, rng),
		period:       1,
		nextDistMark: map[*Player]float64{},
	}
}

// Run plays the match to full time, emitting every event in clock order.
func (e *Engine) Run(out func(Event)) {
	e.emitMatchStart(out)
	e.poss = e.midfielder(e.home)
	for e.clockMs < fullMs {
		if e.period == 1 && e.clockMs >= halfMs {
			e.clockMs = halfMs
			e.emit(out, Event{Type: EventHalfEnd})
			e.period = 2
			e.emit(out, Event{Type: EventHalfStart})
			e.poss = e.midfielder(e.away)
		}
		e.step(out)
	}
	e.clockMs = fullMs
	e.emit(out, Event{Type: EventMatchEnd, Attrs: map[string]any{
		"score_home": e.scoreHome,
		"score_away": e.scoreAway,
	}})
}

func (e *Engine) step(out func(Event)) {
	team := e.teamOf(e.poss)
	roll := e.rng.Float64()
	switch {
	case e.attackingThird(e.poss, team) && roll < 0.18:
		e.shoot(out, team)
	case roll < 0.30:
		e.carry(out, team)
	default:
		e.pass(out, team)
	}
	if e.teamOf(e.poss) == team && e.attackingThird(e.poss, team) && e.rng.Float64() < 0.35 {
		e.emitPressure(out, team)
	}
}

func (e *Engine) pass(out func(Event), team *Team) {
	from := e.poss
	to := e.pickPassTarget(from, team)
	dist := math.Hypot(to.X-from.X, to.Y-from.Y)
	difficulty := clamp(dist/60, 0, 1)*0.7 + e.pressureNear(from)*0.3
	speed := 40 + e.rng.Float64()*45
	success := e.rng.Float64() < 0.97-0.5*difficulty

	outcome := "complete"
	if !success {
		outcome = "incomplete"
	}
	e.emit(out, Event{
		Type: EventPass, TeamID: team.ID, PlayerID: from.ID,
		X: from.X, Y: from.Y, EndX: to.X, EndY: to.Y,
		BallSpeedKmh: round1(speed),
		Attrs: map[string]any{
			"outcome":    outcome,
			"distance_m": round1(dist),
			"difficulty": round2(difficulty),
			"target_id":  to.ID,
		},
	})
	e.clockMs += int64(1500 + e.rng.Intn(2500))

	if success {
		e.moveTeam(team)
		e.addDistance(from, dist*0.3)
		e.poss = to
		return
	}
	defender := e.nearestOpponent(to)
	e.resolveTurnover(out, team, defender)
}

func (e *Engine) carry(out func(Event), team *Team) {
	p := e.poss
	dir := e.attackDir(team)
	dist := 4 + e.rng.Float64()*8
	speedMs := 4.5 + e.rng.Float64()*4
	oldX, oldY := p.X, p.Y
	p.X = clamp(p.X+dir*dist, 2, PitchLength-2)
	p.Y = clamp(p.Y+(e.rng.Float64()-0.5)*6, 2, PitchWidth-2)
	actual := math.Hypot(p.X-oldX, p.Y-oldY)
	e.clockMs += int64(actual / speedMs * 1000)
	e.addDistance(p, actual)

	speedKmh := speedMs * 3.6
	if speedKmh > 29 {
		e.emit(out, Event{Type: EventPlayerSpeed, TeamID: team.ID, PlayerID: p.ID,
			X: p.X, Y: p.Y,
			Attrs: map[string]any{"speed_kmh": round1(speedKmh), "threshold": "sprint"}})
	}

	mark := e.nextDistMark[p]
	if mark == 0 {
		mark = 5000
	}
	for p.DistM >= mark {
		e.emit(out, Event{Type: EventPlayerDistance, TeamID: team.ID, PlayerID: p.ID,
			Attrs: map[string]any{"distance_m": mark, "threshold": "milestone"}})
		mark += 5000
		e.clockMs += 200
	}
	e.nextDistMark[p] = mark
}

func (e *Engine) shoot(out func(Event), team *Team) {
	p := e.poss
	goalX := PitchLength
	if e.attackDir(team) < 0 {
		goalX = 0
	}
	distGoal := math.Abs(goalX - p.X)
	speed := 80 + e.rng.Float64()*50
	onTarget := e.rng.Float64() < 0.55
	goalProb := clamp(0.45-distGoal/90, 0.05, 0.40)
	scored := onTarget && e.rng.Float64() < goalProb

	outcome := "off_target"
	if onTarget {
		outcome = "on_target"
	}
	if scored {
		outcome = "goal"
	}
	e.emit(out, Event{Type: EventShot, TeamID: team.ID, PlayerID: p.ID,
		X: p.X, Y: p.Y, BallSpeedKmh: round1(speed),
		Attrs: map[string]any{"outcome": outcome, "distance_to_goal_m": round1(distGoal)}})
	e.clockMs += 2000

	if scored {
		if team == e.home {
			e.scoreHome++
		} else {
			e.scoreAway++
		}
		e.emit(out, Event{Type: EventGoal, TeamID: team.ID, PlayerID: p.ID,
			X: p.X, Y: p.Y,
			Attrs: map[string]any{
				"score_home": e.scoreHome,
				"score_away": e.scoreAway,
			}})
		e.clockMs += 20000 // celebration and restart
		e.poss = e.midfielder(e.opponent(team))
		return
	}
	// Save or miss: the defending keeper collects or takes the goal kick.
	e.poss = e.opponent(team).Players[0]
}

func (e *Engine) resolveTurnover(out func(Event), losing *Team, defender *Player) {
	winning := e.opponent(losing)
	e.emit(out, Event{Type: EventTackle, TeamID: winning.ID, PlayerID: defender.ID,
		X: defender.X, Y: defender.Y,
		Attrs: map[string]any{"outcome": "won"}})
	e.clockMs += 1200

	if e.rng.Float64() < 0.12 {
		e.emit(out, Event{Type: EventFoul, TeamID: winning.ID, PlayerID: defender.ID,
			X: defender.X, Y: defender.Y,
			Attrs: map[string]any{"card": "none"}})
		e.clockMs += 8000
		if e.rng.Float64() < 0.25 {
			card := "yellow"
			if e.rng.Float64() < 0.04 {
				card = "red"
			}
			e.emit(out, Event{Type: EventCard, TeamID: winning.ID, PlayerID: defender.ID,
				Attrs: map[string]any{"card": card}})
			e.clockMs += 5000
		}
	}

	e.emit(out, Event{Type: EventPossessionChange, TeamID: winning.ID,
		Attrs: map[string]any{"from_team": losing.ID, "to_team": winning.ID}})
	e.poss = defender
}

func (e *Engine) emitPressure(out func(Event), team *Team) {
	e.emit(out, Event{Type: EventPressure, TeamID: team.ID,
		X: e.poss.X, Y: e.poss.Y,
		Attrs: map[string]any{
			"intensity": round2(0.5 + e.rng.Float64()*0.5),
			"zone":      "attacking_third",
		}})
	e.clockMs += 1000
}

func (e *Engine) emitMatchStart(out func(Event)) {
	e.emit(out, Event{Type: EventMatchStart, Attrs: map[string]any{
		"home_team": e.home,
		"away_team": e.away,
		"formation": "4-4-2",
	}})
}

func (e *Engine) emit(out func(Event), ev Event) {
	e.seq++
	ev.EventID = fmt.Sprintf("e-%06d", e.seq)
	ev.MatchID = e.matchID
	ev.Seq = e.seq
	ev.Period = e.period
	ev.ClockMs = e.clockMs
	out(ev)
}

func (e *Engine) attackDir(t *Team) float64 {
	homePositive := e.period == 1
	if t == e.home {
		if homePositive {
			return 1
		}
		return -1
	}
	if homePositive {
		return -1
	}
	return 1
}

func (e *Engine) attackingThird(p *Player, t *Team) bool {
	if e.attackDir(t) > 0 {
		return p.X > 70
	}
	return p.X < 35
}

func (e *Engine) teamOf(p *Player) *Team {
	if p.TeamID == e.home.ID {
		return e.home
	}
	return e.away
}

func (e *Engine) opponent(t *Team) *Team {
	if t == e.home {
		return e.away
	}
	return e.home
}

func (e *Engine) midfielder(t *Team) *Player {
	return t.Players[6]
}

func (e *Engine) pickPassTarget(from *Player, team *Team) *Player {
	dir := e.attackDir(team)
	var best *Player
	bestScore := math.Inf(-1)
	for _, p := range team.Players {
		if p == from {
			continue
		}
		d := math.Hypot(p.X-from.X, p.Y-from.Y)
		forward := dir * (p.X - from.X) / 30
		score := 1/(d+5) + forward*0.4 + e.rng.Float64()*0.3
		if score > bestScore {
			bestScore = score
			best = p
		}
	}
	return best
}

func (e *Engine) nearestOpponent(to *Player) *Player {
	opp := e.opponent(e.teamOf(to))
	var best *Player
	bestDist := math.Inf(1)
	for _, p := range opp.Players {
		d := math.Hypot(p.X-to.X, p.Y-to.Y)
		if d < bestDist {
			bestDist = d
			best = p
		}
	}
	return best
}

func (e *Engine) pressureNear(p *Player) float64 {
	opp := e.opponent(e.teamOf(p))
	near := 0
	for _, o := range opp.Players {
		if math.Hypot(o.X-p.X, o.Y-p.Y) < 8 {
			near++
		}
	}
	return clamp(float64(near)/3, 0, 1)
}

func (e *Engine) moveTeam(team *Team) {
	dir := e.attackDir(team)
	for _, p := range team.Players {
		p.X = clamp(p.X+dir*(1+e.rng.Float64()*3), 2, PitchLength-2)
		p.Y = clamp(p.Y+(e.rng.Float64()-0.5)*2, 2, PitchWidth-2)
	}
}

func (e *Engine) addDistance(p *Player, m float64) {
	p.DistM += m
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func round1(v float64) float64 { return math.Round(v*10) / 10 }
func round2(v float64) float64 { return math.Round(v*100) / 100 }
