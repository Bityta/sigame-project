/**
 * Game Feature - PlayerList
 * Список игроков с очками
 */

import { Card } from '@/shared/ui';
import type { PlayerState } from '@/shared/types';
import './PlayerList.css';

interface PlayerListProps {
  players: PlayerState[];
  activePlayer?: string;
}

export const PlayerList = ({ players, activePlayer }: PlayerListProps) => {
  // Сортируем игроков по очкам
  const sortedPlayers = [...players].sort((a, b) => b.score - a.score);

  return (
    <Card className="player-list" padding="medium">
      <h3 className="player-list__title">Игроки</h3>
      <div className="player-list__items">
        {sortedPlayers.map((player) => (
          <div
            key={player.userId}
            className={`player-item ${
              player.userId === activePlayer ? 'player-item--active' : ''
            } ${player.isReady ? 'player-item--ready' : ''}`}
          >
            <div className="player-item__info">
              <span className="player-item__name">{player.username}</span>
              {player.role === 'host' && (
                <span className="player-item__badge">👑</span>
              )}
            </div>
            <div className="player-item__score">{player.score}</div>
          </div>
        ))}
      </div>
    </Card>
  );
};

