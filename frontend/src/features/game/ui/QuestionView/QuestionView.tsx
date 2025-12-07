/**
 * Game Feature - QuestionView
 * Отображение текущего вопроса
 */

import { Card, Button } from '@/shared/ui';
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
  return (
    <Card className="question-view" padding="large">
      <div className="question-view__header">
        <div className="question-view__price">{question.price} очков</div>
        {timeRemaining !== undefined && (
          <div className="question-view__timer">{timeRemaining}с</div>
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
          <Button
            variant="danger"
            size="large"
            fullWidth
            onClick={onPressButton}
          >
            🔴 Нажать кнопку!
          </Button>
        )}
      </div>
    </Card>
  );
};
