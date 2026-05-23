import { useEffect, useState } from 'react';
import type { RecognitionResult } from '../services/api';
import type { HistoryItem } from '../types';

const STORAGE_KEY = 'shazam-history';
const MAX_HISTORY_ITEMS = 50;

function readHistory(): HistoryItem[] {
  if (typeof window === 'undefined') {
    return [];
  }

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return [];
    }

    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) {
      return [];
    }

    return parsed.filter(
      (item): item is HistoryItem =>
        item &&
        typeof item.id === 'string' &&
        typeof item.songId === 'number' &&
        typeof item.title === 'string' &&
        typeof item.artist === 'string' &&
        typeof item.confidence === 'number' &&
        typeof item.capturedAt === 'string'
    );
  } catch {
    return [];
  }
}

export function useHistory() {
  const [history, setHistory] = useState<HistoryItem[]>(readHistory);

  useEffect(() => {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(history));
  }, [history]);

  const addToHistory = (result: RecognitionResult) => {
    setHistory((current) => {
      const nextItem: HistoryItem = {
        id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
        songId: result.song_id,
        title: result.title,
        artist: result.artist,
        confidence: result.confidence,
        timeInSong: result.time_in_song,
        capturedAt: new Date().toISOString(),
      };

      return [nextItem, ...current].slice(0, MAX_HISTORY_ITEMS);
    });
  };

  const removeFromHistory = (id: string) => {
    setHistory((current) => current.filter((item) => item.id !== id));
  };

  const clearHistory = () => {
    setHistory([]);
  };

  return {
    history,
    addToHistory,
    removeFromHistory,
    clearHistory,
  };
}
