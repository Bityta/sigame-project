import { useParams, useNavigate } from 'react-router-dom';
import { useGameWebSocket } from '@/entities/game';
import { useCurrentUser } from '@/entities/user';
import { GameBoard, PlayerList, QuestionView, RoundsOverview, RoundIntro, GameEnd } from '@/features/game';
import { Button, Spinner } from '@/shared/ui';
import { ROUTES, TEXTS } from '@/shared/config';
import './GamePage.css';

export const GamePage = () => {
  const { gameId } = useParams<{ gameId: string }>();
  const navigate = useNavigate();
  const { data: user } = useCurrentUser();

  const {
    isConnected,
    gameState,
    selectQuestion,
    pressButton,
    judgeAnswer,
  } = useGameWebSocket({
    gameId: gameId!,
    userId: user?.id || '',
    onError: (error) => {
      console.error('Game error:', error);
    },
  });

  if (!isConnected || !gameState) {
    return (
      <div className="game-page">
        <div className="game-page__connecting">
          <Spinner size="large" center />
          <p>{TEXTS.GAME.CONNECTING}</p>
        </div>
      </div>
    );
  }

  const currentPlayer = gameState.players.find((p) => p.userId === user?.id);
  const isHost = currentPlayer?.role === 'host';
  
  // Only host can select questions
  const canSelectQuestion = gameState.status === 'question_select' && isHost;
  
  // Only players (not host) can press button
  const canPressButton = gameState.status === 'button_press' && !isHost;
  
  // Host judges answers
  const canJudgeAnswer = gameState.status === 'answer_judging' && isHost;

  const handleQuestionSelect = (themeId: string, questionId: string) => {
    selectQuestion(themeId, questionId);
  };

  const handleLeaveGame = () => {
    navigate(ROUTES.LOBBY);
  };

  // Determine turn indicator text
  const getTurnIndicator = () => {
    switch (gameState.status) {
      case 'question_select':
        return isHost ? 'Выберите вопрос' : 'Ведущий выбирает вопрос...';
      case 'button_press':
        return isHost ? 'Ждём, пока игрок нажмёт кнопку...' : 'Жмите кнопку!';
      case 'answer_judging':
        return isHost ? 'Оцените ответ игрока' : 'Ждём решения ведущего...';
      default:
        return '';
    }
  };

  return (
    <div className="game-page">
      {/* Header */}
      <header className="game-page__header">
        <div className="game-page__header-left">
          <div className="game-page__header-info">
            <h1 className="game-page__title">{TEXTS.GAME.ROUND(gameState.currentRound)}</h1>
            {gameState.roundName && (
              <p className="game-page__round-name">{gameState.roundName}</p>
            )}
          </div>
          {/* Role Indicator */}
          <div className={`game-page__role-indicator ${isHost ? 'game-page__role-indicator--host' : 'game-page__role-indicator--player'}`}>
            {isHost ? '👑 Ведущий' : '🎮 Игрок'}
          </div>
          {/* Global Timer */}
          {gameState.timeRemaining !== undefined && gameState.timeRemaining > 0 && (
            <div className={`game-page__header-timer ${
              gameState.timeRemaining <= 3 ? 'game-page__header-timer--danger' : 
              gameState.timeRemaining <= 5 ? 'game-page__header-timer--warning' : ''
            }`}>
              <span className="game-page__header-timer-icon">⏱</span>
              <span>{gameState.timeRemaining}с</span>
            </div>
          )}
        </div>
        <Button variant="danger" size="small" onClick={handleLeaveGame}>
          {TEXTS.GAME.LEAVE_GAME}
        </Button>
      </header>

      {/* Main Content - New Layout */}
      <div className="game-page__content">
        {/* Players Panel - Top */}
        <PlayerList
          players={gameState.players}
          activePlayer={gameState.activePlayer}
          currentUserId={user?.id}
        />

        {/* Turn Indicator with Timer Bar */}
        {getTurnIndicator() && (
          <div className="game-page__turn-indicator-wrapper">
            <span className="game-page__turn-indicator-text">{getTurnIndicator()}</span>
            {gameState.status === 'question_select' && gameState.timeRemaining !== undefined && (
              <div className="game-page__timer-bar">
                <div 
                  className={`game-page__timer-bar-fill ${
                    gameState.timeRemaining <= 3 ? 'game-page__timer-bar-fill--danger' :
                    gameState.timeRemaining <= 5 ? 'game-page__timer-bar-fill--warning' : ''
                  }`}
                  style={{ 
                    width: `${Math.min(100, (gameState.timeRemaining / 60) * 100)}%` 
                  }}
                />
              </div>
            )}
          </div>
        )}

        {/* Waiting Screen - kept for backward compatibility but should not appear */}
        {gameState.status === 'waiting' && (
          <div className="game-page__waiting">
            <h2>{TEXTS.GAME.WAITING_PLAYERS}</h2>
            <p>Подготовка игры...</p>
          </div>
        )}

        {/* Rounds Overview */}
        {gameState.status === 'rounds_overview' && gameState.allRounds && (
          <RoundsOverview rounds={gameState.allRounds} />
        )}

        {/* Round Intro */}
        {gameState.status === 'round_start' && (
          <RoundIntro
            roundNumber={gameState.currentRound}
            roundName={gameState.roundName}
          />
        )}

        {/* Game Board */}
        {gameState.status === 'question_select' && (
          <GameBoard
            themes={gameState.themes}
            onQuestionSelect={handleQuestionSelect}
            canSelectQuestion={canSelectQuestion}
          />
        )}

        {/* Question View */}
        {gameState.currentQuestion && (
          <QuestionView
            question={gameState.currentQuestion}
            canPressButton={canPressButton}
            onPressButton={pressButton}
            timeRemaining={gameState.timeRemaining}
            isHost={isHost}
          />
        )}

        {/* Judging Panel (for host) */}
        {canJudgeAnswer && gameState.activePlayer && (
          <div className="game-page__judging">
            <div className="game-page__judging-player">
              <span className="game-page__judging-label">Отвечает</span>
              <span className="game-page__judging-name">
                {gameState.players.find(p => p.userId === gameState.activePlayer)?.username}
              </span>
            </div>
            {/* Show correct answer to host */}
            {gameState.currentQuestion?.answer && (
              <div className="game-page__judging-answer">
                <span className="game-page__judging-answer-label">Правильный ответ:</span>
                <span className="game-page__judging-answer-text">{gameState.currentQuestion.answer}</span>
              </div>
            )}
            <p className="game-page__judging-hint">Игрок говорит ответ вслух. Оцените:</p>
            <div className="game-page__judging-buttons">
              <button 
                className="game-page__judge-btn game-page__judge-btn--correct"
                onClick={() => judgeAnswer(gameState.activePlayer!, true)}
              >
                ✓ Верно
              </button>
              <button 
                className="game-page__judge-btn game-page__judge-btn--wrong"
                onClick={() => judgeAnswer(gameState.activePlayer!, false)}
              >
                ✗ Неверно
              </button>
            </div>
          </div>
        )}

        {/* Waiting for Host (for players) */}
        {gameState.status === 'answer_judging' && !isHost && (
          <div className="game-page__waiting-host">
            <div className="game-page__waiting-host-icon">🎤</div>
            <p>Скажите ответ вслух!</p>
            <p className="game-page__waiting-host-hint">Ждём решения ведущего...</p>
          </div>
        )}

        {/* Game End */}
        {gameState.status === 'game_end' && gameState.winners && gameState.finalScores && (
          <GameEnd
            winners={gameState.winners}
            finalScores={gameState.finalScores}
            currentUserId={user?.id}
          />
        )}

        {/* Message */}
        {gameState.message && (
          <div className="game-page__message">{gameState.message}</div>
        )}
      </div>
    </div>
  );
};
