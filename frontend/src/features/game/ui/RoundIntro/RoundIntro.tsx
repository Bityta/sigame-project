/**
 * Game Feature - RoundIntro
 * Отображение названия раунда перед показом вопросов
 */

import './RoundIntro.css';

interface RoundIntroProps {
  roundNumber: number;
  roundName?: string;
}

export const RoundIntro = ({ roundNumber, roundName }: RoundIntroProps) => {
  return (
    <div className="round-intro">
      <div className="round-intro__content">
        <div className="round-intro__number">
          <span className="round-intro__number-label">Раунд</span>
          <span className="round-intro__number-value">{roundNumber}</span>
        </div>
        {roundName && (
          <h1 className="round-intro__name">{roundName}</h1>
        )}
        <div className="round-intro__decoration">
          <div className="round-intro__line round-intro__line--left"></div>
          <div className="round-intro__star">★</div>
          <div className="round-intro__line round-intro__line--right"></div>
        </div>
      </div>
    </div>
  );
};

