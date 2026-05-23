export type AppPage = 'recognize' | 'history' | 'library';

export type RecognitionStatus =
  | 'idle'
  | 'recording'
  | 'processing'
  | 'match'
  | 'no_match'
  | 'error';

export interface HistoryItem {
  id: string;
  songId: number;
  title: string;
  artist: string;
  confidence: number;
  timeInSong?: number;
  capturedAt: string;
}

export interface RecordingPayload {
  blob: Blob;
  durationMs: number;
  filename: string;
}
