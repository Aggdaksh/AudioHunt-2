package songs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{service: NewService()}
}

func RegisterRoutes(r *mux.Router) {
	h := NewHandler()
	r.HandleFunc("/api/songs", h.AddSong).Methods("POST")
	r.HandleFunc("/api/songs", h.ListSongs).Methods("GET")
	r.HandleFunc("/api/recognize", h.Recognize).Methods("POST")
}

func (h *Handler) ListSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := h.service.ListSongs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type songItem struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Artist string `json:"artist"`
	}
	out := make([]songItem, len(songs))
	for i, s := range songs {
		out[i] = songItem{ID: s.ID, Title: s.Title, Artist: s.Artist}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "File too large or form parse error", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio_file")
	if err != nil {
		http.Error(w, "Missing audio_file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	artist := r.FormValue("artist")
	if title == "" || artist == "" {
		http.Error(w, "title and artist are required", http.StatusBadRequest)
		return
	}

	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	ext := filepath.Ext(header.Filename)
	base := header.Filename[:len(header.Filename)-len(ext)]
	uniqueFilename := fmt.Sprintf("%s_%d%s", base, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	songID, err := h.service.AddSong(title, artist, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      songID,
		"message": "Song added successfully",
	})
}

func (h *Handler) Recognize(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio_snippet")
	if err != nil {
		http.Error(w, "Missing audio_snippet", http.StatusBadRequest)
		return
	}
	defer file.Close()

	queryDir := "./uploads/queries"
	if err := os.MkdirAll(queryDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to prepare query directory", http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".wav"
	}
	queryPath := filepath.Join(queryDir, fmt.Sprintf("query_%d%s", time.Now().Unix(), ext))

	tmpFile, err := os.Create(queryPath)
	if err != nil {
		http.Error(w, "Failed to create query file", http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(w, "Failed to save query audio", http.StatusInternalServerError)
		return
	}
	if err := tmpFile.Sync(); err != nil {
		http.Error(w, "Failed to finalize query audio", http.StatusInternalServerError)
		return
	}

	log.Printf("Saved recognition query to %s", queryPath)

	result, err := h.service.Recognize(queryPath)
	if err != nil {
		log.Printf("!!! RECOGNIZE ERROR: %v", err)
		msg := err.Error()
		if msg == "no fingerprints extracted — audio too short or too quiet" {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"song_id":      result.Song.ID,
		"title":        result.Song.Title,
		"artist":       result.Song.Artist,
		"confidence":   result.Confidence,
		"is_match":     result.IsMatch,
		"time_in_song": result.TimeInSong,
	})
}
