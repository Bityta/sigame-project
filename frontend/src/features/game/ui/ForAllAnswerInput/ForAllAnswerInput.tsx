/**
 * ForAllAnswerInput - Input panel for players to answer forAll question
 */

import React, { useState } from 'react';
import { Card, Button } from '@/shared/ui';
import './ForAllAnswerInput.css';

interface ForAllAnswerInputProps {
  onSubmit: (answer: string) => void;
  timeRemaining?: number;
  hasSubmitted: boolean;
  isHost: boolean;
}

export const ForAllAnswerInput: React.FC<ForAllAnswerInputProps> = ({
  onSubmit,
  timeRemaining,
  hasSubmitted,
  isHost,
}) => {
  const [answer, setAnswer] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (answer.trim()) {
      onSubmit(answer.trim());
    }
  };

  // Host doesn't participate
  if (isHost) {
    return (
      <Card className="for-all-answer-input for-all-answer-input--host" padding="large">
        <div className="for-all-answer-input__header">
          <span className="for-all-answer-input__icon">👥</span>
          <h2 className="for-all-answer-input__title">Вопрос для всех</h2>
        </div>
        <p className="for-all-answer-input__description">
          Игроки отвечают на вопрос...
        </p>
        {timeRemaining !== undefined && (
          <div className="for-all-answer-input__timer">
            Осталось времени: <span className="for-all-answer-input__time">{timeRemaining}с</span>
          </div>
        )}
      </Card>
    );
  }

  if (hasSubmitted) {
    return (
      <Card className="for-all-answer-input for-all-answer-input--submitted" padding="large">
        <div className="for-all-answer-input__header">
          <span className="for-all-answer-input__icon">✅</span>
          <h2 className="for-all-answer-input__title">Ответ отправлен!</h2>
        </div>
        <p className="for-all-answer-input__description">
          Ожидаем других игроков...
        </p>
        {timeRemaining !== undefined && (
          <div className="for-all-answer-input__timer">
            Осталось времени: <span className="for-all-answer-input__time">{timeRemaining}с</span>
          </div>
        )}
      </Card>
    );
  }

  return (
    <Card className="for-all-answer-input" padding="large">
      <div className="for-all-answer-input__header">
        <span className="for-all-answer-input__icon">👥</span>
        <h2 className="for-all-answer-input__title">Вопрос для всех</h2>
      </div>

      <p className="for-all-answer-input__description">
        Напишите свой ответ. Все игроки отвечают одновременно!
      </p>

      {timeRemaining !== undefined && (
        <div className="for-all-answer-input__timer">
          Осталось времени: <span className="for-all-answer-input__time">{timeRemaining}с</span>
        </div>
      )}

      <form onSubmit={handleSubmit} className="for-all-answer-input__form">
        <input
          type="text"
          value={answer}
          onChange={(e) => setAnswer(e.target.value)}
          placeholder="Введите ваш ответ..."
          className="for-all-answer-input__input"
          autoFocus
        />
        <Button 
          type="submit" 
          variant="primary" 
          size="large"
          disabled={!answer.trim()}
        >
          Отправить ответ
        </Button>
      </form>
    </Card>
  );
};

