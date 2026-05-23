import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface RecognitionResult {
  song_id: number;
  title: string;
  artist: string;
  confidence: number; // 0–100 percent
  is_match: boolean;
  time_in_song?: number;
}

export const recognizeSong = async (
  audioBlob: Blob,
  filename = 'recording.webm'
): Promise<RecognitionResult> => {
  const formData = new FormData();
  formData.append('audio_snippet', audioBlob, filename);

  const response = await axios.post(
    `${API_URL}/api/recognize`,
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      timeout: 15000,
    }
  );

  return response.data;
};

export interface EnrollmentResult {
  id: number;
  message: string;
}

export interface CatalogSong {
  id: number;
  title: string;
  artist: string;
}

export interface BackendHealth {
  status: string;
}

export const fetchCatalog = async (): Promise<CatalogSong[]> => {
  const response = await axios.get<CatalogSong[]>(`${API_URL}/api/songs`);
  return response.data;
};

export const checkHealth = async (): Promise<BackendHealth> => {
  const response = await axios.get<BackendHealth>(`${API_URL}/api/health`);
  return response.data;
};

export const enrollSong = async (
  audioFile: File,
  title: string,
  artist: string
): Promise<EnrollmentResult> => {
  const formData = new FormData();
  formData.append('audio_file', audioFile);
  formData.append('title', title);
  formData.append('artist', artist);

  const response = await axios.post(
    `${API_URL}/api/songs`,
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      timeout: 150000,
      timeoutErrorMessage:
        'Enrollment took too long on the free cloud server. Please try a shorter file or retry after the backend wakes up.',
    }
  );

  return response.data;
};
