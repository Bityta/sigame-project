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
  canAnswer: boolean;
  onPressButton?: () => void;
  onSubmitAnswer?: (answer: string) => void;
  timeRemaining?: number;
}

export const QuestionView = ({
  question,
  canPressButton,
  canAnswer,
  onPressButton,
  onSubmitAnswer,
  timeRemaining,
}: QuestionViewProps) => {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const answer = formData.get('answer') as string;
    if (answer.trim()) {
      onSubmitAnswer?.(answer);
    }
  };

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

        {canAnswer && (
          <form onSubmit={handleSubmit} className="question-view__answer-form">
            <input
              type="text"
              name="answer"
              placeholder="Ваш ответ..."
              className="question-view__answer-input"
              autoFocus
              autoComplete="off"
            />
            <Button type="submit" variant="primary" size="large">
              Отправить ответ
            </Button>
          </form>
        )}
      </div>
    </Card>
  );
};

