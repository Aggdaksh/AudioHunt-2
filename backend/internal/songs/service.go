package songs

import (
	"fmt"
	"sync"

	"Shazam-Vscode/backend/internal/audio"
	"Shazam-Vscode/backend/internal/audio/dsp"
	"Shazam-Vscode/backend/internal/audio/segment"
	"Shazam-Vscode/backend/internal/db"
	"Shazam-Vscode/backend/internal/matching"
)

type Service struct {
	mu               sync.RWMutex
	songs            []db.Song
	peakFingerprints map[int][]dsp.Fingerprint
}

func NewService() *Service {
	return &Service{
		peakFingerprints: make(map[int][]dsp.Fingerprint),
	}
}

func (s *Service) getSongs() ([]db.Song, error) {
	s.mu.RLock()
	if s.songs != nil {
		defer s.mu.RUnlock()
		return s.songs, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.songs != nil {
		return s.songs, nil
	}

	songs, err := db.GetSongsWithSegments()
	if err != nil {
		return nil, err
	}
	s.songs = songs
	return s.songs, nil
}

func (s *Service) ListSongs() ([]db.Song, error) {
	return s.getSongs()
}

func (s *Service) AddSong(title, artist, tempFilePath string) (int, error) {
	segmentFP, err := segment.FingerprintSongFile(tempFilePath)
	if err != nil {
		return 0, fmt.Errorf("segment fingerprint extraction failed: %w", err)
	}

	peakFPs, err := audio.ExtractFingerprints(tempFilePath, audio.DefaultConfig())
	if err != nil {
		return 0, fmt.Errorf("peak fingerprint extraction failed: %w", err)
	}
	if len(peakFPs) == 0 {
		return 0, fmt.Errorf("peak fingerprint extraction produced no hashes")
	}

	songID, err := db.InsertSongWithPeakFingerprints(
		title,
		artist,
		tempFilePath,
		segmentFP.Fingerprint,
		segmentFP.HashSegments,
		toDBPeakFingerprints(peakFPs),
	)
	if err != nil {
		return 0, fmt.Errorf("db insert song failed: %w", err)
	}

	s.mu.Lock()
	s.songs = nil
	s.peakFingerprints[songID] = peakFPs
	s.mu.Unlock()

	return songID, nil
}

type RecognizeResult struct {
	Song       *db.Song
	Confidence float64
	IsMatch    bool
	TimeInSong float64
}

func (s *Service) Recognize(snippetPath string) (*RecognizeResult, error) {
	queryPeaks, err := audio.ExtractQueryFingerprints(snippetPath, audio.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("fingerprint extraction failed: %w", err)
	}
	if len(queryPeaks) == 0 {
		return nil, fmt.Errorf("no fingerprints extracted — audio too short or too quiet")
	}

	songs, err := s.getSongs()
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 {
		return nil, fmt.Errorf("no songs enrolled — add tracks in Enrollment Flow first")
	}

	best := matching.Result{MatchOffset: -1}
	secondSongHits := 0
	var bestSong *db.Song

	for i := range songs {
		song := &songs[i]
		refPeaks, err := s.getSongPeakFingerprints(song)
		if err != nil {
			continue
		}

		result := matching.MatchPeaksWindowed(refPeaks, queryPeaks)
		if result.HitCount > best.HitCount || (result.HitCount == best.HitCount && result.Confidence > best.Confidence) {
			secondSongHits = best.HitCount
			best = result
			bestSong = song
		} else if result.HitCount > secondSongHits {
			secondSongHits = result.HitCount
		}
	}

	if bestSong == nil || !best.IsMatch || best.HitCount < secondSongHits+5 {
		if bestSong != nil && best.HitCount > 0 {
			return nil, fmt.Errorf(
				"match too weak (%d aligned hashes) — closest was \"%s\" by %s, try 6-8 seconds closer to the speaker",
				best.HitCount, bestSong.Title, bestSong.Artist,
			)
		}
		names := make([]string, len(songs))
		for i := range songs {
			names[i] = songs[i].Title + " — " + songs[i].Artist
		}
		return nil, fmt.Errorf("no match — play one of the enrolled songs: %v", names)
	}

	displayConfidence := 100 * float64(best.HitCount) / float64(best.HitCount+maxInt(secondSongHits, 1))

	return &RecognizeResult{
		Song:       bestSong,
		Confidence: displayConfidence,
		IsMatch:    true,
		TimeInSong: best.TimeInSong,
	}, nil
}

func (s *Service) getSongPeakFingerprints(song *db.Song) ([]dsp.Fingerprint, error) {
	s.mu.RLock()
	if fps, ok := s.peakFingerprints[song.ID]; ok {
		s.mu.RUnlock()
		return fps, nil
	}
	s.mu.RUnlock()

	storedFingerprints, err := db.GetPeakFingerprints(song.ID)
	if err == nil && len(storedFingerprints) > 0 {
		fps := toDSPFingerprints(storedFingerprints)
		s.mu.Lock()
		s.peakFingerprints[song.ID] = fps
		s.mu.Unlock()
		return fps, nil
	}

	fps, err := audio.ExtractFingerprints(song.FilePath, audio.DefaultConfig())
	if err != nil {
		return nil, err
	}
	if len(fps) == 0 {
		return nil, fmt.Errorf("no peak fingerprints generated for %s", song.Title)
	}

	s.mu.Lock()
	s.peakFingerprints[song.ID] = fps
	s.mu.Unlock()

	// Best-effort backfill for legacy songs that were enrolled before peak fingerprints
	// were persisted in Postgres.
	_ = db.InsertPeakFingerprints(song.ID, toDBPeakFingerprints(fps))

	return fps, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func toDBPeakFingerprints(fps []dsp.Fingerprint) []db.PeakFingerprint {
	out := make([]db.PeakFingerprint, 0, len(fps))
	for _, fp := range fps {
		out = append(out, db.PeakFingerprint{
			Hash:   fp.Hash,
			Offset: fp.Offset,
		})
	}
	return out
}

func toDSPFingerprints(fps []db.PeakFingerprint) []dsp.Fingerprint {
	out := make([]dsp.Fingerprint, 0, len(fps))
	for _, fp := range fps {
		out = append(out, dsp.Fingerprint{
			Hash:   fp.Hash,
			Offset: fp.Offset,
		})
	}
	return out
}
