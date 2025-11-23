/**
 * Pack Domain Types
 */

export type MediaType = 'text' | 'image' | 'audio' | 'video';
export type PackStatus = 'processing' | 'approved' | 'rejected';

export interface PackQuestion {
  id: string;
  themeId: string;
  price: number;
  questionText: string;
  answerText: string;
  mediaType: MediaType;
  mediaUrl?: string;
}

export interface PackTheme {
  id: string;
  roundId: string;
  themeName: string;
  questions?: PackQuestion[];
}

export interface PackRound {
  id: string;
  packId: string;
  roundNumber: number;
  roundName: string;
  themes?: PackTheme[];
}

export interface Pack {
  id: string;
  name: string;
  author: string;
  description: string;
  thumbnailUrl?: string;
  uploadedBy: string;
  downloadsCount: number;
  rating: number;
  ratingCount: number;
  status: PackStatus;
  createdAt: string;
  updatedAt: string;
  rounds?: PackRound[];
  tags?: string[];
}

export interface PackListQuery {
  search?: string;
  author?: string;
  tags?: string[];
  page?: number;
  size?: number;
  sort?: 'rating' | 'downloads' | 'created_at';
}

