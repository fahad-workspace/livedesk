# LiveDesk: Explainable AI Match Intelligence for Live Football

Real-time match narratives that prove every claim. A multi-agent system turns synthetic football events into explainable, multilingual commentary, recaps and tactical insight for studios, broadcasters and fans.

Built for the Microsoft x Premier League **Inside the Game: Developer Hackathon** (October 2026). Synthetic data only. No real Premier League match data is used.

## What it does

LiveDesk is a live broadcast desk powered by AI agents. It takes in a real-time stream of synthetic football match events and turns them into narratives that never invent a fact. Every generated sentence links back to the exact match events and stats behind it, so a producer can always ask: why do you say that?

## How it works

Four agents work together under Microsoft Agent Framework on Microsoft Foundry. They share a common state built on moment cards.

| Agent | Role |
|---|---|
| Match Analyst | Queries Microsoft Fabric (Eventhouse, KQL) to find key moments and compute tactical features such as momentum and pressure |
| Narrative Writer | Writes commentary in four personas: studio analyst, neutral broadcast, fan mode and kids |
| Fact Verifier | Checks every sentence against the underlying match events and rejects anything it cannot support, which sends the text back for a rewrite |
| Localizer | Translates approved output into five languages using a football glossary |

## Architecture

```
Match Simulator (Go)
        |
        v
Azure Event Hubs
        |
        v
Microsoft Fabric Eventhouse (KQL)  <---- Match Analyst agent (tool calls)
        |                                   |
        | moment cards (shared state)       |
        v                                   v
Narrative Writer  <== rewrite loop ==>  Fact Verifier
        |
        v
Localizer (Azure AI Translator, 5 languages, football glossary)
        |
        v
SSE stream -> Studio Console (React on Azure Static Web Apps)
```

The studio console streams live commentary, shows a live momentum graph and produces half-time and full-time recaps on its own.

## Stack

- **Simulator and services:** Go
- **Ingestion:** Azure Event Hubs
- **Real-time analytics:** Microsoft Fabric Eventhouse, KQL
- **Agent orchestration:** Microsoft Agent Framework on Microsoft Foundry
- **Generation:** Azure OpenAI (via Foundry)
- **Translation:** Azure AI Translator
- **Delivery:** SSE, React, Azure Container Apps, Azure Static Web Apps
- **CI/CD:** GitHub Actions, built with GitHub Copilot

## Repo layout

- `cmd/simulator` - Go match simulator that replays a synthetic match timeline at configurable speed
- `docs/` - architecture and design notes

## Status

Hackathon in progress. Submissions close October 27, 2026.
