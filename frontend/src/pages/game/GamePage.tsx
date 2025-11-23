import { useParams, useNavigate } from 'react-router-dom';
import { useGameWebSocket } from '@/entities/game';
import { useCurrentUser } from '@/entities/user';
import { GameBoard, PlayerList, QuestionView } from '@/features/game';
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
    sendReady,
    selectQuestion,
    pressButton,
    submitAnswer,
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
  const canSelectQuestion =
    gameState.status === 'question_select' &&
    gameState.activePlayer === user?.id;
  const canPressButton = gameState.status === 'button_press';
  const canAnswer =
    gameState.status === 'answering' &&
    gameState.activePlayer === user?.id;

  const handleQuestionSelect = (themeId: string, questionId: string) => {
    selectQuestion(themeId, questionId);
  };

  const handleReady = () => {
    sendReady();
  };

  const handleLeaveGame = () => {
    navigate(ROUTES.LOBBY);
  };

  return (
    <div className="game-page">
      <header className="game-page__header">
        <div className="game-page__header-info">
          <h1 className="game-page__title">{TEXTS.GAME.ROUND(gameState.currentRound)}</h1>
          {gameState.roundName && (
            <p className="game-page__round-name">{gameState.roundName}</p>
          )}
        </div>
        <Button variant="danger" size="small" onClick={handleLeaveGame}>
          {TEXTS.GAME.LEAVE_GAME}
        </Button>
      </header>

      <div className="game-page__content">
        <aside className="game-page__sidebar">
          <PlayerList
            players={gameState.players}
            activePlayer={gameState.activePlayer}
          />
        </aside>

        <main className="game-page__main">
          {gameState.status === 'waiting' && (
            <div className="game-page__waiting">
              <h2>{TEXTS.GAME.WAITING_PLAYERS}</h2>
              <p>
                {TEXTS.GAME.READY_COUNT(
                  gameState.players.filter((p) => p.isReady).length,
                  gameState.players.length
                )}
              </p>
              {!currentPlayer?.isReady && (
                <Button variant="primary" size="large" onClick={handleReady}>
                  {TEXTS.GAME.READY_BUTTON}
                </Button>
              )}
            </div>
          )}

          {(gameState.status === 'question_select' ||
            gameState.status === 'round_start') && (
            <GameBoard
              themes={gameState.themes}
              onQuestionSelect={handleQuestionSelect}
              canSelectQuestion={canSelectQuestion}
            />
          )}

          {gameState.currentQuestion && (
            <QuestionView
              question={gameState.currentQuestion}
              canPressButton={canPressButton}
              canAnswer={canAnswer}
              onPressButton={pressButton}
              onSubmitAnswer={submitAnswer}
              timeRemaining={gameState.timeRemaining}
            />
          )}

          {gameState.status === 'game_end' && (
            <div className="game-page__game-end">
              <h2>{TEXTS.GAME.GAME_FINISHED}</h2>
              <Button
                variant="primary"
                size="large"
                onClick={handleLeaveGame}
              >
                {TEXTS.GAME.RETURN_TO_LOBBY}
              </Button>
            </div>
          )}

          {gameState.message && (
            <div className="game-page__message">{gameState.message}</div>
          )}
        </main>
      </div>
    </div>
  );
};
