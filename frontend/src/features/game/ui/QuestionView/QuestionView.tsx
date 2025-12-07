/**
 * Game Feature - QuestionView
 * Отображение текущего вопроса с большой кнопкой
 */

import { Card } from '@/shared/ui';
import type { QuestionState } from '@/shared/types';
import './QuestionView.css';

interface QuestionViewProps {
  question: QuestionState;
  canPressButton: boolean;
  onPressButton?: () => void;
  timeRemaining?: number;
}

export const QuestionView = ({
  question,
  canPressButton,
  onPressButton,
  timeRemaining,
}: QuestionViewProps) => {
  const isTimerWarning = timeRemaining !== undefined && timeRemaining <= 5;
  const isTimerDanger = timeRemaining !== undefined && timeRemaining <= 3;

  return (
    <Card className="question-view" padding="large">
      <div className="question-view__header">
        <div className="question-view__price">{question.price} очков</div>
        {timeRemaining !== undefined && (
          <div className={`question-view__timer ${isTimerDanger ? 'question-view__timer--danger' : isTimerWarning ? 'question-view__timer--warning' : ''}`}>
            {timeRemaining}с
          </div>
        )}
      </div>

      {question.text && (
        <div className="question-view__text">{question.text}</div>
      )}

      {question.mediaType && question.mediaType !== 'text' && (
        <div className="question-view__media">
          <p>Медиа: {question.mediaType}</p>
        </div>
      )}

      <div className="question-view__actions">
        {canPressButton && (
          <button
            className="question-view__buzz-button"
            onClick={onPressButton}
          >
            🔴
            <span>ОТВЕТИТЬ!</span>
          </button>
        )}
      </div>
    </Card>
  );
};
