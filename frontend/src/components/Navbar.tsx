import type { AppPage } from '../types';

interface NavbarProps {
  activePage: AppPage;
  onChange: (page: AppPage) => void;
}

const items: Array<{ id: AppPage; label: string; short: string }> = [
  { id: 'recognize', label: 'Recognize', short: 'REC' },
  { id: 'history', label: 'History', short: 'LOG' },
  { id: 'library', label: 'Library', short: 'LIB' },
];

export default function Navbar({ activePage, onChange }: NavbarProps) {
  return (
    <nav aria-label="Primary" className="bottom-nav">
      {items.map((item) => (
        <button
          className={`bottom-nav__item ${activePage === item.id ? 'bottom-nav__item--active' : ''}`}
          key={item.id}
          onClick={() => onChange(item.id)}
          type="button"
        >
          <span className="bottom-nav__mono">{item.short}</span>
          <span>{item.label}</span>
        </button>
      ))}
    </nav>
  );
}
