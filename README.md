# Football League Simulation

A Go-based football league simulation system with a Vue.js frontend.

## Features

- League table with team statistics
- Match simulation
- Week-by-week match results
- RESTful API endpoints
- Modern Vue.js frontend

## Prerequisites

- Go 1.16 or higher
- Node.js 14 or higher
- npm 6 or higher

## Backend Setup

1. Navigate to the project root directory
2. Install Go dependencies:
   ```bash
   go mod download
   ```
3. Run the backend server:
   ```bash
   go run cmd/main/main.go
   ```
   The server will start on http://localhost:8080

## Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```
   The frontend will be available at http://localhost:5173

## API Endpoints

- `GET /api/teams` - Get all teams
- `GET /api/matches` - Get all matches
- `GET /api/league` - Get league statistics
- `POST /api/matches/simulate/{week}` - Simulate matches for a specific week
- `POST /api/matches/simulate-all` - Simulate all remaining matches
- `PUT /api/matches/{id}` - Update match result

## Project Structure

```
.
├── cmd/
│   └── main/
│       └── main.go
├── internal/
│   ├── database/
│   │   └── sqlite.go
│   ├── handlers/
│   │   └── api.go
│   └── models/
│       ├── league.go
│       ├── match.go
│       └── team.go
├── frontend/
│   ├── src/
│   │   ├── App.vue
│   │   └── main.js
│   └── package.json
└── README.md
```

## Database Schema

The application uses SQLite for data storage. The schema includes:

- Teams table
- Matches table
- League standings table

## Running Tests

```bash
go test ./...
```
