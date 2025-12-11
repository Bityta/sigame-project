/**
 * StakeBettingPanel - Panel for player to place a stake on stake question
 */

import React, { useState } from 'react';
import type { StakeInfo } from '@/shared/types';
import { Card, Button } from '@/shared/ui';
import './StakeBettingPanel.css';

interface StakeBettingPanelProps {
  stakeInfo: StakeInfo;
  playerScore: number;
  onPlaceStake: (amount: number, allIn: boolean) => void;
  timeRemaining?: number;
  isActivePlayer: boolean;
}

export const StakeBettingPanel: React.FC<StakeBettingPanelProps> = ({
  stakeInfo,
  playerScore,
  onPlaceStake,
  timeRemaining,
  isActivePlayer,
}) => {
  const [selectedAmount, setSelectedAmount] = useState(stakeInfo.minBet);

  // Quick bet options
  const quickBets = [
    { label: 'Минимум', value: stakeInfo.minBet },
    { label: '×2', value: Math.min(stakeInfo.minBet * 2, stakeInfo.maxBet) },
    { label: '×3', value: Math.min(stakeInfo.minBet * 3, stakeInfo.maxBet) },
    { label: 'Половина', value: Math.min(Math.floor(playerScore / 2), stakeInfo.maxBet) },
  ].filter((bet, idx, arr) => {
    // Remove duplicates
    return arr.findIndex(b => b.value === bet.value) === idx;
  });

  const handleSliderChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedAmount(parseInt(e.target.value, 10));
  };

  const handlePlaceBet = () => {
    onPlaceStake(selectedAmount, false);
  };

  const handleAllIn = () => {
    onPlaceStake(playerScore, true);
  };

  if (!isActivePlayer) {
    return (
      <Card className="stake-betting-panel stake-betting-panel--waiting" padding="large">
        <div className="stake-betting-panel__header">
          <span className="stake-betting-panel__icon">💰</span>
          <h2 className="stake-betting-panel__title">Ва-банк!</h2>
        </div>
        <p className="stake-betting-panel__description">
          Ожидаем, пока игрок сделает ставку...
        </p>
      </Card>
    );
  }

  return (
    <Card className="stake-betting-panel" padding="large">
      <div className="stake-betting-panel__header">
        <span className="stake-betting-panel__icon">💰</span>
        <h2 className="stake-betting-panel__title">Ва-банк!</h2>
      </div>

      <p className="stake-betting-panel__description">
        Сделайте ставку на этот вопрос.
        <br />
        Если ответите правильно — получите ставку, неправильно — потеряете.
      </p>

      {timeRemaining !== undefined && (
        <div className="stake-betting-panel__timer">
          Осталось времени: <span className="stake-betting-panel__time">{timeRemaining}с</span>
        </div>
      )}

      <div className="stake-betting-panel__info">
        <div className="stake-betting-panel__info-item">
          <span className="stake-betting-panel__info-label">Ваш счёт:</span>
          <span className="stake-betting-panel__info-value">{playerScore}</span>
        </div>
        <div className="stake-betting-panel__info-item">
          <span className="stake-betting-panel__info-label">Мин. ставка:</span>
          <span className="stake-betting-panel__info-value">{stakeInfo.minBet}</span>
        </div>
        <div className="stake-betting-panel__info-item">
          <span className="stake-betting-panel__info-label">Макс. ставка:</span>
          <span className="stake-betting-panel__info-value">{stakeInfo.maxBet}</span>
        </div>
      </div>

      <div className="stake-betting-panel__quick-bets">
        {quickBets.map((bet) => (
          <Button
            key={bet.label}
            variant={selectedAmount === bet.value ? 'primary' : 'secondary'}
            size="medium"
            onClick={() => setSelectedAmount(bet.value)}
          >
            {bet.label} ({bet.value})
          </Button>
        ))}
      </div>

      <div className="stake-betting-panel__slider-container">
        <input
          type="range"
          min={stakeInfo.minBet}
          max={stakeInfo.maxBet}
          value={selectedAmount}
          onChange={handleSliderChange}
          className="stake-betting-panel__slider"
        />
        <div className="stake-betting-panel__selected-amount">
          Ваша ставка: <span className="stake-betting-panel__amount">{selectedAmount}</span>
        </div>
      </div>

      <div className="stake-betting-panel__actions">
        <Button variant="primary" size="large" onClick={handlePlaceBet}>
          Сделать ставку {selectedAmount}
        </Button>
        {playerScore > 0 && (
          <Button variant="secondary" size="large" onClick={handleAllIn}>
            Ва-банк! ({playerScore})
          </Button>
        )}
      </div>
    </Card>
  );
};

