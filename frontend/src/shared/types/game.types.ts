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
  | 'finished'
  // Special question type statuses
  | 'secret_transfer'    // Host chooses player for secret question
  | 'stake_betting'      // Player places a stake
  | 'for_all_answering'  // All players answering
  | 'for_all_results';   // Showing forAll results

export type QuestionType = 'normal' | 'secret' | 'stake' | 'forAll';

export type WSMessageType = 
  // Client -> Server
  | 'READY'
  | 'SELECT_QUESTION'
  | 'PRESS_BUTTON'
  | 'SUBMIT_ANSWER'
  | 'JUDGE_ANSWER'
  | 'PONG'           // Response to PING for RTT measurement
  | 'MEDIA_LOAD_PROGRESS'
  | 'MEDIA_LOAD_COMPLETE'
  | 'TRANSFER_SECRET'      // Host transfers secret question
  | 'PLACE_STAKE'          // Player places stake
  | 'SUBMIT_FOR_ALL_ANSWER' // Player submits forAll answer
  // Server -> Client
  | 'STATE_UPDATE'
  | 'QUESTION_SELECTED'
  | 'BUTTON_PRESSED'
  | 'ANSWER_RESULT'
  | 'ROUND_COMPLETE'
  | 'GAME_COMPLETE'
  | 'ERROR'
  | 'PING'           // RTT measurement ping
  | 'ROUND_MEDIA_MANIFEST'
  | 'START_MEDIA'
  | 'SECRET_TRANSFERRED'   // Secret question transferred notification
  | 'STAKE_PLACED'         // Stake placed notification
  | 'FOR_ALL_RESULTS';     // ForAll results

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
  // Special question type fields
  stakeInfo?: StakeInfo;
  secretTarget?: string; // Target player ID for secret question
  forAllResults?: ForAllAnswerResult[];
}

// Stake betting info
export interface StakeInfo {
  minBet: number;
  maxBet: number;
  currentBet: number;
  isAllIn: boolean;
}

// ForAll question result for a single player
export interface ForAllAnswerResult {
  userId: string;
  username: string;
  answer: string;
  isCorrect: boolean;
  scoreDelta: number;
}

export interface PlayerState {
  userId: string;
  username: string;
  avatarUrl?: string;
  role: 'host' | 'player';
  score: number;
  isActive: boolean;
  isReady: boolean;
  isConnected: boolean;
}

export interface ThemeState {
  name: string;
  questions: QuestionState[];
}

export interface QuestionState {
  id: string;
  price: number;
  available: boolean;
  type: QuestionType; // normal, secret, stake, forAll
  text?: string;
  mediaType?: string;
  mediaUrl?: string; // URL to media file (image/audio/video)
  mediaDurationMs?: number; // Duration in ms (for audio/video)
  answer?: string; // Only for host
}

// PING/PONG for RTT measurement
export interface PingPayload {
  server_time: number;  // Server timestamp in milliseconds
}

export interface PongPayload {
  server_time: number;  // Echo back from PING
  client_time: number;  // Client's current timestamp
}

// Button press info for a single player
export interface PressInfo {
  user_id: string;
  username: string;
  time_ms: number;  // Adjusted reaction time in milliseconds
}

export interface ButtonPressedPayload {
  winner_id: string;
  winner_name: string;
  reaction_time_ms: number;  // Winner's adjusted reaction time
  all_presses: PressInfo[];  // All button presses sorted by adjusted time
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

// Media sync types
export interface QuestionRef {
  theme: number;
  price: number;
}

export interface MediaItem {
  id: string;
  type: 'image' | 'audio' | 'video';
  url: string;
  size: number;
  question_ref: QuestionRef;
}

export interface RoundMediaManifestPayload {
  round: number;
  media: MediaItem[];
  total_size: number;
  total_count: number;
}

export interface MediaLoadProgressPayload {
  loaded: number;
  total: number;
  bytes_loaded: number;
  percent: number;
}

export interface MediaLoadCompletePayload {
  round: number;
  loaded_count: number;
}

export interface StartMediaPayload {
  media_id: string;
  media_type: 'image' | 'audio' | 'video';
  url: string;
  start_at: number;  // Unix timestamp in ms
  duration_ms: number;
}

// Special question type payloads

// Client -> Server: Transfer secret question to player
export interface TransferSecretPayload {
  target_user_id: string;
}

// Client -> Server: Place stake
export interface PlaceStakePayload {
  amount: number;
  all_in: boolean;
}

// Client -> Server: Submit forAll answer
export interface SubmitForAllAnswerPayload {
  answer: string;
}

// Server -> Client: Secret transferred notification
export interface SecretTransferredPayload {
  from_user_id: string;
  from_username: string;
  to_user_id: string;
  to_username: string;
}

// Server -> Client: Stake placed notification
export interface StakePlacedPayload {
  user_id: string;
  username: string;
  amount: number;
  all_in: boolean;
}

// Server -> Client: ForAll results
export interface ForAllResultsPayload {
  correct_answer: string;
  results: ForAllAnswerResult[];
}
