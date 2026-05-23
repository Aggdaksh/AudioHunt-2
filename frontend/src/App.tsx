import { useEffect, useState } from 'react';
import Navbar from './components/Navbar';
import { useHistory } from './hooks/useHistory';
import HistoryPage from './pages/HistoryPage';
import LibraryPage from './pages/LibraryPage';
import RecognizePage from './pages/RecognizePage';
import { checkHealth, fetchCatalog, type CatalogSong } from './services/api';
import type { AppPage } from './types';
import './App.css';

function backendStatusLabel(status: 'checking' | 'online' | 'offline') {
  if (status === 'online') {
    return 'Backend online';
  }

  if (status === 'offline') {
    return 'Backend offline';
  }

  return 'Checking backend';
}

function App() {
  const [page, setPage] = useState<AppPage>('recognize');
  const [catalog, setCatalog] = useState<CatalogSong[]>([]);
  const [backendStatus, setBackendStatus] = useState<'checking' | 'online' | 'offline'>('checking');
  const { history, addToHistory, removeFromHistory, clearHistory } = useHistory();

  const refreshCatalog = async () => {
    try {
      const songs = await fetchCatalog();
      setCatalog(songs);
    } catch {
      setCatalog([]);
    }
  };

  const refreshHealth = async () => {
    try {
      await checkHealth();
      setBackendStatus('online');
    } catch {
      setBackendStatus('offline');
    }
  };

  useEffect(() => {
    void refreshCatalog();
    void refreshHealth();

    const intervalId = window.setInterval(() => {
      void refreshHealth();
    }, 20_000);

    return () => {
      window.clearInterval(intervalId);
    };
  }, []);

  const handlePageChange = (nextPage: AppPage) => {
    setPage(nextPage);
  };

  const handleBackendCheck = () => {
    setBackendStatus('checking');
    void refreshHealth();
  };

  return (
    <div className="app-shell">
      <div className="app-shell__glow app-shell__glow--left" aria-hidden="true" />
      <div className="app-shell__glow app-shell__glow--right" aria-hidden="true" />

      <div className="app-frame">
        <header className="app-header">
          <div className="app-brand">
            <div className="app-brand__logo-shell" aria-hidden="true">
              <img alt="" className="app-brand__logo" src="/audiohunt-mark.svg" />
            </div>
            <div className="app-brand__copy">
              <p className="app-header__kicker">Realtime audio fingerprinting</p>
              <h1 className="app-brand__wordmark">AudioHunt</h1>
            </div>
          </div>
          <button
            className={`app-header__backend-button app-header__backend-button--${backendStatus}`}
            onClick={handleBackendCheck}
            type="button"
          >
            <span className="app-header__backend-dot" aria-hidden="true" />
            <span>{backendStatusLabel(backendStatus)}</span>
          </button>
        </header>

        <main className="app-main">
          {page === 'recognize' ? (
            <RecognizePage addToHistory={addToHistory} history={history} />
          ) : null}

          {page === 'history' ? (
            <HistoryPage
              history={history}
              onClear={clearHistory}
              onDelete={removeFromHistory}
              onGoRecognize={() => handlePageChange('recognize')}
            />
          ) : null}

          {page === 'library' ? (
            <LibraryPage
              catalog={catalog}
              onEnrolled={() => {
                void refreshCatalog();
              }}
            />
          ) : null}
        </main>

        <footer className="site-footer">
          <section className="site-footer__bar">
            <div className="site-footer__socials" aria-label="Contact links">
              <a
                aria-label="Email Daksh Aggarwal"
                className="site-footer__social-link"
                href="mailto:daksh121105@gmail.com"
              >
                <svg fill="none" viewBox="0 0 24 24">
                  <path
                    d="M4 7.5 11.07 12.45a1.6 1.6 0 0 0 1.86 0L20 7.5M5.6 19h12.8A1.6 1.6 0 0 0 20 17.4V8.6A1.6 1.6 0 0 0 18.4 7H5.6A1.6 1.6 0 0 0 4 8.6v8.8A1.6 1.6 0 0 0 5.6 19Z"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="1.8"
                  />
                </svg>
              </a>

              <a
                aria-label="Daksh Aggarwal on LinkedIn"
                className="site-footer__social-link"
                href="https://www.linkedin.com/in/daksh-aggarwal-a5938028b/"
                rel="noreferrer"
                target="_blank"
              >
                <svg fill="currentColor" viewBox="0 0 24 24">
                  <path d="M6.94 8.5A1.44 1.44 0 1 1 6.95 5.62a1.44 1.44 0 0 1-.01 2.88ZM5.74 9.8h2.43V18H5.74V9.8Zm3.95 0h2.34v1.12h.03c.33-.62 1.13-1.28 2.33-1.28 2.5 0 2.96 1.64 2.96 3.77V18h-2.43v-4.07c0-.97-.02-2.22-1.35-2.22-1.35 0-1.55 1.05-1.55 2.15V18H9.69V9.8Z" />
                </svg>
              </a>
            </div>

            <div className="site-footer__credit">
              <span>Made by Daksh Aggarwal</span>
              <strong>© 2026 AudioHunt</strong>
            </div>
          </section>
        </footer>

        <Navbar activePage={page} onChange={handlePageChange} />
      </div>
    </div>
  );
}

export default App;
