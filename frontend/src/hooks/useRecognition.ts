import { useEffect, useRef, useState } from 'react';
import { recognizeSong, type RecognitionResult } from '../services/api';
import type { RecordingPayload, RecognitionStatus } from '../types';

interface UseRecognitionOptions {
  minDurationMs: number;
  onMatch: (result: RecognitionResult) => void;
}

export function useRecognition({ minDurationMs, onMatch }: UseRecognitionOptions) {
  const [status, setStatus] = useState<RecognitionStatus>('idle');
  const [result, setResult] = useState<RecognitionResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [sessionId, setSessionId] = useState(0);
  const [stopRequestId, setStopRequestId] = useState(0);
  const [latencyMs, setLatencyMs] = useState<number | null>(null);
  const stopTimeoutRef = useRef<number | null>(null);
  const recordingStartedAtRef = useRef<number | null>(null);
  const shortStopLockedRef = useRef(false);

  const clearStopTimeout = () => {
    if (stopTimeoutRef.current === null) {
      return;
    }

    window.clearTimeout(stopTimeoutRef.current);
    stopTimeoutRef.current = null;
  };

  useEffect(() => {
    return () => {
      clearStopTimeout();
    };
  }, []);

  const start = () => {
    if (status === 'recording' || status === 'processing') {
      return;
    }

    clearStopTimeout();
    recordingStartedAtRef.current = performance.now();
    shortStopLockedRef.current = false;
    setResult(null);
    setError(null);
    setLatencyMs(null);
    setStatus('recording');
    setSessionId((current) => current + 1);
  };

  const reset = () => {
    clearStopTimeout();
    recordingStartedAtRef.current = null;
    shortStopLockedRef.current = false;
    setResult(null);
    setError(null);
    setLatencyMs(null);
    setStatus('idle');
  };

  const requestStop = () => {
    if (status !== 'recording') {
      return;
    }

    clearStopTimeout();
    setError(null);
    setStopRequestId((current) => current + 1);

    const elapsedMs =
      recordingStartedAtRef.current === null ? 0 : performance.now() - recordingStartedAtRef.current;

    if (elapsedMs < minDurationMs) {
      recordingStartedAtRef.current = null;
      shortStopLockedRef.current = true;
      setResult(null);
      setLatencyMs(null);
      setError(`Recording time not enough. Record for at least ${Math.ceil(minDurationMs / 1000)} seconds.`);
      setStatus('error');
      return;
    }

    stopTimeoutRef.current = window.setTimeout(() => {
      stopTimeoutRef.current = null;
      recordingStartedAtRef.current = null;
      setResult(null);
      setLatencyMs(null);
      setError('Recording could not be finalized. Please try again.');
      setStatus((current) => (current === 'recording' ? 'error' : current));
    }, 2500);
  };

  const handleRecordingComplete = async ({ blob, filename }: RecordingPayload) => {
    if (shortStopLockedRef.current) {
      return;
    }

    const startedAt = performance.now();

    clearStopTimeout();
    recordingStartedAtRef.current = null;
    setStatus('processing');
    setError(null);

    try {
      const response = await recognizeSong(blob, filename);
      const duration = Math.round(performance.now() - startedAt);

      setLatencyMs(duration);
      setResult(response);

      if (response.is_match) {
        setStatus('match');
        onMatch(response);
        return;
      }

      setStatus('no_match');
      setError('No confident match came back from the recognizer.');
    } catch (err: unknown) {
      const responseStatus = (err as { response?: { status?: number } })?.response?.status;
      const responseData = (err as { response?: { data?: unknown } })?.response?.data;
      const errorCode = (err as { code?: string })?.code;
      const serverMessage =
        typeof responseData === 'string'
          ? responseData
          : errorCode === 'ECONNABORTED'
            ? 'The backend took too long to respond. Please try again.'
            : err instanceof Error
              ? err.message
              : 'Recognition failed';

      setResult(null);
      setLatencyMs(Math.round(performance.now() - startedAt));
      setError(serverMessage);

      if (responseStatus === 404) {
        setStatus('no_match');
        return;
      }

      setStatus('error');
    }
  };

  const handleRecordingError = (message: string) => {
    if (shortStopLockedRef.current) {
      return;
    }

    clearStopTimeout();
    recordingStartedAtRef.current = null;
    setResult(null);
    setLatencyMs(null);
    setError(message);
    setStatus('error');
  };

  return {
    status,
    result,
    error,
    sessionId,
    stopRequestId,
    latencyMs,
    start,
    requestStop,
    reset,
    handleRecordingComplete,
    handleRecordingError,
  };
}
