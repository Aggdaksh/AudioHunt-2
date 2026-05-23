import type { HistoryItem } from '../types';

interface HistoryListProps {
  compact?: boolean;
  items: HistoryItem[];
  onDelete?: (id: string) => void;
}

function formatRelativeTime(timestamp: string) {
  const deltaMs = Date.now() - new Date(timestamp).getTime();
  const minute = 60_000;
  const hour = 60 * minute;
  const day = 24 * hour;

  if (deltaMs < minute) {
    return 'just now';
  }
  if (deltaMs < hour) {
    return `${Math.round(deltaMs / minute)} min ago`;
  }
  if (deltaMs < day) {
    return `${Math.round(deltaMs / hour)} hr ago`;
  }

  return `${Math.round(deltaMs / day)} day ago`;
}

export default function HistoryList({ compact = false, items, onDelete }: HistoryListProps) {
  if (items.length === 0) {
    return (
      <div className="empty-state empty-state--inline">
        <p>Nothing recognized yet. Your matches will start stacking here after the first successful scan.</p>
      </div>
    );
  }

  return (
    <div className={`history-list ${compact ? 'history-list--compact' : ''}`}>
      {items.map((item) => (
        <article className="history-row" key={item.id}>
          <div className="history-row__main">
            <div className="history-row__title-line">
              <strong>{item.title}</strong>
              <span className="history-row__confidence">{Math.round(item.confidence)}%</span>
            </div>
            <p>{item.artist}</p>
          </div>
          <div className="history-row__meta">
            <span>{formatRelativeTime(item.capturedAt)}</span>
            {typeof item.timeInSong === 'number' ? <span>{Math.floor(item.timeInSong)}s mark</span> : null}
          </div>
          {!compact && onDelete ? (
            <button className="history-row__delete" onClick={() => onDelete(item.id)} type="button">
              Remove
            </button>
          ) : null}
        </article>
      ))}
    </div>
  );
}
