---
name: Work Packet (Copilot Coding Agent)
about: Atomic task for Copilot coding agent (one PR)
title: "WPXX – <short imperative>"
labels: ["tui", "work-packet", "ready", "copilot"]
---

## Objective
One sentence describing the concrete outcome.

## Files to create/modify (only these)
- tui/...

## Requirements (MUST)
- Bullet list of MUST requirements only.
- No “maybe”, no “nice to have”.

## Validation
Run:
- (from tui/) `go test ./...`

## Done when
- Tests pass
- Requirements met
- Godoc added for exported symbols

## Stop condition
Do not continue to other tasks. Stop when done.