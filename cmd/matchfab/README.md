# MatchFab

MatchFab is LiveDesk's synthetic football data engine. The hackathon redistributes no real match data, so MatchFab creates football-realistic event streams instead. It is both our ingestion source and our data innovation deliverable.

## Run it

```sh
# Full match to stdout, as fast as possible
go run ./cmd/matchfab

# Write a match to a file
go run ./cmd/matchfab -out match.ndjson

# Replay in near real time for the live demo (60 = one match minute per second)
go run ./cmd/matchfab -speed 60

# A different match every time
go run ./cmd/matchfab -seed 2026
```

The same seed always produces the same match, which keeps demos and tests reproducible.

## Output

NDJSON, one event per line. Coordinates are metres on a 105 x 68 pitch. The home team attacks towards x = 105 in the first period and swaps at half time.

```json
{"event_id":"e-000042","match_id":"match-001","seq":42,"period":1,"clock_ms":512400,"type":"pass","team_id":"home","player_id":"home-p07","x":52.1,"y":30.4,"end_x":68.3,"end_y":41.2,"ball_speed_kmh":61.4,"attrs":{"outcome":"complete","distance_m":19.4,"difficulty":0.23,"target_id":"home-p10"}}
```

## Event types

| Type | What it carries |
|---|---|
| `match_start` | Full lineups and formation metadata |
| `pass` | Outcome, distance, difficulty rating, ball speed, target |
| `shot` | Outcome, distance to goal, shot speed |
| `goal` | Scorer and running score |
| `tackle` | Outcome (won possession) |
| `possession_change` | From team, to team |
| `pressure` | Intensity (0 to 1) and zone |
| `foul` / `card` | Foul and any card shown |
| `player_speed` | Sprint indicator when a player crosses 29 km/h |
| `player_distance` | Distance milestones every 5 km |
| `half_end` / `half_start` / `match_end` | Match flow markers |

The event set mirrors the hackathon brief: passes, shots, tackles, possession changes, pressure events and match metadata, plus the called-out features (pass distance, accuracy, difficulty rating, ball and shot speed, player speed and distance thresholds).

## Roadmap

- [ ] Azure Event Hubs producer sink (env var config, batches)
- [ ] Configurable team styles (possession vs counter attack)
- [ ] Multi-match batch generation for Fabric load testing
