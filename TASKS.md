# TASKS.md: Build Queue

Ordered build queue. Each task is one GitHub issue and one pull request. Start the next task only after the previous one merges. Tasks marked [azure] need the Azure subscription configured first.

## Phase 1: Data and plumbing

### Task 1: Event Hubs producer sink for MatchFab [azure]
Add a sink to `cmd/matchfab` that streams NDJSON events to Azure Event Hubs in batches. Config from env vars (`EVENTHUB_CONNECTION_STRING`, `EVENTHUB_NAME`). Stdout stays the default sink.
Done when: `go run ./cmd/matchfab -sink eventhubs` delivers a full match to a real hub, verified with the Azure portal metrics.

### Task 2: Fabric Eventhouse schema and ingestion [azure]
KQL scripts under `fabric/` that create the events table (matching `internal/matchfab/event.go`), plus the ingestion mapping from Event Hubs. Include a `fabric/README.md` with the exact portal steps.
Done when: a replayed match from Task 1 is queryable in the Eventhouse within seconds of ingestion.

### Task 3: KQL feature queries
`fabric/features.kql` with materialized views or functions for: momentum index (rolling shot and pressure balance), pressure windows, pass accuracy and difficulty by team and player, player speed and distance aggregates, possession phases.
Done when: each query runs against Task 2 data and returns sane results, documented with sample output in comments.

## Phase 2: Shared state and agents

### Task 4: Moment card store
Go service under `cmd/cards` implementing the moment card schema from `docs/architecture.md`: create, read, update status (draft, approved, rejected with reasons). In-memory plus JSON file persistence is enough for the hackathon.
Done when: unit tests cover the status lifecycle and `go test ./...` is green.

### Task 5: Match Analyst agent [azure]
Foundry agent (Microsoft Agent Framework) with tools that call the KQL features from Task 3. Watches the event stream, detects key moments (goals, pressure swings, milestones, sprint and distance thresholds), writes moment cards with evidence.
Done when: running it against a replayed match produces cards for every goal and at least three non-goal moments, each with event IDs and stats.

### Task 6: Narrative Writer agent [azure]
Foundry agent that turns approved moment cards into commentary in five modes: studio, broadcast, fan, player focus, kids. Sees only moment cards, never raw events.
Done when: each persona produces visibly different text for the same card, shown in a test or sample output file.

### Task 7: Fact Verifier agent and rewrite loop [azure]
Foundry agent that checks each generated sentence against the card evidence. Rejects with reasons, triggers rewrite (max 2), then template fallback.
Done when: a deliberately corrupted claim in a test gets rejected, rewritten and approved, with the loop visible in logs.

### Task 8: Localizer agent [azure]
Azure AI Translator integration with a football glossary (`localization/glossary.json`). Five languages: Hindi, Spanish, Arabic, French plus English source.
Done when: one approved narrative exists in all five languages with glossary terms preserved, committed as a sample under `localization/samples/`.

## Phase 3: Delivery

### Task 9: Orchestrator and SSE API [azure]
Go service under `cmd/orchestrator`: runs the agent pipeline per moment card, streams approved narratives over SSE (`/stream?persona=&lang=`). Deployable to Azure Container Apps with a Dockerfile.
Done when: `curl -N` against the running service shows narratives arriving during a MatchFab replay.

### Task 10: Studio console [azure]
React app under `web/`: live ticker (SSE), momentum graph, Why evidence drawer, persona and language switchers, recap tab. Deploys to Azure Static Web Apps.
Done when: the console renders a full replay live against the Task 9 API with no console errors.

### Task 11: Overlay JSON output
Timed cue sheet endpoint (`/overlay`) emitting machine-readable overlay cues: clock, type (stat card, momentum strip, player tag), payload, display duration.
Done when: one replayed match produces a valid cue sheet, validated by a JSON schema test.

### Task 12: Recap generator
Half-time and full-time recap from approved moment cards, per persona.
Done when: a full match produces both recaps in at least studio and kids modes, committed as samples.

## Phase 4: Polish and submission

### Task 13: Demo seed and replay script
Pick the most dramatic seed (late winner, momentum swings). `scripts/demo.sh` runs the whole flow end to end in under 3 minutes.
Done when: a fresh clone plus documented Azure setup runs the demo with one script.

### Task 14: README and architecture diagram
Final README pass: architecture diagram (committed image), stack list, judge-focused quickstart.
Done when: a stranger can understand and run the project from the README alone.

### Task 15: Video shot list
`docs/video-shotlist.md`: second-by-second plan for the under 2 minute demo video following the agreed structure.
Done when: reviewed by the project owner.
