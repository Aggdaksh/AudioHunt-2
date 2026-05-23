import { useRef } from 'react';
import AudioRecorder, { type AudioRecorderHandle } from '../components/AudioRecorder';
import HistoryList from '../components/HistoryList';
import RecognizeButton from '../components/RecognizeButton';
import SongCard from '../components/SongCard';
import WaveformVisualizer from '../components/WaveformVisualizer';
import { useRecognition } from '../hooks/useRecognition';
import { useWaveform } from '../hooks/useWaveform';
import type { RecognitionResult } from '../services/api';
import type { HistoryItem } from '../types';

interface RecognizePageProps {
  addToHistory: (result: RecognitionResult) => void;
  history: HistoryItem[];
}

export default function RecognizePage({ addToHistory, history }: RecognizePageProps) {
  const minRecordingMs = 8_000;
  const { frequencyData, connectStream } = useWaveform();
  const recognition = useRecognition({ minDurationMs: minRecordingMs, onMatch: addToHistory });
  const recorderRef = useRef<AudioRecorderHandle | null>(null);
  const recentHistory = history.slice(0, 3);

  const isRecording = recognition.status === 'recording';

  return (
    <section className="page">
      <div className="hero-panel">
        <div className="hero-panel__copy">
          <p className="eyebrow">Recognition flow</p>
          <h1>Catch the song in one clean hit.</h1>
          <p className="hero-panel__lead">
            Play the music loudly near the mic, keep the source steady, and stop only after a clean 8 to 10 second clip.
          </p>
        </div>
        <div className="hero-panel__guidance">
          <div className="hero-panel__guide-card">
            <strong>Recording tip</strong>
            <p>Play the song close to the mic for at least 8 seconds. Aim for a clean 8 to 10 second window.</p>
          </div>
          <div className="hero-panel__guide-card">
            <strong>Explore the app</strong>
            <p>Open Library to manage enrolled songs and use History to revisit your latest recognitions.</p>
          </div>
        </div>
      </div>

      <div className="recognize-layout">
        <section className="section-card section-card--stage">
          <div className="section-heading">
            <div>
              <p className="eyebrow">Live scan</p>
              <h2>Point the mic toward the music source</h2>
            </div>
          </div>

          <AudioRecorder
            ref={recorderRef}
            minDurationMs={minRecordingMs}
            onError={recognition.handleRecordingError}
            onRecordingComplete={recognition.handleRecordingComplete}
            onStreamReady={connectStream}
            sessionId={recognition.sessionId}
            stopRequestId={recognition.stopRequestId}
          />

          <div className="recognize-stage">
            <WaveformVisualizer frequencyData={frequencyData} visible={isRecording} />
            <RecognizeButton onPress={recognition.start} status={recognition.status} />
            {isRecording ? (
              <div className="recording-toolbar recording-toolbar--single">
                <button
                  className="stop-pill"
                  onClick={() => {
                    recognition.requestStop();
                    recorderRef.current?.stop();
                  }}
                  type="button"
                >
                  Stop
                </button>
              </div>
            ) : (
              recognition.status !== 'processing' ? (
                <p className="recognize-stage__caption">
                  Start recording, play the music loudly near the mic, and press Stop after at least 8 seconds.
                </p>
              ) : null
            )}
          </div>

          {recognition.status === 'match' && recognition.result ? (
            <SongCard
              latencyMs={recognition.latencyMs}
              onDismiss={recognition.reset}
              result={recognition.result}
            />
          ) : null}

          {recognition.status === 'no_match' ? (
            <div className="feedback-card feedback-card--warning">
              <strong>No confident match this time.</strong>
              <p>{recognition.error ?? 'Try another 8 to 10 second pass with the music source closer to the mic.'}</p>
              <button className="ghost-btn" onClick={recognition.reset} type="button">
                Reset scan
              </button>
            </div>
          ) : null}

          {recognition.status === 'error' ? (
            <div className="feedback-card feedback-card--error">
              <strong>Recording hit a problem.</strong>
              <p>{recognition.error ?? 'Microphone access or upload failed.'}</p>
              <button className="ghost-btn" onClick={recognition.reset} type="button">
                Try again
              </button>
            </div>
          ) : null}
        </section>

        <section className="section-card">
          <div className="section-heading">
            <div>
              <p className="eyebrow">Recent finds</p>
              <h2>Last three recognitions</h2>
            </div>
          </div>
          <HistoryList compact items={recentHistory} />
        </section>
      </div>
    </section>
  );
}
