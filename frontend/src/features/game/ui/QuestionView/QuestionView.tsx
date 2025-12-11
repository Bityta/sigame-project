/**
 * Game Feature - QuestionView
 * Отображение текущего вопроса
 */

import { Card, Button } from '@/shared/ui';
import type { QuestionState, StartMediaPayload } from '@/shared/types';
import { SyncMediaPlayer } from '../SyncMediaPlayer/SyncMediaPlayer';
import './QuestionView.css';

interface QuestionViewProps {
  question: QuestionState;
  canPressButton: boolean;
  onPressButton?: () => void;
  timeRemaining?: number;
  startMedia?: StartMediaPayload | null;
}

export const QuestionView = ({
  question,
  canPressButton,
  onPressButton,
  timeRemaining,
  startMedia,
}: QuestionViewProps) => {
  const hasMedia = question.mediaType && question.mediaType !== 'text';

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

      {hasMedia && (
        <div className="question-view__media">
          <SyncMediaPlayer
            startMedia={startMedia || null}
            fallbackUrl={question.mediaUrl}
            autoPlay={true}
            muted={false}
            controls={true}
          />
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
