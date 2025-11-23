/**
 * Input Component
 * Переиспользуемый UI компонент инпута
 */

import { forwardRef } from 'react';
import type { InputHTMLAttributes } from 'react';
import './Input.css';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  hint?: string;
  hintType?: 'error' | 'success' | 'info';
  fullWidth?: boolean;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, hint, hintType = 'info', fullWidth = false, className = '', ...props }, ref) => {
    const wrapperClasses = [
      'input-wrapper',
      fullWidth && 'input-wrapper--full-width',
    ]
      .filter(Boolean)
      .join(' ');

    const inputClasses = [
      'input',
      error && 'input--error',
      className,
    ]
      .filter(Boolean)
      .join(' ');

    const hintClasses = [
      'input-hint',
      `input-hint--${hintType}`,
    ]
      .filter(Boolean)
      .join(' ');

    return (
      <div className={wrapperClasses}>
        {(label || hint) && (
          <div className="input-header">
            {label && <label className="input-label">{label}</label>}
            {hint && <span className={hintClasses}>{hint}</span>}
          </div>
        )}
        <input ref={ref} className={inputClasses} {...props} />
        {error && <span className="input-error">{error}</span>}
      </div>
    );
  }
);

Input.displayName = 'Input';

