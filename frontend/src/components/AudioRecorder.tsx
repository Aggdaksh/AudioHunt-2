import { forwardRef, useEffect, useImperativeHandle, useRef } from 'react';
import type { RecordingPayload } from '../types';

interface AudioRecorderProps {
  minDurationMs?: number;
  onError: (error: string) => void;
  onRecordingComplete: (payload: RecordingPayload) => void;
  onStreamReady?: (stream: MediaStream | null) => void;
  sessionId: number;
  stopRequestId?: number;
}

export interface AudioRecorderHandle {
  stop: () => void;
}

function pickRecorderMimeType(): string {
  const candidates = [
    'audio/webm;codecs=opus',
    'audio/webm',
    'audio/mp4',
    'audio/ogg;codecs=opus',
    'audio/ogg',
  ];

  if (typeof MediaRecorder === 'undefined') {
    return '';
  }

  for (const mimeType of candidates) {
    if (MediaRecorder.isTypeSupported(mimeType)) {
      return mimeType;
    }
  }

  return '';
}

function fileNameForMimeType(mimeType: string): string {
  if (mimeType.includes('mp4')) {
    return 'recording.mp4';
  }
  if (mimeType.includes('ogg')) {
    return 'recording.ogg';
  }
  return 'recording.webm';
}

