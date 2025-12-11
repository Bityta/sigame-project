/**
 * Game Entity - WebSocket Service
 * Управление WebSocket соединением для игровой сессии
 */

import { API_CONFIG } from '@/shared/config';
import type {
  WSMessage,
  WSMessageType,
  GameState,
  RoundMediaManifestPayload,
  StartMediaPayload,
} from '@/shared/types';
import { mediaCache, type MediaLoadProgress } from './mediaCache';

type MessageHandler<T = unknown> = (payload: T) => void;

export class GameWebSocket {
  private ws: WebSocket | null = null;
  private handlers: Map<WSMessageType, MessageHandler[]> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private gameId: string;
  private userId: string;
  private progressIntervalId: number | null = null;

  constructor(gameId: string, userId: string) {
    this.gameId = gameId;
    this.userId = userId;
    this.setupMediaCallbacks();
  }

  /**
   * Setup media cache callbacks
   */
  private setupMediaCallbacks(): void {
    // Progress callback - send updates to server every 500ms
    let lastReportTime = 0;
    mediaCache.setProgressCallback((progress: MediaLoadProgress) => {
      const now = Date.now();
      if (now - lastReportTime >= 500) {
        this.sendMediaLoadProgress(progress);
        lastReportTime = now;
      }
    });

    // Complete callback - notify server when done
    mediaCache.setCompleteCallback((round: number, loadedCount: number) => {
      this.sendMediaLoadComplete(round, loadedCount);
    });
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
          // Backend может отправить несколько JSON разделённых \n
          const messages = event.data.split('\n').filter((s: string) => s.trim());
          
          for (const msgStr of messages) {
            try {
              const message: WSMessage = JSON.parse(msgStr);
              console.log('[GameWS] Получено:', message.type, message.payload);
              this.handleMessage(message);
            } catch (error) {
              console.error('[GameWS] Ошибка парсинга сообщения:', error);
            }
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
    // Handle media manifest - start preloading
    if (message.type === 'ROUND_MEDIA_MANIFEST') {
      this.handleMediaManifest(message.payload as RoundMediaManifestPayload);
    }

    const handlers = this.handlers.get(message.type);
    if (handlers) {
      handlers.forEach((handler) => handler(message.payload));
    }
  }

  /**
   * Handle ROUND_MEDIA_MANIFEST - start preloading media
   */
  private handleMediaManifest(payload: RoundMediaManifestPayload): void {
    console.log(`[GameWS] Received media manifest for round ${payload.round}: ${payload.total_count} files`);
    
    // Start preloading in background
    mediaCache.preloadRound(payload).catch((error) => {
      console.error('[GameWS] Media preload error:', error);
    });
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
    this.sendGameMessage('READY');
  }

  /**
   * Выбрать вопрос
   */
  selectQuestion(themeId: string, questionId: string): void {
    this.sendGameMessage('SELECT_QUESTION', {
      theme_id: themeId,
      question_id: questionId,
    });
  }

  /**
   * Нажать на кнопку
   */
  pressButton(): void {
    this.sendGameMessage('PRESS_BUTTON');
  }

  /**
   * Отправить ответ
   */
  submitAnswer(answer: string): void {
    this.sendGameMessage('SUBMIT_ANSWER', { answer });
  }

  /**
   * Оценить ответ (только для хоста)
   */
  judgeAnswer(answerUserId: string, correct: boolean): void {
    this.sendGameMessage('JUDGE_ANSWER', {
      user_id: answerUserId,
      correct,
    });
  }

  /**
   * Отправить игровое сообщение с user_id и game_id на верхнем уровне
   */
  private sendGameMessage(type: WSMessageType, payload?: Record<string, unknown>): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('[GameWS] Соединение не открыто');
      return;
    }

    const message = {
      type,
      user_id: this.userId,
      game_id: this.gameId,
      ...(payload && { payload }),
    };
    console.log('[GameWS] Отправка:', type, message);
    this.ws.send(JSON.stringify(message));
  }

  // --- Media sync methods ---

  /**
   * Send media loading progress to server
   */
  private sendMediaLoadProgress(progress: MediaLoadProgress): void {
    this.sendGameMessage('MEDIA_LOAD_PROGRESS', {
      loaded: progress.loaded,
      total: progress.total,
      bytes_loaded: progress.bytesLoaded,
      percent: progress.percent,
    });
  }

  /**
   * Send media loading complete notification
   */
  private sendMediaLoadComplete(round: number, loadedCount: number): void {
    this.sendGameMessage('MEDIA_LOAD_COMPLETE', {
      round,
      loaded_count: loadedCount,
    });
  }

  /**
   * Get media cache for accessing preloaded media
   */
  getMediaCache() {
    return mediaCache;
  }
}

