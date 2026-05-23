import HistoryList from '../components/HistoryList';
import type { HistoryItem } from '../types';

interface HistoryPageProps {
  history: HistoryItem[];
  onClear: () => void;
  onDelete: (id: string) => void;
  onGoRecognize: () => void;
}

function averageConfidence(history: HistoryItem[]) {
  if (history.length === 0) {
    return 0;
  }

  const total = history.reduce((sum, item) => sum + item.confidence, 0);
  return Math.round(total / history.length);
}

export default function HistoryPage({ history, onClear, onDelete, onGoRecognize }: HistoryPageProps) {
  return (
    <section className="page">
      <div className="page-banner">
        <div>
          <p className="eyebrow">History</p>
          <h1>Every recognition, in one running log.</h1>
          <p className="hero-panel__lead">Keep a lightweight journal of what the app has identified on this device.</p>
        </div>
        <div className="summary-grid">
          <div className="summary-card">
            <span>Total matches</span>
            <strong>{history.length}</strong>
          </div>
          <div className="summary-card">
            <span>Avg confidence</span>
            <strong>{averageConfidence(history)}%</strong>
          </div>
        </div>
      </div>

      <section className="section-card">
        <div className="section-heading">
          <div>
            <p className="eyebrow">Recognition archive</p>
            <h2>Newest matches first</h2>
          </div>
          <div className="inline-actions">
            <button className="ghost-btn" onClick={onGoRecognize} type="button">
              New scan
            </button>
            <button className="ghost-btn" disabled={history.length === 0} onClick={onClear} type="button">
              Clear all
            </button>
          </div>
        </div>

        {history.length === 0 ? (
          <div className="empty-state">
            <strong>Nothing recognized yet.</strong>
            <p>Run your first scan from the Recognize page and the results will land here automatically.</p>
            <button className="primary-btn" onClick={onGoRecognize} type="button">
              Go to Recognize
            </button>
          </div>
        ) : (
          <HistoryList items={history} onDelete={onDelete} />
        )}
      </section>
    </section>
  );
}