const AudioRecorder = forwardRef<AudioRecorderHandle, AudioRecorderProps>(function AudioRecorder(
  {
    minDurationMs = 6500,
    onError,
    onRecordingComplete,
    onStreamReady,
    sessionId,
    stopRequestId = 0,
  }: AudioRecorderProps,
  ref
) {
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const chunksRef = useRef<Blob[]>([]);
  const mimeTypeRef = useRef<string>('');
  const recordingStartMsRef = useRef<number>(0);
  const lastHandledSessionIdRef = useRef(0);
  const lastHandledStopRequestIdRef = useRef(0);
  const activeSessionIdRef = useRef(0);
  const pendingManualStopRef = useRef(false);
  const manualStopCommittedRef = useRef(false);
  const shouldUploadRef = useRef(false);
  const isUnmountedRef = useRef(false);
  const finalizeTimeoutRef = useRef<number | null>(null);
  const outcomeDeliveredRef = useRef(false);

  const onErrorRef = useRef(onError);
  const onRecordingCompleteRef = useRef(onRecordingComplete);
  const onStreamReadyRef = useRef(onStreamReady);

  useEffect(() => {
    onErrorRef.current = onError;
  }, [onError]);

  useEffect(() => {
    onRecordingCompleteRef.current = onRecordingComplete;
  }, [onRecordingComplete]);

  useEffect(() => {
    onStreamReadyRef.current = onStreamReady;
  }, [onStreamReady]);

  const emitError = (message: string) => {
    if (isUnmountedRef.current) {
      return;
    }

    onErrorRef.current(message);
  };

  const emitRecordingComplete = (payload: RecordingPayload) => {
    if (isUnmountedRef.current) {
      return;
    }

    onRecordingCompleteRef.current(payload);
  };

  const emitStreamReady = (stream: MediaStream | null) => {
    void onStreamReadyRef.current?.(stream);
  };

  const clearFinalizeTimeout = () => {
    if (finalizeTimeoutRef.current === null) {
      return;
    }

    window.clearTimeout(finalizeTimeoutRef.current);
    finalizeTimeoutRef.current = null;
  };

  const stopStream = () => {
    if (!streamRef.current) {
      return;
    }

    streamRef.current.getTracks().forEach((track) => track.stop());
    streamRef.current = null;
  };

  const clearRecorderRefs = () => {
    clearFinalizeTimeout();
    mediaRecorderRef.current = null;
    chunksRef.current = [];
    mimeTypeRef.current = '';
    recordingStartMsRef.current = 0;
    activeSessionIdRef.current = 0;
    pendingManualStopRef.current = false;
    manualStopCommittedRef.current = false;
    shouldUploadRef.current = false;
  };

  const releaseRecorderResources = () => {
    stopStream();
    emitStreamReady(null);
  };

  const failRecording = (message: string) => {
    clearRecorderRefs();
    releaseRecorderResources();
    emitError(message);
  };

  const finalizeRecording = (overrideShouldUpload?: boolean) => {
    if (outcomeDeliveredRef.current) {
      return;
    }

    outcomeDeliveredRef.current = true;

    const recorder = mediaRecorderRef.current;
    const blobType = mimeTypeRef.current || recorder?.mimeType || 'audio/webm';
    const elapsedMs = recordingStartMsRef.current ? Date.now() - recordingStartMsRef.current : 0;
    const shouldUpload = overrideShouldUpload ?? shouldUploadRef.current;
    const payload: RecordingPayload = {
      blob: new Blob(chunksRef.current, { type: blobType }),
      durationMs: elapsedMs,
      filename: fileNameForMimeType(blobType),
    };

    clearRecorderRefs();
    releaseRecorderResources();

    if (!shouldUpload) {
      return;
    }

    if (elapsedMs < minDurationMs) {
      emitError(`Recording time not enough. Record for at least ${(minDurationMs / 1000).toFixed(1)} seconds.`);
      return;
    }

    emitRecordingComplete(payload);
  };

  const requestStop = (reason: 'manual' | 'cleanup') => {
    const recorder = mediaRecorderRef.current;

    if (reason === 'manual') {
      pendingManualStopRef.current = true;
      manualStopCommittedRef.current = true;
    }

    if (!recorder || recorder.state === 'inactive') {
      if (reason === 'cleanup') {
        shouldUploadRef.current = false;
      }
      return;
    }

    if (reason === 'cleanup' && manualStopCommittedRef.current) {
      return;
    }

    shouldUploadRef.current = reason !== 'cleanup';

    try {
      recorder.requestData();
    } catch {
      // Some browsers reject requestData while the recorder is spinning up.
    }

    try {
      recorder.stop();
      if (reason === 'manual') {
        clearFinalizeTimeout();
        finalizeTimeoutRef.current = window.setTimeout(() => {
          try {
            recorder.requestData();
          } catch {
            // Ignore extra flush failures during fallback finalization.
          }

          finalizeRecording(true);
        }, 600);
      }
    } catch {
      failRecording('Recording could not be stopped cleanly. Please try again.');
    }
  };

  useImperativeHandle(
    ref,
    () => ({
      stop: () => {
        requestStop('manual');
      },
    }),
    []
  );

  useEffect(() => {
    return () => {
      isUnmountedRef.current = true;
      clearFinalizeTimeout();
      shouldUploadRef.current = false;
      outcomeDeliveredRef.current = true;

      const recorder = mediaRecorderRef.current;
      if (recorder && recorder.state === 'recording') {
        try {
          recorder.stop();
        } catch {
          // Ignore unmount stop errors and just release the stream below.
        }
      }

      releaseRecorderResources();
      clearRecorderRefs();
    };
  }, []);

  useEffect(() => {
    if (!sessionId || sessionId === lastHandledSessionIdRef.current) {
      return;
    }

    lastHandledSessionIdRef.current = sessionId;
    activeSessionIdRef.current = sessionId;
    pendingManualStopRef.current = false;
    manualStopCommittedRef.current = false;
    shouldUploadRef.current = false;
    outcomeDeliveredRef.current = false;
    chunksRef.current = [];
    mimeTypeRef.current = '';
    recordingStartMsRef.current = 0;
    clearFinalizeTimeout();

    const startRecording = async () => {
      if (typeof MediaRecorder === 'undefined') {
        emitError('This browser does not support audio recording.');
        return;
      }

      try {
        const stream = await navigator.mediaDevices.getUserMedia({
          audio: {
            echoCancellation: false,
            noiseSuppression: false,
            autoGainControl: false,
          },
        });

        if (isUnmountedRef.current || activeSessionIdRef.current !== sessionId) {
          stream.getTracks().forEach((track) => track.stop());
          return;
        }

        const mimeType = pickRecorderMimeType();
        const recorder = mimeType ? new MediaRecorder(stream, { mimeType }) : new MediaRecorder(stream);

        streamRef.current = stream;
        mediaRecorderRef.current = recorder;
        chunksRef.current = [];
        mimeTypeRef.current = mimeType || recorder.mimeType || 'audio/webm';
        recordingStartMsRef.current = Date.now();
        shouldUploadRef.current = true;

        emitStreamReady(stream);

        recorder.ondataavailable = (event) => {
          if (event.data.size > 0) {
            chunksRef.current.push(event.data);
          }
        };

        recorder.onerror = () => {
          failRecording('Recording failed before it could be processed. Please try again.');
        };

        recorder.onstop = () => {
          finalizeRecording();
        };

        recorder.start(250);

        if (pendingManualStopRef.current || manualStopCommittedRef.current) {
          requestStop('manual');
        }
      } catch {
        releaseRecorderResources();
        clearRecorderRefs();
        emitError('Microphone access was denied or unavailable.');
      }
    };

    void startRecording();
  }, [minDurationMs, sessionId]);

  useEffect(() => {
    if (!stopRequestId || stopRequestId === lastHandledStopRequestIdRef.current) {
      return;
    }

    lastHandledStopRequestIdRef.current = stopRequestId;
    requestStop('manual');
  }, [stopRequestId]);

  return null;
});

export default AudioRecorder;
