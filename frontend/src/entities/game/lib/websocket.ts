/**
 * Game Entity - WebSocket Service
 * Управление WebSocket соединением для игровой сессии
 */

import { API_CONFIG } from '@/shared/config';
import type {
  WSMessage,
  WSMessageType,
  GameState,
} from '@/shared/types';

type MessageHandler<T = unknown> = (payload: T) => void;

export class GameWebSocket {
  private ws: WebSocket | null = null;
  private handlers: Map<WSMessageType, MessageHandler[]> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private gameId: string;
  private userId: string;

  constructor(gameId: string, userId: string) {
    this.gameId = gameId;
    this.userId = userId;
  }

  /**
   * Подключиться к WebSocket
   */
  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        // Формируем WebSocket URL
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsHost = API_CONFIG.GAME_BASE_URL.replace('http://', '').replace('https://', '');
        const url = `${wsProtocol}//${wsHost}/api/game/${this.gameId}/ws?user_id=${this.userId}`;
        
        console.log('[GameWS] Подключение к:', url);
        this.ws = new WebSocket(url);

        this.ws.onopen = () => {
          console.log('[GameWS] Соединение установлено');
          this.reconnectAttempts = 0;
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message: WSMessage = JSON.parse(event.data);
            console.log('[GameWS] Получено:', message.type, message.payload);
            this.handleMessage(message);
          } catch (error) {
            console.error('[GameWS] Ошибка парсинга сообщения:', error);
          }
        };

        this.ws.onerror = (error) => {
          console.error('[GameWS] Ошибка WebSocket:', error);
          reject(error);
        };

        this.ws.onclose = (event) => {
          console.log('[GameWS] Соединение закрыто:', event.code, event.reason);
          if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.attemptReconnect();
          }
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  /**
   * Отключиться от WebSocket
   */
  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.handlers.clear();
  }

  /**
   * Отправить сообщение на сервер
   */
  send<T = unknown>(type: WSMessageType, payload?: T): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('[GameWS] Соединение не открыто');
      return;
    }

    const message: WSMessage<T> = { type, payload };
    console.log('[GameWS] Отправка:', message.type, message.payload);
    this.ws.send(JSON.stringify(message));
  }

  /**
   * Подписаться на тип сообщения
   */
  on<T = unknown>(type: WSMessageType, handler: MessageHandler<T>): () => void {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, []);
    }
    this.handlers.get(type)!.push(handler as MessageHandler);

    // Возвращаем функцию отписки
    return () => this.off(type, handler);
  }

  /**
   * Отписаться от типа сообщения
   */
  off<T = unknown>(type: WSMessageType, handler: MessageHandler<T>): void {
    const handlers = this.handlers.get(type);
    if (handlers) {
      const index = handlers.indexOf(handler as MessageHandler);
      if (index !== -1) {
        handlers.splice(index, 1);
      }
    }
  }

  /**
   * Обработать входящее сообщение
   */
  private handleMessage(message: WSMessage): void {
    const handlers = this.handlers.get(message.type);
    if (handlers) {
      handlers.forEach((handler) => handler(message.payload));
    }
  }

  /**
   * Попытка переподключения
   */
  private attemptReconnect(): void {
    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

    console.log(`[GameWS] Переподключение через ${delay}ms (попытка ${this.reconnectAttempts})`);

    setTimeout(() => {
      console.log('[GameWS] Попытка переподключения...');
      this.connect().catch((error) => {
        console.error('[GameWS] Не удалось переподключиться:', error);
      });
    }, delay);
  }

  /**
   * Проверка статуса соединения
   */
  get isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }

  // --- Игровые действия ---

  /**
   * Отправить готовность
   */
  sendReady(): void {
    this.send('READY', {
      user_id: this.userId,
      game_id: this.gameId,
    });
  }

  /**
   * Выбрать вопрос
   */
  selectQuestion(themeId: string, questionId: string): void {
    this.send('SELECT_QUESTION', {
      user_id: this.userId,
      game_id: this.gameId,
      payload: {
        theme_id: themeId,
        question_id: questionId,
      },
    });
  }

  /**
   * Нажать на кнопку
   */
  pressButton(): void {
    this.send('PRESS_BUTTON', {
      user_id: this.userId,
      game_id: this.gameId,
    });
  }

  /**
   * Отправить ответ
   */
  submitAnswer(answer: string): void {
    this.send('SUBMIT_ANSWER', {
      user_id: this.userId,
      game_id: this.gameId,
      payload: {
        answer,
      },
    });
  }

  /**
   * Оценить ответ (только для хоста)
   */
  judgeAnswer(answerUserId: string, correct: boolean): void {
    this.send('JUDGE_ANSWER', {
      user_id: this.userId,
      game_id: this.gameId,
      payload: {
        user_id: answerUserId,
        correct,
      },
    });
  }
}

