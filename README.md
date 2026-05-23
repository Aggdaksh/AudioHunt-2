# 🎵 AudioHunt

A full-stack audio recognition project inspired by Shazam. AudioHunt lets users record a short clip in the browser, send it to a Go backend, and match it against enrolled songs using custom fingerprinting logic.

## ✨ Features

- **Audio Recognition Engine:** A Go backend that uses microphone-friendly preprocessing and persisted peak fingerprints for reliable matching.
- **Premium User Interface:** A modern, immersive glassmorphic UI with dark mode, smooth micro-animations, and dynamic visual feedback during audio recording.
- **Microphone Integration:** Seamless in-browser audio recording using the browser `MediaRecorder` pipeline.
- **Scalable Architecture:** Clean separation of concerns with a Go backend REST API, a React frontend, and a PostgreSQL database.

## 🛠 Tech Stack

### Frontend
- **Framework:** React 19 + TypeScript + Vite
- **Styling:** Custom CSS with glassmorphic and premium dark-mode aesthetics
- **Audio Capture:** RecordRTC for browser microphone access
- **HTTP Client:** Axios

### Backend
- **Language:** Go (Golang)
- **Routing:** Gorilla Mux
- **Audio Processing:** Segment-based fingerprinting for the live MVP, plus an experimental Shazam-style peak-pair DSP path that can be wired in later.
- **Database:** PostgreSQL

## 📂 Project Structure

```
.
├── backend/                  # Go backend application
│   ├── cmd/                  # Entry points (main.go)
│   ├── internal/             # Application code (API routing, DSP logic, song management)
│   ├── migrations/           # Database setup and schema files
│   ├── uploads/              # Local upload storage for development
│   ├── Dockerfile            # Render-ready backend image with ffmpeg
│   ├── .env                  # Backend environment variables
│   └── go.mod                # Go module dependencies
├── frontend/                 # React frontend application
│   ├── public/               # Static assets
│   ├── src/                  # React components (AudioRecorder, MatchResults, etc.)
│   ├── .env.example          # Frontend environment example
│   ├── package.json          # Node dependencies
│   └── vite.config.ts        # Vite configuration
├── docker-compose.yml        # Docker composition for the PostgreSQL database
├── render.yaml               # Render service blueprint
└── README.md                 # Project documentation
```

## 🚀 Getting Started

### Prerequisites

- [Go 1.25+](https://golang.org/doc/install)
- [Node.js](https://nodejs.org/) & npm
- [Docker Desktop](https://www.docker.com/products/docker-desktop) (for running PostgreSQL)

### 1. Start the Database
The backend relies on PostgreSQL. A `docker-compose.yml` file is provided to spin up the database easily with the correct schema migrations.

```bash
# From the root directory
docker-compose up -d
```
*Note: The database runs on port `5434` to avoid conflicts with local Postgres installations.*

### 2. Run the Go Backend

The backend server processes audio chunks and communicates with the database.

```bash
cd backend
go mod download
go run ./cmd/api/main.go
```
*The backend server will typically start on `http://localhost:8080`.*

### 3. Run the React Frontend

Open a new terminal window to start the Vite development server.

```bash
cd frontend
npm install
npm run dev
```
*The frontend will be accessible at `http://localhost:5173`. Open this in your browser to interact with the application.*

## 🌐 Deployment (Render + Vercel + Supabase)

This repo is now prepared for the stack you picked:

- **Frontend:** Vercel
- **Backend:** Render
- **Database:** Supabase Postgres

### 1. Create the Supabase database

Create a Supabase project, then copy the Postgres connection string and use it as `DB_DSN` on Render.

Run the SQL from:

- `backend/migrations/001_init.sql`
- `backend/migrations/002_hash_segments.sql`

### 2. Deploy the backend to Render

This repo includes:

- `backend/Dockerfile`
- `render.yaml`
- `backend/.env.example`

Why Docker on Render:

- the backend needs `ffmpeg`
- Docker makes the backend environment reproducible
- Render can build directly from the backend folder using the included blueprint

Required Render env vars:

- `PORT=8080`
- `DB_DSN=<your-supabase-connection-string>`

### 3. Deploy the frontend to Vercel

Use the `frontend` folder as the Vercel project root.

Required Vercel env var:

- `VITE_API_URL=https://<your-render-backend>.onrender.com`

The example file is:

- `frontend/.env.example`

### 4. Re-enroll demo songs after deploy

For a fresh cloud deployment, upload songs again through the Library page after the backend and database are live.

Important note:

- peak fingerprints are now persisted in Postgres, so recognition can survive free-host restarts better
- original uploaded audio files on a free backend host should still be treated as temporary

## ⚠️ Portfolio note

If this is a public resume project, use audio you are allowed to host and demo. For a portfolio deployment, royalty-free tracks or your own sample audio are the safest choice.

## 🧠 How it Works (Under the Hood)

1. **Recording:** The user clicks the record button on the frontend. `RecordRTC` captures the audio via the browser's MediaRecorder API and converts it into a valid audio format (like `.wav`).
2. **Transmission:** The chunk is sent to the backend `/api/recognize` endpoint via a `multipart/form-data` POST request.
3. **Signal Processing:** The Go backend fingerprints enrolled songs from clean source audio and fingerprints microphone queries with extra denoise and loudness normalization.
4. **Matching:** The backend compares the query against stored song fingerprints. For best results, play an enrolled song near the mic for 8-10 seconds.
5. **Results:** The backend returns the best-matching song with its confidence score, which the frontend displays to the user in a beautiful result card.

## 📝 License
This project is for educational and portfolio purposes.
