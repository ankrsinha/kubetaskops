import { Sun, Moon } from 'lucide-react';

export function ThemeToggle({ theme, onToggle }) {
  return (
    <button
      id="btn-theme-toggle"
      className="btn-theme-toggle"
      onClick={onToggle}
      title={theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
      aria-label="Toggle theme"
    >
      {theme === 'dark' ? (
        <Sun size={18} className="theme-icon sun" />
      ) : (
        <Moon size={18} className="theme-icon moon" />
      )}
    </button>
  );
}
