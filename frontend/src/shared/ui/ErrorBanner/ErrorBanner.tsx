/**
 * ErrorBanner Component
 * Красная плашка сверху для отображения ошибок API
 */

import { useErrorStore } from '../../lib/error-store';
import './ErrorBanner.css';

export const ErrorBanner = () => {
  const { error, clearError } = useErrorStore();

  if (!error) return null;

  return (
    <div className="error-banner">
      <div className="error-banner__content">
        <span className="error-banner__message">{error.message}</span>
        {error.code && (
          <span className="error-banner__code">Код: {error.code}</span>
        )}
      </div>
      <button
        className="error-banner__close"
        onClick={clearError}
        aria-label="Закрыть"
      >
        ×
      </button>
    </div>
  );
};

