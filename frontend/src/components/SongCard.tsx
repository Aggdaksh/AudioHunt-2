import type { RecognitionResult } from '../services/api';

interface SongCardProps {
  latencyMs: number | null;
  onDismiss: () => void;
  result: RecognitionResult;
}

function toneForConfidence(confidence: number) {
  if (confidence >= 70) {
    return 'mint';
  }
  if (confidence >= 40) {
    return 'amber';
  }
  return 'coral';
}

export default function SongCard({ latencyMs, onDismiss, result }: SongCardProps) {
  const confidence = Math.max(0, Math.min(100, Math.round(result.confidence)));
  const confidenceTone = toneForConfidence(confidence);

  return (
    <article className="song-card">
      <div className="song-card__header">
        <div>
          <p className="eyebrow">Match found</p>
          <h2>{result.title}</h2>
          <p className="song-card__artist">{result.artist}</p>
        </div>
        <button className="ghost-btn" onClick={onDismiss} type="button">
          Listen again
        </button>
      </div>

      <div className="song-card__meta">
        <span className={`stat-pill stat-pill--${confidenceTone}`}>{confidence}% confidence</span>
        {latencyMs !== null ? <span className="stat-pill">{latencyMs} ms round trip</span> : null}
        {typeof result.time_in_song === 'number' ? (
          <span className="stat-pill">Around {Math.floor(result.time_in_song)}s into the track</span>
        ) : null}
      </div>

      <div className="confidence-meter" aria-label={`Confidence ${confidence} percent`}>
        <div className="confidence-meter__track">
          <div
            className={`confidence-meter__fill confidence-meter__fill--${confidenceTone}`}
            style={{ width: `${confidence}%` }}
          />
        </div>
      </div>
    </article>
  );
}
