/**
 * Game Feature - QuestionView
 * Отображение текущего вопроса с большой кнопкой
 */

import { useState } from 'react';
import { Card } from '@/shared/ui';
import type { QuestionState } from '@/shared/types';
import './QuestionView.css';

interface QuestionViewProps {
  question: QuestionState;
  canPressButton: boolean;
  onPressButton?: () => void;
  timeRemaining?: number;
  isHost?: boolean;
  hideAnswer?: boolean; // Hide answer when judging panel is shown
}

export const QuestionView = ({
  question,
  canPressButton,
  onPressButton,
  timeRemaining,
  isHost = false,
  hideAnswer = false,
}: QuestionViewProps) => {
  const isTimerWarning = timeRemaining !== undefined && timeRemaining <= 5;
  const isTimerDanger = timeRemaining !== undefined && timeRemaining <= 3;
  const [imageLoading, setImageLoading] = useState(true);
  const [imageError, setImageError] = useState(false);

  const renderMedia = () => {
    const { mediaType, mediaUrl } = question;

    // No media or text-only question
    if (!mediaType || mediaType === 'text' || !mediaUrl) {
      return null;
    }

    switch (mediaType) {
      case 'image':
        return (
          <div className="question-view__media question-view__media--image">
            {imageLoading && !imageError && (
              <div className="question-view__media-loading">
                <span className="question-view__media-spinner" />
                Загрузка изображения...
              </div>
            )}
            {imageError ? (
              <div className="question-view__media-error">
                Не удалось загрузить изображение
              </div>
            ) : (
              <img
                src={mediaUrl}
                alt="Вопрос"
                className={`question-view__image ${imageLoading ? 'question-view__image--loading' : ''}`}
                onLoad={() => setImageLoading(false)}
                onError={() => {
                  setImageLoading(false);
                  setImageError(true);
                }}
              />
            )}
          </div>
        );

      case 'audio':
        return (
          <div className="question-view__media question-view__media--audio">
            <div className="question-view__media-icon">🎵</div>
            <audio
              src={mediaUrl}
              controls
              autoPlay
              className="question-view__audio"
            >
              Ваш браузер не поддерживает аудио
            </audio>
          </div>
        );

      case 'video':
        return (
          <div className="question-view__media question-view__media--video">
            <video
              src={mediaUrl}
              controls
              autoPlay
              className="question-view__video"
            >
              Ваш браузер не поддерживает видео
            </video>
          </div>
        );

      default:
        return (
          <div className="question-view__media">
            <p>Неизвестный тип медиа: {mediaType}</p>
          </div>
        );
    }
  };

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

      {renderMedia()}

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
