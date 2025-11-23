/**
 * Card Component
 * Переиспользуемый UI компонент карточки
 */

import { forwardRef } from 'react';
import type { HTMLAttributes } from 'react';
import './Card.css';

export interface CardProps extends HTMLAttributes<HTMLDivElement> {
  padding?: 'none' | 'small' | 'medium' | 'large';
}

export const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ children, padding = 'medium', className = '', ...props }, ref) => {
    const classes = [
      'card',
      `card--padding-${padding}`,
      className,
    ]
      .filter(Boolean)
      .join(' ');

    return (
      <div ref={ref} className={classes} {...props}>
        {children}
      </div>
    );
  }
);

Card.displayName = 'Card';

