# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.

## Commands
- Run Database (PostgreSQL via Docker): `make db`
- Run API server: `make api` (or `go run ./cmd/api/`)
- Run background worker: `make worker` (or `go run ./cmd/worker/`)
- Build binaries (creates `bin/api` and `bin/worker`): `make build`
- Clean build outputs: `make clean`
- Run all tests: `go test ./...`
- Run a single test: `go test -v -run <TestName> <package_path>`

## High-Level Architecture
This repository is a Go monorepo that manages an automated bus ticket booking system. The application is split into two main executables and several shared internal packages:

### Executables
1. **API Server (`cmd/api/main.go`)**: 
   - A standard Go `net/http` web server providing REST endpoints (under `/api/*`) for creating/managing booking schedules, checking payment status, and retrieving stats.
   - It also embeds and serves a static Single Page Application (SPA) frontend from the `static/` directory using `//go:embed`.
2. **Background Worker (`cmd/worker/main.go`)**: 
   - A long-running process that periodically polls the database for `pending` booking schedules.
   - It orchestrates the complex logic of searching for routes, finding available trips, filtering by user preferences (time range, seat type, floor, window preference), and executing the reservation via the FUTA API.

### Core Internal Packages (`internal/`)
- **`futa/`**: The core API client for communicating with the FUTA bus booking system. Handles tasks like searching pickup points, routes, trips, retrieving seat diagrams, and executing the booking.
- **`database/`**: Handles all PostgreSQL operations via `pgx/v5`. Manages the state of schedules (`pending`, `searching`, `success`, `paid`, etc.) and system stats.
- **`email/`**: Sends email notifications (like payment links after successful bookings) using the Resend API.
- **`webhook/`**: Handles webhook callbacks to notify external systems of booking updates.
- **`config/`**: Manages configuration loaded from `config.yaml`.