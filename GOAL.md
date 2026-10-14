# GOAL.md: Definition of Done

This file is the anchor for every build session. A task is done only when its acceptance criteria here and in TASKS.md are met, CI is green, and the change is merged to main through a pull request.

## The product

LiveDesk is a real-time, explainable, multilingual match intelligence system built for the Microsoft x Premier League Inside the Game Developer Hackathon. It turns synthetic football events (from MatchFab) into narratives that never invent a fact. Submissions close October 27, 2026. Code complete target: October 23, 2026.

## Done means all of this is true

### 1. Technological implementation (rubric, 20%)

- [ ] `go build ./...` and `go test ./...` pass in CI on every merge
- [ ] Web console builds in CI with no warnings treated as errors
- [ ] Every service runs from a container definition committed to this repo
- [ ] README documents setup, architecture and how to run the demo end to end

### 2. Agentic design and innovation (rubric, 20%)

- [ ] Four agents run under Microsoft Agent Framework on Microsoft Foundry: Match Analyst, Narrative Writer, Fact Verifier, Localizer
- [ ] Agents coordinate through a shared moment card store, not through prompts alone
- [ ] The Fact Verifier rewrite loop works: a rejected sentence is rewritten and re-verified, max two attempts, then template fallback
- [ ] The Match Analyst recovers from a failed KQL query (retry, then cached features)

### 3. Real-world impact (rubric, 20%)

- [ ] The whole system deploys to Azure with one documented command sequence
- [ ] A full synthetic match replays live in under 3 minutes with narratives streaming
- [ ] Half-time and full-time recaps generate automatically

### 4. User experience and presentation (rubric, 20%)

- [ ] Studio console shows: live commentary ticker, momentum graph, Why evidence drawer, persona switcher (studio, broadcast, fan, player focus, kids), language switcher (English, Hindi, Spanish, Arabic, French), recap tab
- [ ] Every narrative sentence links to its evidence (event IDs and stats)
- [ ] Output is also available as timed, machine-readable overlay JSON
- [ ] Demo video under 2 minutes, uploaded and linked on the project page

### 5. Category adherence (rubric, 20%)

- [ ] Only synthetic data from MatchFab is used anywhere in the system
- [ ] The five brief stages are all present: ingest, interpret, explain, render, personalize
- [ ] Brief-named features exist: pass difficulty rating, ball and shot speed, player speed and distance thresholds, multi-language storytelling

## Microsoft stack that must be genuinely used

Azure Event Hubs (ingestion), Microsoft Fabric Eventhouse and KQL (real-time analytics), Microsoft Foundry with Agent Framework (orchestration), Azure OpenAI via Foundry (generation), Azure AI Translator (localization), Azure Container Apps (services), Azure Static Web Apps (console), GitHub Actions (CI), GitHub Copilot (development).

## Hard rules for every agent working on this repo

1. Work on a branch named `task/<task-name>`. Never push to main.
2. Open a pull request when the task's acceptance criteria pass locally.
3. Do not add a service or dependency that is not in TASKS.md or this file.
4. Do not commit secrets. Azure config comes from environment variables.
5. If a task cannot meet its criteria, open a draft PR with a comment explaining the blocker instead of forcing a merge.
