/**
 * Game Entity - WebSocket Hook
 * React хук для управления WebSocket соединением
 */

import { useEffect, useRef, useState, useCallback } from 'react';
import { GameWebSocket } from '../lib/websocket';
import type { WSMessageType, GameState } from '@/shared/types';

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
  const wsRef = useRef<GameWebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [gameState, setGameState] = useState<GameState | null>(null);

  // Инициализация WebSocket
  useEffect(() => {
    if (!gameId || !userId) return;

    const ws = new GameWebSocket(gameId, userId);
    wsRef.current = ws;

    // Подключаемся
    ws.connect()
      .then(() => {
        setIsConnected(true);
      })
      .catch((error) => {
        console.error('Ошибка подключения WebSocket:', error);
        onError?.('Не удалось подключиться к игре');
      });

    // Подписываемся на обновления состояния
    const unsubStateUpdate = ws.on<GameState>('STATE_UPDATE', (state) => {
      setGameState(state);
      onStateUpdate?.(state);
    });

    // Подписываемся на ошибки
    const unsubError = ws.on<{ message: string }>('ERROR', (error) => {
      onError?.(error.message);
    });

    // Cleanup при размонтировании
    return () => {
      unsubStateUpdate();
      unsubError();
      ws.disconnect();
      setIsConnected(false);
    };
  }, [gameId, userId, onStateUpdate, onError]);

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
    sendReady,
    selectQuestion,
    pressButton,
    submitAnswer,
    judgeAnswer,
    subscribe,
  };
};

