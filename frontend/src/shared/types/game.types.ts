/**
 * Game Session Domain Types
 */

export type GameStatus = 
  | 'waiting'
  | 'rounds_overview'
  | 'round_start'
  | 'question_select'
  | 'question_show'
  | 'button_press'
  | 'answering'
  | 'answer_judging'
  | 'round_end'
  | 'game_end'
  | 'finished';

export type WSMessageType = 
  // Client -> Server
  | 'READY'
  | 'SELECT_QUESTION'
  | 'PRESS_BUTTON'
  | 'SUBMIT_ANSWER'
  | 'JUDGE_ANSWER'
  // Server -> Client
  | 'STATE_UPDATE'
  | 'QUESTION_SELECTED'
  | 'BUTTON_PRESSED'
  | 'ANSWER_RESULT'
  | 'ROUND_COMPLETE'
  | 'GAME_COMPLETE'
  | 'ERROR';

export interface WSMessage<T = any> {
  type: WSMessageType;
  payload?: T;
}

export interface RoundOverview {
  roundNumber: number;
  name: string;
  themeNames: string[];
}

export interface PlayerScore {
  userId: string;
  username: string;
  score: number;
  rank: number;
}

export interface GameState {
  gameId: string;
  status: GameStatus;
  currentRound: number;
  roundName?: string;
  themes?: ThemeState[];
  players: PlayerState[];
  activePlayer?: string;
  currentQuestion?: QuestionState;
  timeRemaining?: number;
  message?: string;
  allRounds?: RoundOverview[];
  winners?: PlayerScore[];
  finalScores?: PlayerScore[];
}

export interface PlayerState {
  userId: string;
  username: string;
  role: 'host' | 'player';
  score: number;
  isActive: boolean;
  isReady: boolean;
}

export interface ThemeState {
  name: string;
  questions: QuestionState[];
}

export interface QuestionState {
  id: string;
  price: number;
  available: boolean;
  text?: string;
  mediaType?: string;
  mediaUrl?: string; // URL to media file (image/audio/video)
  mediaDurationMs?: number; // Duration in ms (for audio/video)
  answer?: string; // Only for host
}

export interface ButtonPressedPayload {
  userId: string;
  username: string;
  latencyMs: number;
}

export interface AnswerResultPayload {
  userId: string;
  username: string;
  correct: boolean;
  answer?: string;
  score: number;
  scoreDelta: number;
}

export interface RoundCompletePayload {
  roundNumber: number;
  scores: Array<{
    userId: string;
    username: string;
    score: number;
    rank: number;
  }>;
  nextRound?: number;
}

export interface GameCompletePayload {
  winners: Array<{
    userId: string;
    username: string;
    score: number;
    rank: number;
  }>;
  scores: Array<{
    userId: string;
    username: string;
    score: number;
    rank: number;
  }>;
}

export interface ErrorPayload {
  message: string;
  code?: string;
}

export interface CreateGameResponse {
  gameId: string;
  websocketUrl: string;
  status: string;
}

