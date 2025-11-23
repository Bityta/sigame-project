/**
 * Game Feature - GameBoard
 * Игровое поле с темами и вопросами
 */

import { Card } from '@/shared/ui';
import type { ThemeState, QuestionState } from '@/shared/types';
import './GameBoard.css';

interface GameBoardProps {
  themes?: ThemeState[];
  onQuestionSelect?: (themeId: string, questionId: string) => void;
  canSelectQuestion: boolean;
}

export const GameBoard = ({
  themes,
  onQuestionSelect,
  canSelectQuestion,
}: GameBoardProps) => {
  if (!themes || themes.length === 0) {
    return (
      <Card className="game-board-empty">
        <p>Загрузка игрового поля...</p>
      </Card>
    );
  }

  return (
    <div className="game-board">
      {themes.map((theme, themeIndex) => (
        <div key={themeIndex} className="game-board__theme">
          <div className="game-board__theme-name">{theme.name}</div>
          <div className="game-board__questions">
            {theme.questions.map((question, questionIndex) => (
              <button
                key={questionIndex}
                className={`game-board__question ${
                  !question.available ? 'game-board__question--disabled' : ''
                }`}
                onClick={() =>
                  question.available &&
                  canSelectQuestion &&
                  onQuestionSelect?.(theme.name, question.id)
                }
                disabled={!question.available || !canSelectQuestion}
              >
                {question.price}
              </button>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};

