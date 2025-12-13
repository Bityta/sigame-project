import { useParams, useNavigate } from 'react-router-dom';
import { useRef, useEffect, useState } from 'react';
import { useGameWebSocket } from '@/entities/game';
import { useCurrentUser } from '@/entities/user';
import { 
  GameBoard, 
  PlayerList, 
  QuestionView, 
  RoundsOverview, 
  RoundIntro, 
  GameEnd,
  SecretTransferPanel,
  StakeBettingPanel,
  ForAllAnswerInput,
  ForAllResults,
} from '@/features/game';
import { Button, Spinner } from '@/shared/ui';
import { ROUTES, TEXTS } from '@/shared/config';
import type { ForAllResultsPayload, SecretTransferredPayload, StakePlacedPayload } from '@/shared/types';
import './GamePage.css';

export const GamePage = () => {
  const { gameId } = useParams<{ gameId: string }>();
  const navigate = useNavigate();
  const { data: user } = useCurrentUser();
  
  // Track timer for CSS animation sync
  const [timerDuration, setTimerDuration] = useState<number>(10); // Default 10 seconds
  const [timerElapsed, setTimerElapsed] = useState<number>(0); // How much time has passed
  const [timerKey, setTimerKey] = useState<number>(0); // Key to reset animation
  const lastStatusRef = useRef<string>('');
  const maxTimeSeenRef = useRef<number>(10); // Track highest timeRemaining seen

  // State for special question types
  const [forAllResults, setForAllResults] = useState<ForAllResultsPayload | null>(null);
  const [hasSubmittedForAll, setHasSubmittedForAll] = useState(false);

  const {
    isConnected,
    gameState,
    startMedia,
    selectQuestion,
    pressButton,
    judgeAnswer,
    transferSecret,
    placeStake,
    submitForAllAnswer,
    subscribe,
  } = useGameWebSocket({
    gameId: gameId!,
    userId: user?.id || '',
    onError: (error) => {
      console.error('Game error:', error);
    },
  });

  // Subscribe to special question type events
  useEffect(() => {
    const unsubForAllResults = subscribe<ForAllResultsPayload>('FOR_ALL_RESULTS', (payload) => {
      setForAllResults(payload);
    });

    const unsubSecretTransferred = subscribe<SecretTransferredPayload>('SECRET_TRANSFERRED', (payload) => {
      console.log('Secret transferred:', payload);
    });

    const unsubStakePlaced = subscribe<StakePlacedPayload>('STAKE_PLACED', (payload) => {
      console.log('Stake placed:', payload);
    });

    return () => {
      unsubForAllResults();
      unsubSecretTransferred();
      unsubStakePlaced();
    };
  }, [subscribe]);

  // Reset forAll state when status changes
  useEffect(() => {
    if (gameState?.status !== 'for_all_answering' && gameState?.status !== 'for_all_results') {
      setHasSubmittedForAll(false);
    }
    if (gameState?.status !== 'for_all_results') {
      setForAllResults(null);
    }
  }, [gameState?.status]);
  
  // Start CSS animation when entering question_select or button_press
  useEffect(() => {
    const currentTime = gameState?.timeRemaining || 0;
    const currentStatus = gameState?.status || '';
    
    // Track phases that need timer animation
    if (currentStatus === 'question_select' || currentStatus === 'button_press') {
      if (lastStatusRef.current !== currentStatus) {
        // New phase - start animation
        maxTimeSeenRef.current = currentTime;
        setTimerDuration(currentTime);
        setTimerElapsed(0);
        setTimerKey(prev => prev + 1); // Reset animation
      } else if (currentTime > maxTimeSeenRef.current) {
        // We saw a higher time, update max
        maxTimeSeenRef.current = currentTime;
      }
    }
    lastStatusRef.current = currentStatus;
  }, [gameState?.status, gameState?.timeRemaining]);

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
  
  // Debug logging for answer_judging state
  useEffect(() => {
    if (gameState.status === 'answer_judging') {
      console.log('[GamePage] answer_judging state:', {
        status: gameState.status,
        isHost,
        activePlayer: gameState.activePlayer,
        activePlayerType: typeof gameState.activePlayer,
        activePlayerTruthy: !!gameState.activePlayer,
        canJudgeAnswer,
        user: user?.id,
        condition: canJudgeAnswer && gameState.activePlayer
      });
    }
  }, [gameState.status, gameState.activePlayer, isHost, canJudgeAnswer, user?.id]);

  // Debug logging for question_select state
  useEffect(() => {
    if (gameState.status === 'question_select') {
      console.log('[GamePage] question_select state:', {
        status: gameState.status,
        isHost,
        themes: gameState.themes?.length || 0,
        hasThemes: !!gameState.themes,
        canSelectQuestion
      });
    }
  }, [gameState.status, gameState.themes, isHost, canSelectQuestion]);

  const handleQuestionSelect = (themeId: string, questionId: string) => {
    selectQuestion(themeId, questionId);
  };

  const handleLeaveGame = () => {
    navigate(ROUTES.LOBBY);
  };

  // Handler for submitting forAll answer
  const handleSubmitForAllAnswer = (answer: string) => {
    submitForAllAnswer(answer);
    setHasSubmittedForAll(true);
  };

  // Determine turn indicator text
  const getTurnIndicator = () => {
    switch (gameState.status) {
      case 'question_select':
        return isHost ? 'Выберите вопрос' : 'Ведущий выбирает вопрос...';
      case 'button_press':
        return isHost ? 'Ждём, пока игрок нажмёт кнопку...' : 'Жмите кнопку!';
      case 'answering':
        const isActivePlayer = gameState.activePlayer === user?.id;
        return isHost ? 'Игрок отвечает...' : (isActivePlayer ? 'Говорите ваш ответ!' : 'Ждём ответа игрока...');
      case 'answer_judging':
        return isHost ? 'Оцените ответ игрока' : 'Ждём решения ведущего...';
      case 'secret_transfer':
        return isHost ? 'Выберите игрока для передачи вопроса' : 'Кот в мешке! Ждём выбора ведущего...';
      case 'stake_betting':
        const isActiveForStake = gameState.activePlayer === user?.id;
        return isActiveForStake ? 'Сделайте ставку!' : 'Ждём ставку игрока...';
      case 'for_all_answering':
        return isHost ? 'Игроки отвечают...' : 'Введите ваш ответ!';
      case 'for_all_results':
        return 'Результаты';
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

        {/* Turn Indicator with Timer Bar - always visible to prevent layout shift */}
        <div className="game-page__turn-indicator-wrapper">
          <span className="game-page__turn-indicator-text">
            {getTurnIndicator() || '\u00A0'}
          </span>
          {(gameState.status === 'question_select' || gameState.status === 'button_press') && timerDuration > 0 && (
            <div className="game-page__timer-bar">
              <div 
                key={timerKey}
                className={`game-page__timer-bar-fill ${
                  (gameState.timeRemaining ?? 0) <= 3 ? 'game-page__timer-bar-fill--danger' :
                  (gameState.timeRemaining ?? 0) <= 5 ? 'game-page__timer-bar-fill--warning' : ''
                }`}
                style={{ 
                  animationDuration: `${timerDuration}s`
                }}
              />
            </div>
          )}
        </div>

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
          <>
            {console.log('[GamePage] Rendering GameBoard:', {
              status: gameState.status,
              themesLength: gameState.themes?.length || 0,
              themes: gameState.themes,
              canSelectQuestion
            })}
            <GameBoard
              themes={gameState.themes}
              onQuestionSelect={handleQuestionSelect}
              canSelectQuestion={canSelectQuestion}
            />
          </>
        )}

        {/* Secret Transfer Panel */}
        {gameState.status === 'secret_transfer' && isHost && (
          <SecretTransferPanel
            players={gameState.players}
            onTransfer={transferSecret}
            timeRemaining={gameState.timeRemaining}
            currentUserId={user?.id}
          />
        )}

        {/* Secret Transfer - waiting for host (non-host players) */}
        {gameState.status === 'secret_transfer' && !isHost && (
          <div className="game-page__special-waiting">
            <span className="game-page__special-waiting-icon">🐱</span>
            <h2>Кот в мешке!</h2>
            <p>Ведущий выбирает, кому передать вопрос...</p>
          </div>
        )}

        {/* Stake Betting Panel */}
        {gameState.status === 'stake_betting' && gameState.stakeInfo && (
          <StakeBettingPanel
            stakeInfo={gameState.stakeInfo}
            playerScore={currentPlayer?.score || 0}
            onPlaceStake={placeStake}
            timeRemaining={gameState.timeRemaining}
            isActivePlayer={gameState.activePlayer === user?.id}
          />
        )}

        {/* ForAll Answer Input */}
        {gameState.status === 'for_all_answering' && (
          <ForAllAnswerInput
            onSubmit={handleSubmitForAllAnswer}
            timeRemaining={gameState.timeRemaining}
            hasSubmitted={hasSubmittedForAll}
            isHost={isHost}
          />
        )}

        {/* ForAll Results */}
        {gameState.status === 'for_all_results' && forAllResults && (
          <ForAllResults
            results={forAllResults.results}
            correctAnswer={forAllResults.correct_answer}
          />
        )}

        {/* Question View - answer always visible for host */}
        {gameState.currentQuestion && (
          <QuestionView
            question={gameState.currentQuestion}
            canPressButton={canPressButton}
            onPressButton={pressButton}
            // Only show timer during button_press phase (not during question_show reading time)
            timeRemaining={gameState.status === 'button_press' ? gameState.timeRemaining : undefined}
            isHost={isHost}
            hideAnswer={false}
            startMedia={startMedia}
          />
        )}

        {/* Judging/Waiting Panel - fixed height container to prevent layout shift */}
        <div className="game-page__action-panel">
          {/* Judging Panel (for host) */}
          {canJudgeAnswer && gameState.activePlayer && (
            <div className="game-page__judging">
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

          {/* Answering Panel (for active player) */}
          {gameState.status === 'answering' && gameState.activePlayer === user?.id && (
            <div className="game-page__answering">
              <div className="game-page__answering-icon">🎤</div>
              <p className="game-page__answering-text">Говорите ваш ответ!</p>
              {gameState.timeRemaining !== undefined && (
                <p className="game-page__answering-timer">Осталось: {gameState.timeRemaining}с</p>
              )}
            </div>
          )}

          {/* Waiting for Active Player (for other players) */}
          {gameState.status === 'answering' && gameState.activePlayer !== user?.id && (
            <div className="game-page__waiting-player">
              <div className="game-page__waiting-player-icon">🎤</div>
              <p>Игрок отвечает...</p>
            </div>
          )}

          {/* Waiting for Host (for players) */}
          {gameState.status === 'answer_judging' && !isHost && (
            <div className="game-page__waiting-host">
              <div className="game-page__waiting-host-icon">🎤</div>
              <p>Скажите ответ вслух!</p>
            </div>
          )}
        </div>

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
