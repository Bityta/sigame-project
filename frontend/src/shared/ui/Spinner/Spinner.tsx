/**
 * Spinner Component
 * Индикатор загрузки
 */

import './Spinner.css';

export interface SpinnerProps {
  size?: 'small' | 'medium' | 'large';
  center?: boolean;
}

export const Spinner = ({ size = 'medium', center = false }: SpinnerProps) => {
  const classes = [
    'spinner',
    `spinner--${size}`,
    center && 'spinner--center',
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classes}>
      <div className="spinner__circle"></div>
    </div>
  );
};

