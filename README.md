# AudioHunt

AudioHunt is a full-stack music recognition app inspired by Shazam. It lets a user enroll songs into a personal catalog, record a short live microphone clip in the browser, and identify the closest matching song using a custom audio fingerprinting pipeline.

Live site: [audiohuntnew.vercel.app](https://audiohuntnew.vercel.app/)

Repository: [github.com/Aggdaksh/AudioHunt-2](https://github.com/Aggdaksh/AudioHunt-2)

## What It Does

AudioHunt has two core flows:

- Library flow: upload a clean song file, add a title and artist, and store its audio fingerprints in PostgreSQL.
- Recognition flow: record music through the browser microphone, send the clip to the backend, fingerprint it, and compare it against enrolled tracks.

The result includes the matched song, artist, confidence score, and an estimated position inside the song.

## How It Works

AudioHunt does not call an external recognition API. The recognition logic is implemented in the Go backend.

1. A song is enrolled from the Library page.
2. The backend converts audio with FFmpeg into a consistent mono sample format.
3. The DSP pipeline creates a spectrogram, finds strong frequency peaks, and turns peak pairs into compact hashes.
4. These hashes are stored in PostgreSQL with the song metadata.
5. During recognition, the browser records a short microphone clip with `MediaRecorder`.
6. The backend normalizes and fingerprints the query clip.
7. The matcher compares query hashes against stored song hashes and picks the strongest aligned match.
8. The frontend displays the match result and saves it in local recognition history.

```mermaid
flowchart LR
  User["User plays music"] --> Browser["React frontend"]
  Browser --> Recorder["Browser MediaRecorder"]
  Recorder --> API["Go REST API"]
  API --> FFmpeg["FFmpeg audio normalization"]
  FFmpeg --> DSP["Spectrogram + peak fingerprinting"]
  DSP --> Matcher["Windowed hash matching"]
  Matcher --> DB[("PostgreSQL fingerprint catalog")]
  DB --> Matcher
  Matcher --> Result["Song, artist, confidence, timestamp"]
  Result --> Browser
```

## System Architecture

```mermaid
flowchart TB
  subgraph Client["Client Layer"]
    UI["React + TypeScript + Vite"]
    Mic["Microphone capture"]
    History["Local recognition history"]
  end

  subgraph API["Backend Layer"]
    Router["Go API with Gorilla Mux"]
    Songs["Song enrollment service"]
    Recognition["Recognition service"]
    Audio["FFmpeg + DSP pipeline"]
  end

  subgraph Data["Data Layer"]
    Postgres[("PostgreSQL")]
    SongRows["songs table"]
    Fingerprints["stored peak fingerprints"]
  end

  UI --> Mic
  UI --> Router
  Router --> Songs
  Router --> Recognition
  Songs --> Audio
  Recognition --> Audio
  Songs --> Postgres
  Recognition --> Postgres
  Postgres --> SongRows
  Postgres --> Fingerprints
```

## Tech Stack

Frontend:

- React 19 for the interactive single-page interface.
- TypeScript for safer component and API contracts.
- Vite for fast local development and optimized production builds.
- Axios for API requests.
- Browser `MediaRecorder` for microphone capture.
- Custom CSS for the Shazam-style visual experience.

Backend:

- Go for the REST API and audio-processing pipeline.
- Gorilla Mux for routing.
- FFmpeg for reliable audio conversion across MP3, WAV, M4A, WEBM, and other common formats.
- Custom DSP code for spectrogram generation, peak detection, and hash matching.
- PostgreSQL for durable song metadata and fingerprint storage.
- Docker for reproducible backend packaging with FFmpeg available in the runtime.

Database:

- PostgreSQL stores songs, metadata, segment fingerprints, and peak fingerprints.
- Fingerprints are persisted so the recognizer can survive backend restarts without requiring songs to be reprocessed immediately.

## Why This Stack

React + Vite was chosen because the frontend needs a responsive browser experience with microphone access, fast iteration, and clean production output. TypeScript keeps the UI state machine, API responses, and recognition result objects easier to maintain.

Go was chosen for the backend because audio fingerprinting is CPU-heavy and benefits from a compiled, simple, concurrent runtime. Go also keeps the API small, fast, and easy to package as a single service.

PostgreSQL was chosen because fingerprints are structured data that need durable storage, indexing, and predictable querying. A relational database also makes it easy to keep song metadata and fingerprint data consistent.

FFmpeg was chosen because browser and uploaded audio files can arrive in many formats. Converting everything into a consistent sample rate and channel layout makes the fingerprinting pipeline more predictable.

## Scalability And System Design

The current system is intentionally split into frontend, backend, and database layers. That separation makes the project easier to scale without rewriting the whole app.

- The frontend is static and can be served globally through a CDN.
- The backend is stateless for API requests, so multiple Go instances can run behind a load balancer.
- PostgreSQL is the source of truth for the song catalog and fingerprints.
- The backend keeps an in-memory cache of loaded song fingerprints for faster repeated recognition.
- Fingerprints are persisted, which reduces dependence on local uploaded files after restart.
- Recognition and enrollment are separate flows, so enrollment can later move to a background worker without changing the user-facing recognition API.
- Large audio files can later be stored in object storage such as S3 or Supabase Storage, while PostgreSQL keeps metadata and fingerprints.
- For bigger catalogs, fingerprint hashes can be indexed or moved into a search-optimized store while keeping the same API shape.

Future scaling direction:

- Add a job queue for long-running song enrollment.
- Store original audio files in object storage.
- Add worker services for CPU-heavy fingerprint extraction.
- Keep the Go API focused on fast recognition requests.
- Add rate limiting and authentication for public usage.
- Add observability around enrollment time, recognition latency, and match confidence.

## Project Structure

```text
AudioHunt/
├── backend/
│   ├── cmd/api/                 # Go API entry point
│   ├── internal/api/            # Router, CORS, health endpoint
│   ├── internal/audio/          # FFmpeg conversion and DSP pipeline
│   ├── internal/db/             # PostgreSQL access
│   ├── internal/matching/       # Hash matching logic
│   ├── internal/songs/          # Enrollment and recognition services
│   ├── migrations/              # PostgreSQL schema
│   ├── uploads/                 # Local development uploads
│   ├── Dockerfile               # Backend image with FFmpeg
│   └── go.mod
├── frontend/
│   ├── public/                  # Logo and favicon assets
│   ├── src/components/          # UI components
│   ├── src/hooks/               # Recording, waveform, history logic
│   ├── src/pages/               # Recognize, History, Library pages
│   ├── src/services/api.ts      # API client
│   ├── src/App.tsx
│   └── package.json
├── docker-compose.yml           # Local PostgreSQL database
├── render.yaml                  # Backend service blueprint
└── README.md
```

## Run Locally

### Prerequisites

- Go 1.25+
- Node.js and npm
- Docker Desktop
- FFmpeg, if running the backend directly outside Docker

### 1. Clone the project

```bash
git clone https://github.com/Aggdaksh/AudioHunt-2.git
cd AudioHunt-2
```

### 2. Start PostgreSQL

```bash
docker compose up -d
```

The local database runs on port `5434` and automatically loads the SQL files in `backend/migrations`.

### 3. Configure and run the backend

```bash
cd backend
cp .env.example .env
go mod download
go run ./cmd/api
```

Expected backend URL:

```text
http://localhost:8080
```

The local backend `.env` should contain:

```env
PORT=8080
DB_DSN=postgres://shazam_user:pass123@localhost:5434/shazam_clone?sslmode=disable
```

### 4. Configure and run the frontend

Open a second terminal:

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Expected frontend URL:

```text
http://localhost:5173
```

The local frontend `.env` should contain:

```env
VITE_API_URL=http://localhost:8080
```

### 5. Try the app

1. Open the Library tab.
2. Upload a song and enter a clean title and artist.
3. Open the Recognize tab.
4. Play that enrolled song loudly near the mic for 8 to 10 seconds.
5. Press Stop and wait for the match result.

## API Overview

```text
GET  /api/health      Check backend status
GET  /api/songs       List enrolled songs
POST /api/songs       Enroll a song file with title and artist
POST /api/recognize   Recognize a microphone recording
```

## Current Limitations

- Enrollment can take time for full-length songs because the backend extracts fingerprints from the uploaded audio.
- Free hosted backend instances may sleep after inactivity, so the first request can be slower.
- Recognition quality depends on microphone quality, speaker distance, background noise, and whether the song has already been enrolled.
- This is a portfolio/educational project, so only use audio that you are allowed to upload and demo.

## Author

Made by Daksh Aggarwal.

- Email: [daksh121105@gmail.com](mailto:daksh121105@gmail.com)
- LinkedIn: [daksh-aggarwal-a5938028b](https://www.linkedin.com/in/daksh-aggarwal-a5938028b/)

## License

This project is for educational and portfolio use.
