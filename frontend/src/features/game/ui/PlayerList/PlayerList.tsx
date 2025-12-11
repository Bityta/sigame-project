/**
 * Game Feature - PlayerList
 * Ведущий слева, игроки справа
 */

import type { PlayerState } from '@/shared/types';
import './PlayerList.css';

interface PlayerListProps {
  players: PlayerState[];
  activePlayer?: string;
  currentUserId?: string;
}

export const PlayerList = ({ players, activePlayer, currentUserId }: PlayerListProps) => {
  // Отделяем ведущего от игроков
  const host = players.find(p => p.role === 'host');
  const gamePlayers = players
    .filter(p => p.role !== 'host')
    .sort((a, b) => b.score - a.score);

  // Рассчитываем ранги
  const getRank = (index: number) => index + 1;

  const getRankClass = (rank: number) => {
    if (rank === 1) return 'player-rank--1';
    if (rank === 2) return 'player-rank--2';
    if (rank === 3) return 'player-rank--3';
    return 'player-rank--other';
  };

  const getScoreClass = (score: number) => {
    if (score < 0) return 'player-score--negative';
    if (score === 0) return 'player-score--zero';
    return '';
  };

  return (
    <div className="players-panel">
      {/* Ведущий слева */}
      {host && (
        <div className="players-panel__host">
          <div className={`player-card player-card--host ${currentUserId === host.userId ? 'player-card--you' : ''}`}>
            <div className="player-avatar player-avatar--host">👑</div>
            <div className="player-info">
              <div className="player-name">{host.username}</div>
              <div className="player-status">Ведущий</div>
            </div>
          </div>
        </div>
      )}

      {/* Игроки справа */}
      <div className="players-panel__players">
        {gamePlayers.map((player, index) => {
          const rank = getRank(index);
          const isAnswering = player.userId === activePlayer;
          
          return (
            <div
              key={player.userId}
              className={`player-card ${
                isAnswering ? 'player-card--answering' : ''
              } ${currentUserId === player.userId ? 'player-card--you' : ''}`}
            >
              {/* Бейдж позиции */}
              <span className={`player-rank ${getRankClass(rank)}`}>
                {rank}
              </span>
              
              <div className="player-avatar player-avatar--player">
                {player.username.substring(0, 2).toUpperCase()}
              </div>
              <div className="player-info">
                <div className="player-name">{player.username}</div>
                <div className={`player-status ${isAnswering ? 'player-status--answering' : ''}`}>
                  {isAnswering ? '🎤 Отвечает!' : 'Игрок'}
                </div>
              </div>
              <span className={`player-score ${getScoreClass(player.score)}`}>
                {player.score}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
};
