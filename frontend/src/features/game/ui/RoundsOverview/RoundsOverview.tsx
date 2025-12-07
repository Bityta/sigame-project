/**
 * Game Feature - RoundsOverview
 * Отображение обзора всех раундов перед началом игры
 */

import type { RoundOverview } from '@/shared/types';
import './RoundsOverview.css';

interface RoundsOverviewProps {
  rounds: RoundOverview[];
}

export const RoundsOverview = ({ rounds }: RoundsOverviewProps) => {
  return (
    <div className="rounds-overview">
      <div className="rounds-overview__header">
        <h1 className="rounds-overview__title">Обзор игры</h1>
        <p className="rounds-overview__subtitle">Сегодня вас ждут следующие раунды:</p>
      </div>

      <div className="rounds-overview__rounds">
        {rounds.map((round, index) => (
          <div
            key={round.roundNumber}
            className="rounds-overview__round"
            style={{ animationDelay: `${index * 0.15}s` }}
          >
            <div className="rounds-overview__round-header">
              <span className="rounds-overview__round-number">Раунд {round.roundNumber}</span>
              <h2 className="rounds-overview__round-name">{round.name}</h2>
            </div>
            <div className="rounds-overview__themes">
              {round.themeNames.map((theme, themeIndex) => (
                <span
                  key={themeIndex}
                  className="rounds-overview__theme"
                  style={{ animationDelay: `${index * 0.15 + themeIndex * 0.05}s` }}
                >
                  {theme}
                </span>
              ))}
            </div>
          </div>
        ))}
      </div>

      <div className="rounds-overview__footer">
        <div className="rounds-overview__countdown">
          <span className="rounds-overview__countdown-text">Игра начнётся автоматически...</span>
          <div className="rounds-overview__progress-bar">
            <div className="rounds-overview__progress-fill"></div>
          </div>
        </div>
      </div>
    </div>
  );
};

