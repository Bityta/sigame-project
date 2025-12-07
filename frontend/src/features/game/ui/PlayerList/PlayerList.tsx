/**
 * Game Feature - PlayerList
 * Горизонтальная панель игроков (новый дизайн)
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
  const gamePlayers = players.filter(p => p.role !== 'host').sort((a, b) => b.score - a.score);

  const getStatusText = (player: PlayerState) => {
    if (player.userId === activePlayer && player.role !== 'host') {
      return '🎤 Отвечает!';
    }
    if (!player.isReady) {
      return 'Не готов';
    }
    return 'Ожидание';
  };

  return (
    <div className="players-panel">
      {/* Карточка ведущего */}
      {host && (
        <div className={`player-card player-card--host ${currentUserId === host.userId ? 'player-card--you' : ''}`}>
          <div className="player-avatar player-avatar--host">👑</div>
          <div className="player-info">
            <div className="player-name">
              {host.username}
              <span className="player-role player-role--host">HOST</span>
            </div>
            <div className="player-status">Ведущий</div>
          </div>
        </div>
      )}

      {/* Карточки игроков */}
      {gamePlayers.map((player) => (
        <div
          key={player.userId}
          className={`player-card ${
            player.userId === activePlayer ? 'player-card--answering' : ''
          } ${currentUserId === player.userId ? 'player-card--you' : ''}`}
        >
          <div className="player-avatar player-avatar--player">
            {player.username.substring(0, 2).toUpperCase()}
          </div>
          <div className="player-info">
            <div className="player-name">{player.username}</div>
            <div className="player-status">{getStatusText(player)}</div>
          </div>
          <span className={`player-score ${player.score < 0 ? 'player-score--negative' : ''}`}>
            {player.score}
          </span>
        </div>
      ))}
    </div>
  );
};
