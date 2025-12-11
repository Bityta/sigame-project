/**
 * ForAllResults - Shows results of forAll question for all players
 */

import React from 'react';
import type { ForAllAnswerResult } from '@/shared/types';
import { Card } from '@/shared/ui';
import './ForAllResults.css';

interface ForAllResultsProps {
  results: ForAllAnswerResult[];
  correctAnswer: string;
}

export const ForAllResults: React.FC<ForAllResultsProps> = ({
  results,
  correctAnswer,
}) => {
  // Sort by score delta descending (correct answers first)
  const sortedResults = [...results].sort((a, b) => b.scoreDelta - a.scoreDelta);

  return (
    <Card className="for-all-results" padding="large">
      <div className="for-all-results__header">
        <span className="for-all-results__icon">📊</span>
        <h2 className="for-all-results__title">Результаты</h2>
      </div>

      <div className="for-all-results__correct-answer">
        <span className="for-all-results__correct-label">Правильный ответ:</span>
        <span className="for-all-results__correct-value">{correctAnswer}</span>
      </div>

      <div className="for-all-results__list">
        {sortedResults.map((result) => (
          <div 
            key={result.userId} 
            className={`for-all-results__item ${
              result.isCorrect 
                ? 'for-all-results__item--correct' 
                : 'for-all-results__item--wrong'
            }`}
          >
            <div className="for-all-results__item-icon">
              {result.isCorrect ? '✅' : '❌'}
            </div>
            <div className="for-all-results__item-info">
              <span className="for-all-results__item-name">{result.username}</span>
              <span className="for-all-results__item-answer">
                {result.answer || '(нет ответа)'}
              </span>
            </div>
            <div className={`for-all-results__item-delta ${
              result.scoreDelta >= 0 
                ? 'for-all-results__item-delta--positive' 
                : 'for-all-results__item-delta--negative'
            }`}>
              {result.scoreDelta >= 0 ? '+' : ''}{result.scoreDelta}
            </div>
          </div>
        ))}
      </div>

      {results.length === 0 && (
        <p className="for-all-results__no-answers">
          Никто не успел ответить на вопрос
        </p>
      )}
    </Card>
  );
};

