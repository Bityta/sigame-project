/**
 * Game Feature - GameEnd
 * Экран окончания игры с поздравлением победителей
 */

import { useNavigate } from 'react-router-dom';
import type { PlayerScore } from '@/shared/types';
import { Button } from '@/shared/ui';
import { ROUTES } from '@/shared/config';
import './GameEnd.css';

interface GameEndProps {
  winners: PlayerScore[];
  finalScores: PlayerScore[];
  currentUserId?: string;
}

export const GameEnd = ({ winners, finalScores, currentUserId }: GameEndProps) => {
  const navigate = useNavigate();
  const topWinner = winners[0];
  const isCurrentUserWinner = winners.some(w => w.userId === currentUserId);

  const getMedalEmoji = (rank: number) => {
    switch (rank) {
      case 1: return '🥇';
      case 2: return '🥈';
      case 3: return '🥉';
      default: return `#${rank}`;
    }
  };

  const getPlaceText = (rank: number) => {
    switch (rank) {
      case 1: return '1-е место';
      case 2: return '2-е место';
      case 3: return '3-е место';
      default: return `${rank}-е место`;
    }
  };

  return (
    <div className="game-end">
      {/* Confetti effect */}
      <div className="game-end__confetti">
        {Array.from({ length: 50 }).map((_, i) => (
          <div
            key={i}
            className="game-end__confetti-piece"
            style={{
              left: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 3}s`,
              backgroundColor: ['#fbbf24', '#ef4444', '#3b82f6', '#10b981', '#8b5cf6'][Math.floor(Math.random() * 5)],
            }}
          />
        ))}
      </div>

      {/* Header */}
      <div className="game-end__header">
        <h1 className="game-end__title">🎉 Игра завершена! 🎉</h1>
        {isCurrentUserWinner && (
          <p className="game-end__congrats">Поздравляем с победой!</p>
        )}
      </div>

      {/* Winner Podium */}
      {winners.length > 0 && (
        <div className="game-end__podium">
          {/* Second place */}
          {winners[1] && (
            <div className="game-end__podium-place game-end__podium-place--second">
              <div className="game-end__podium-medal">{getMedalEmoji(2)}</div>
              <div className="game-end__podium-name">{winners[1].username}</div>
              <div className="game-end__podium-score">{winners[1].score}</div>
              <div className="game-end__podium-bar game-end__podium-bar--second" />
            </div>
          )}
          
          {/* First place */}
          <div className="game-end__podium-place game-end__podium-place--first">
            <div className="game-end__podium-crown">👑</div>
            <div className="game-end__podium-medal">{getMedalEmoji(1)}</div>
            <div className="game-end__podium-name">{topWinner.username}</div>
            <div className="game-end__podium-score">{topWinner.score}</div>
            <div className="game-end__podium-bar game-end__podium-bar--first" />
          </div>
          
          {/* Third place */}
          {winners[2] && (
            <div className="game-end__podium-place game-end__podium-place--third">
              <div className="game-end__podium-medal">{getMedalEmoji(3)}</div>
              <div className="game-end__podium-name">{winners[2].username}</div>
              <div className="game-end__podium-score">{winners[2].score}</div>
              <div className="game-end__podium-bar game-end__podium-bar--third" />
            </div>
          )}
        </div>
      )}

      {/* Full Scoreboard */}
      <div className="game-end__scoreboard">
        <h2 className="game-end__scoreboard-title">Итоговая таблица</h2>
        <div className="game-end__scoreboard-list">
          {finalScores.map((player) => (
            <div
              key={player.userId}
              className={`game-end__scoreboard-row ${
                player.userId === currentUserId ? 'game-end__scoreboard-row--you' : ''
              } ${player.rank <= 3 ? `game-end__scoreboard-row--top${player.rank}` : ''}`}
            >
              <span className="game-end__scoreboard-rank">
                {getMedalEmoji(player.rank)}
              </span>
              <span className="game-end__scoreboard-name">
                {player.username}
                {player.userId === currentUserId && <span className="game-end__you-badge">Вы</span>}
              </span>
              <span className={`game-end__scoreboard-score ${player.score < 0 ? 'game-end__scoreboard-score--negative' : ''}`}>
                {player.score}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* Actions */}
      <div className="game-end__actions">
        <Button
          variant="primary"
          size="large"
          onClick={() => navigate(ROUTES.LOBBY)}
        >
          Вернуться в лобби
        </Button>
      </div>
    </div>
  );
};

