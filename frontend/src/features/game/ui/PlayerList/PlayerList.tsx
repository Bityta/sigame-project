/**
 * Game Feature - PlayerList
 * Ведущий слева, игроки справа
 * Показывает аватары и статус подключения
 */

import type { PlayerState } from '@/shared/types';
import './PlayerList.css';

interface PlayerListProps {
  players: PlayerState[];
  activePlayer?: string;
  currentUserId?: string;
}

// Функция для получения инициалов
const getInitials = (name: string) => {
  return name.substring(0, 2).toUpperCase();
};

// Компонент аватара
const Avatar = ({ player, isHost = false }: { player: PlayerState; isHost?: boolean }) => {
  const hasAvatar = player.avatarUrl && player.avatarUrl.length > 0;
  
  if (isHost) {
    return (
      <div className={`player-avatar player-avatar--host ${!player.isConnected ? 'player-avatar--disconnected' : ''}`}>
        {hasAvatar ? (
          <img src={player.avatarUrl} alt={player.username} className="player-avatar__img" />
        ) : (
          '👑'
        )}
        {!player.isConnected && <span className="player-avatar__offline-badge">⚫</span>}
      </div>
    );
  }
  
  return (
    <div className={`player-avatar player-avatar--player ${!player.isConnected ? 'player-avatar--disconnected' : ''}`}>
      {hasAvatar ? (
        <img src={player.avatarUrl} alt={player.username} className="player-avatar__img" />
      ) : (
        getInitials(player.username)
      )}
      {!player.isConnected && <span className="player-avatar__offline-badge">⚫</span>}
    </div>
  );
};

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

  const getConnectionStatus = (player: PlayerState) => {
    if (!player.isConnected) return 'Отключён';
    return null;
  };

  return (
    <div className="players-panel">
      {/* Ведущий слева */}
      {host && (
        <div className="players-panel__host">
          <div className={`player-card player-card--host ${currentUserId === host.userId ? 'player-card--you' : ''} ${!host.isConnected ? 'player-card--disconnected' : ''}`}>
            <Avatar player={host} isHost />
            <div className="player-info">
              <div className="player-name">{host.username}</div>
              <div className={`player-status ${!host.isConnected ? 'player-status--disconnected' : ''}`}>
                {getConnectionStatus(host) || 'Ведущий'}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Игроки справа */}
      <div className="players-panel__players">
        {gamePlayers.map((player, index) => {
          const rank = getRank(index);
          const isAnswering = player.userId === activePlayer;
          const isDisconnected = !player.isConnected;
          
          return (
            <div
              key={player.userId}
              className={`player-card ${
                isAnswering ? 'player-card--answering' : ''
              } ${currentUserId === player.userId ? 'player-card--you' : ''} ${isDisconnected ? 'player-card--disconnected' : ''}`}
            >
              {/* Бейдж позиции */}
              <span className={`player-rank ${getRankClass(rank)}`}>
                {rank}
              </span>
              
              <Avatar player={player} />
              
              <div className="player-info">
                <div className="player-name">{player.username}</div>
                <div className={`player-status ${isAnswering ? 'player-status--answering' : ''} ${isDisconnected ? 'player-status--disconnected' : ''}`}>
                  {isDisconnected ? '📵 Отключён' : (isAnswering ? '🎤 Отвечает!' : 'Игрок')}
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
