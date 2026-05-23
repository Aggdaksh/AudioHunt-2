import { useCallback, useEffect, useRef, useState } from 'react';

type WebkitWindow = Window & {
  webkitAudioContext?: typeof AudioContext;
};

export function useWaveform() {
  const [frequencyData, setFrequencyData] = useState<Uint8Array>(new Uint8Array(128));

  const audioContextRef = useRef<AudioContext | null>(null);
  const analyserRef = useRef<AnalyserNode | null>(null);
  const sourceRef = useRef<MediaStreamAudioSourceNode | null>(null);
  const animationFrameRef = useRef<number | null>(null);

  const stop = useCallback(async () => {
    if (animationFrameRef.current !== null) {
      window.cancelAnimationFrame(animationFrameRef.current);
      animationFrameRef.current = null;
    }

    if (sourceRef.current) {
      sourceRef.current.disconnect();
      sourceRef.current = null;
    }

    if (analyserRef.current) {
      analyserRef.current.disconnect();
      analyserRef.current = null;
    }

    if (audioContextRef.current) {
      try {
        await audioContextRef.current.close();
      } catch {
        // Ignore close errors during rapid session changes.
      }
      audioContextRef.current = null;
    }

    setFrequencyData(new Uint8Array(128));
  }, []);

  const connectStream = useCallback(
    async (stream: MediaStream | null) => {
      if (!stream) {
        await stop();
        return;
      }

      await stop();

      const AudioContextConstructor =
        window.AudioContext || (window as WebkitWindow).webkitAudioContext;

      if (!AudioContextConstructor) {
        return;
      }

      const audioContext = new AudioContextConstructor();
      const analyser = audioContext.createAnalyser();
      analyser.fftSize = 256;
      analyser.smoothingTimeConstant = 0.82;

      const source = audioContext.createMediaStreamSource(stream);
      source.connect(analyser);

      audioContextRef.current = audioContext;
      analyserRef.current = analyser;
      sourceRef.current = source;

      if (audioContext.state === 'suspended') {
        await audioContext.resume();
      }

      const tick = () => {
        if (!analyserRef.current) {
          return;
        }

        const next = new Uint8Array(analyserRef.current.frequencyBinCount);
        analyserRef.current.getByteFrequencyData(next);
        setFrequencyData(next);
        animationFrameRef.current = window.requestAnimationFrame(tick);
      };

      tick();
    },
    [stop]
  );

  useEffect(() => {
    return () => {
      void stop();
    };
  }, [stop]);

  return {
    frequencyData,
    connectStream,
    stop,
  };
}
