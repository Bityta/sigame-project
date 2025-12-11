/**
 * SecretTransferPanel - Panel for host to transfer secret question to a player
 */

import React from 'react';
import type { PlayerState } from '@/shared/types';
import { Card, Button } from '@/shared/ui';
import './SecretTransferPanel.css';

interface SecretTransferPanelProps {
  players: PlayerState[];
  onTransfer: (targetUserId: string) => void;
  timeRemaining?: number;
  currentUserId?: string;
}

export const SecretTransferPanel: React.FC<SecretTransferPanelProps> = ({
  players,
  onTransfer,
  timeRemaining,
  currentUserId,
}) => {
  // Filter out host from available players
  const availablePlayers = players.filter(
    (p) => p.role !== 'host' && p.isActive
  );

  return (
    <Card className="secret-transfer-panel" padding="large">
      <div className="secret-transfer-panel__header">
        <span className="secret-transfer-panel__icon">🐱</span>
        <h2 className="secret-transfer-panel__title">Кот в мешке!</h2>
      </div>

      <p className="secret-transfer-panel__description">
        Выберите игрока, которому передать этот вопрос.
        <br />
        Выбранный игрок обязан ответить на вопрос.
      </p>

      {timeRemaining !== undefined && (
        <div className="secret-transfer-panel__timer">
          Осталось времени: <span className="secret-transfer-panel__time">{timeRemaining}с</span>
        </div>
      )}

      <div className="secret-transfer-panel__players">
        {availablePlayers.map((player) => (
          <Button
            key={player.userId}
            variant="primary"
            size="large"
            className="secret-transfer-panel__player-btn"
            onClick={() => onTransfer(player.userId)}
            disabled={player.userId === currentUserId}
          >
            <span className="secret-transfer-panel__player-name">{player.username}</span>
            <span className="secret-transfer-panel__player-score">{player.score} очков</span>
          </Button>
        ))}
      </div>
    </Card>
  );
};

