/**
 * Toggle Switch Component
 * Красивый переключатель вкл/выкл
 */

import { InputHTMLAttributes } from 'react';
import './Toggle.css';

interface ToggleProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  label?: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
}

export const Toggle = ({ label, checked, onChange, disabled, id, ...props }: ToggleProps) => {
  const toggleId = id || `toggle-${Math.random().toString(36).substr(2, 9)}`;

  return (
    <div className={`toggle ${disabled ? 'toggle--disabled' : ''}`}>
      <input
        type="checkbox"
        id={toggleId}
        className="toggle__input"
        checked={checked}
        onChange={(e) => onChange(e.target.checked)}
        disabled={disabled}
        {...props}
      />
      <label htmlFor={toggleId} className="toggle__label">
        <span className="toggle__switch" />
        {label && <span className="toggle__text">{label}</span>}
      </label>
    </div>
  );
};

