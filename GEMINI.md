# Futa Ticket Hunter

## Project Overview

Futa Ticket Hunter is a Go-based application designed to automate the process of hunting and booking bus tickets on the futabus.vn platform. 

The system is composed of two main services:
1.  **API Server (`cmd/api`)**: A web server that provides a REST API to manage ticket hunting schedules, view statistics, and proxy searches to the Futa API. It also serves a static frontend application.
2.  **Background Worker (`cmd/worker`)**: A background process that continuously polls the database for active schedules and interacts with the internal Futa API client (`internal/futa`) to search for available trips, check seat availability, and automatically book tickets.

**Key Technologies:**
*   **Language:** Go (1.25.4)
*   **Database:** PostgreSQL (via `github.com/jackc/pgx/v5`). Local development uses Docker. Production configuration points to CockroachDB.
*   **External Integrations:** 
    *   Resend API for email notifications (`github.com/resend/resend-go/v3`)
    *   Webhooks for triggering external workflows
    *   futabus.vn API for ticket searching and booking

## Directory Structure

*   `cmd/api/`: Entry point for the API web server. Contains the `main.go` and `static/` directory for the frontend.
*   `cmd/worker/`: Entry point for the background worker that processes booking schedules.
*   `internal/config/`: Configuration loading and parsing from `config.yaml`.
*   `internal/database/`: Database connection and repository layer for interacting with PostgreSQL.
*   `internal/futa/`: The core API client for futabus.vn, handling token extraction, endpoint requests, and booking logic.
*   `internal/email/` & `internal/webhook/`: Handlers for sending notifications.
*   `migrations/`: SQL migration files for initializing the database schema.
*   `ui-screenshot/`: Contains screenshots of the UI.
*   `k8s/`: Kubernetes deployment manifests.

## Building and Running

A `Makefile` is provided to simplify development and deployment tasks.

**Prerequisites:**
*   Go 1.25+
*   Docker & Docker Compose (for local database)

**Commands:**

*   **Start Local Database:**
    ```bash
    make db
    ```
    This spins up a PostgreSQL container using `docker-compose.yml` and automatically applies the migrations in the `/migrations` folder.

*   **Run the API Server:**
    ```bash
    make api
    ```
    Starts the API server. By default, it expects a `config.yaml` file in the root directory.

*   **Run the Worker:**
    ```bash
    make worker
    ```
    Starts the background worker to process ticket schedules.

*   **Build Binaries:**
    ```bash
    make build
    ```
    Compiles both the API and Worker binaries into the `bin/` directory.

*   **Clean:**
    ```bash
    make clean
    ```
    Removes the `bin/` directory.

## Development Conventions

*   **Architecture:** The project follows standard Go project layout conventions. The `cmd/` directory contains application entry points, and `internal/` contains private application and library code that cannot be imported by other projects.
*   **Configuration:** Configuration is managed via a `config.yaml` file. This file contains settings for the server port, database connection, worker polling interval, email provider keys, and Futa API constants. Be careful not to commit sensitive secrets to source control.
*   **Futa API Client (`internal/futa`)**: The application acts as a mobile client to the Futa API, extracting tokens from the web frontend and proxying them into API calls using specific headers (e.g., `X-Channel: mobile_app`, `User-Agent`). The `guide_line.md` file contains important details and expected payloads for interacting with Futa's undocumented APIs.
