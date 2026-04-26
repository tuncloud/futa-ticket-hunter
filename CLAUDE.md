# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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