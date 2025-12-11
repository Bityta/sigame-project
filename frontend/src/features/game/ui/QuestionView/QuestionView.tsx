/**
 * Game Feature - QuestionView
 * Отображение текущего вопроса с большой кнопкой
 */

import { useState } from 'react';
import { Card } from '@/shared/ui';
import type { QuestionState, StartMediaPayload } from '@/shared/types';
import { SyncMediaPlayer } from '../SyncMediaPlayer/SyncMediaPlayer';
import './QuestionView.css';

interface QuestionViewProps {
  question: QuestionState;
  canPressButton: boolean;
  onPressButton?: () => void;
  timeRemaining?: number;
  isHost?: boolean;
  hideAnswer?: boolean; // Hide answer when judging panel is shown
  startMedia?: StartMediaPayload | null;
}

export const QuestionView = ({
  question,
  canPressButton,
  onPressButton,
  timeRemaining,
  isHost = false,
  hideAnswer = false,
  startMedia,
}: QuestionViewProps) => {
  const isTimerWarning = timeRemaining !== undefined && timeRemaining <= 5;
  const isTimerDanger = timeRemaining !== undefined && timeRemaining <= 3;
  const hasMedia = question.mediaType && question.mediaType !== 'text' && question.mediaUrl;

  return (
    <Card className="question-view" padding="large">
      <div className="question-view__header">
        <div className="question-view__price">{question.price} очков</div>
        {/* Always show timer container to prevent layout shift */}
        <div className={`question-view__timer ${isTimerDanger ? 'question-view__timer--danger' : isTimerWarning ? 'question-view__timer--warning' : ''}`}>
          {timeRemaining !== undefined ? `${timeRemaining}с` : '\u00A0'}
        </div>
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

      {/* Show correct answer to host (hide when judging panel is shown) */}
      {isHost && question.answer && !hideAnswer && (
        <div className="question-view__answer">
          <span className="question-view__answer-label">Правильный ответ:</span>
          <span className="question-view__answer-text">{question.answer}</span>
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
