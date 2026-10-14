package matchfab

// Event is a single synthetic match event, serialised as one JSON object per
// line (NDJSON). Coordinates are in metres on a 105 x 68 pitch, with the home
// team attacking towards x = 105 in the first period.
type Event struct {
	EventID      string         `json:"event_id"`
	MatchID      string         `json:"match_id"`
	Seq          int            `json:"seq"`
	Period       int            `json:"period"`
	ClockMs      int64          `json:"clock_ms"`
	Type         string         `json:"type"`
	TeamID       string         `json:"team_id,omitempty"`
	PlayerID     string         `json:"player_id,omitempty"`
	X            float64        `json:"x,omitempty"`
	Y            float64        `json:"y,omitempty"`
	EndX         float64        `json:"end_x,omitempty"`
	EndY         float64        `json:"end_y,omitempty"`
	BallSpeedKmh float64        `json:"ball_speed_kmh,omitempty"`
	Attrs        map[string]any `json:"attrs,omitempty"`
}
