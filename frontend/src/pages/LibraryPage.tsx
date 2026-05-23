import EnrollForm from '../components/EnrollForm';
import type { CatalogSong } from '../services/api';

interface LibraryPageProps {
  catalog: CatalogSong[];
  onEnrolled: () => void;
}

export default function LibraryPage({ catalog, onEnrolled }: LibraryPageProps) {
  return (
    <section className="page">
      <div className="page-banner">
        <div>
          <p className="eyebrow">Library</p>
          <h1>Shape the catalog your recognizer listens against.</h1>
          <p className="hero-panel__lead">
            Drag in songs, give them clean labels, and keep your current enrolled tracks visible in one place.
          </p>
        </div>
      </div>

      <div className="library-grid">
        <section className="section-card">
          <div className="section-heading">
            <div>
              <p className="eyebrow">Current catalog</p>
              <h2>{catalog.length} enrolled tracks</h2>
            </div>
          </div>

          {catalog.length === 0 ? (
            <div className="empty-state empty-state--inline">
              <strong>Your library is empty.</strong>
              <p>Add at least one song here before using the recognition flow.</p>
            </div>
          ) : (
            <div className="library-list">
              {catalog.map((song) => (
                <article className="library-row" key={song.id}>
                  <div>
                    <strong>{song.title}</strong>
                    <p>{song.artist}</p>
                  </div>
                  <span className="stat-pill">ID {song.id}</span>
                </article>
              ))}
            </div>
          )}
        </section>

        <section className="section-card">
          <EnrollForm onEnrolled={onEnrolled} />
        </section>
      </div>
    </section>
  );
}
