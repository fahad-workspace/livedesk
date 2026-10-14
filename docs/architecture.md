# LiveDesk Architecture

## Design rules

1. The Narrative Writer never sees raw match events. It only sees verified moment cards. This is the anti-hallucination boundary.
2. Every claim in the output carries an evidence trail: event IDs and computed stats from Microsoft Fabric.
3. Agents recover from failure. The Fact Verifier triggers a rewrite loop on rejection (max two attempts, then a template fallback). The Match Analyst retries KQL queries on transient errors and degrades to cached features.
4. Synthetic data only. The simulator generates or replays football-realistic events; no real match data enters the system.

## Moment card (shared state)

The moment card is the unit of coordination between agents.

```json
{
  "id": "mc-2026-0001",
  "match_id": "match-001",
  "minute": 67,
  "type": "goal",
  "team": "home",
  "players": ["p-1042"],
  "summary": "Home team scores after a sustained pressure window",
  "evidence": {
    "event_ids": ["e-9921", "e-9924", "e-9930"],
    "stats": {
      "momentum_index": 0.81,
      "pressure_window_sec": 214,
      "shots_in_window": 4
    }
  },
  "status": "approved"
}
```

Status moves from `draft` to `approved` (Fact Verifier passed) or `rejected` (rewrite requested, with reasons attached).

## Personas

| Persona | Audience | Tone |
|---|---|---|
| studio | Broadcast producers and analysts | Tactical depth, feature references |
| broadcast | General viewers | Neutral, clear, paced |
| fan | Team supporters | Energetic, team-aware |
| kids | Young fans | Simple words, short sentences |

## Languages

English source, then Azure AI Translator with a custom football glossary: Hindi, Spanish, Arabic, French.
