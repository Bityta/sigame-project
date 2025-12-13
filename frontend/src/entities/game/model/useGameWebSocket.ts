/**
 * Game Entity - WebSocket Hook
 * React хук для управления WebSocket соединением
 */

import { useEffect, useRef, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { GameWebSocket } from '../lib/websocket';
import type { WSMessageType, GameState, StartMediaPayload } from '@/shared/types';
import { ROUTES } from '@/shared/config';

interface UseGameWebSocketOptions {
  gameId: string;
  userId: string;
  onStateUpdate?: (state: GameState) => void;
  onError?: (error: string) => void;
}

export const useGameWebSocket = ({
  gameId,
  userId,
  onStateUpdate,
  onError,
}: UseGameWebSocketOptions) => {
  const navigate = useNavigate();
  const wsRef = useRef<GameWebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [startMedia, setStartMedia] = useState<StartMediaPayload | null>(null);
  const stateReceivedRef = useRef(false);
  
  // Храним колбэки в ref чтобы избежать пересоздания useEffect
  const onStateUpdateRef = useRef(onStateUpdate);
  const onErrorRef = useRef(onError);
  
  useEffect(() => {
    onStateUpdateRef.current = onStateUpdate;
    onErrorRef.current = onError;
  }, [onStateUpdate, onError]);

  // Инициализация WebSocket
  useEffect(() => {
    if (!gameId || !userId) return;

    const ws = new GameWebSocket(gameId, userId);
    wsRef.current = ws;
    stateReceivedRef.current = false;

    // Подключаемся
    ws.connect()
      .then(() => {
        setIsConnected(true);
        
        // Таймаут на получение первого STATE_UPDATE (5 секунд)
        // Если игра не существует, сервер закроет соединение и мы не получим state
        setTimeout(() => {
          if (!stateReceivedRef.current) {
            console.error('[GameWS] Таймаут: STATE_UPDATE не получен');
            onErrorRef.current?.('Игра не найдена или не запущена');
            ws.disconnect();
            navigate(ROUTES.LOBBY);
          }
        }, 5000);
      })
      .catch((error) => {
        console.error('Ошибка подключения WebSocket:', error);
        onErrorRef.current?.('Не удалось подключиться к игре');
        // Если подключение не удалось - редирект в лобби
        navigate(ROUTES.LOBBY);
      });

    // Подписываемся на обновления состояния
    const unsubStateUpdate = ws.on<GameState>('STATE_UPDATE', (state) => {
      stateReceivedRef.current = true;
      
      // Debug logging for answer_judging state
      if (state.status === 'answer_judging') {
        console.log('[useGameWebSocket] STATE_UPDATE answer_judging:', {
          status: state.status,
          activePlayer: state.activePlayer,
          activePlayerType: typeof state.activePlayer,
          userId: userId
        });
      }
      
      setGameState(state);
      onStateUpdateRef.current?.(state);
      
      // Clear startMedia when question changes or game state changes
      if (state.status !== 'question_show' && state.status !== 'button_press') {
        setStartMedia(null);
      }
    });

    // Подписываемся на START_MEDIA для синхронного воспроизведения
    const unsubStartMedia = ws.on<StartMediaPayload>('START_MEDIA', (payload) => {
      console.log('[useGameWebSocket] START_MEDIA received:', payload);
      setStartMedia(payload);
    });

    // Подписываемся на ошибки
    const unsubError = ws.on<{ message: string }>('ERROR', (error) => {
      onErrorRef.current?.(error.message);
      // При ошибке "Game not found" редиректим в лобби
      if (error.message.includes('not found') || error.message.includes('GAME_NOT_FOUND')) {
        navigate(ROUTES.LOBBY);
      }
    });

    // Cleanup при размонтировании
    return () => {
      unsubStateUpdate();
      unsubStartMedia();
      unsubError();
      ws.disconnect();
      setIsConnected(false);
    };
  }, [gameId, userId, navigate]); // Добавили navigate в зависимости

  // Игровые действия
  const sendReady = useCallback(() => {
    wsRef.current?.sendReady();
  }, []);

  const selectQuestion = useCallback((themeId: string, questionId: string) => {
    wsRef.current?.selectQuestion(themeId, questionId);
  }, []);

  const pressButton = useCallback(() => {
    wsRef.current?.pressButton();
  }, []);

  const submitAnswer = useCallback((answer: string) => {
    wsRef.current?.submitAnswer(answer);
  }, []);

  const judgeAnswer = useCallback((answerUserId: string, correct: boolean) => {
    wsRef.current?.judgeAnswer(answerUserId, correct);
  }, []);

  // --- Special question type actions ---

  const transferSecret = useCallback((targetUserId: string) => {
    wsRef.current?.transferSecret(targetUserId);
  }, []);

  const placeStake = useCallback((amount: number, allIn: boolean = false) => {
    wsRef.current?.placeStake(amount, allIn);
  }, []);

  const submitForAllAnswer = useCallback((answer: string) => {
    wsRef.current?.submitForAllAnswer(answer);
  }, []);

  // Подписка на кастомные события
  const subscribe = useCallback(
    <T = unknown>(type: WSMessageType, handler: (payload: T) => void) => {
      return wsRef.current?.on<T>(type, handler) || (() => {});
    },
    []
  );

  return {
    isConnected,
    gameState,
    startMedia,
    sendReady,
    selectQuestion,
    pressButton,
    submitAnswer,
    judgeAnswer,
    // Special question type actions
    transferSecret,
    placeStake,
    submitForAllAnswer,
    subscribe,
  };
};
