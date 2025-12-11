/**
 * SyncMediaPlayer - Synchronized media playback component
 * Starts media playback at a specified timestamp for synchronized viewing across clients
 */

import { useEffect, useRef, useState, useCallback } from 'react';
import type { StartMediaPayload } from '@/shared/types';
import { mediaCache } from '@/entities/game/lib/mediaCache';
import './SyncMediaPlayer.css';

interface SyncMediaPlayerProps {
  startMedia: StartMediaPayload | null;
  fallbackUrl?: string;  // Direct URL if not in cache
  onEnded?: () => void;
  onError?: (error: Error) => void;
  autoPlay?: boolean;
  muted?: boolean;
  controls?: boolean;
  className?: string;
}

export const SyncMediaPlayer: React.FC<SyncMediaPlayerProps> = ({
  startMedia,
  fallbackUrl,
  onEnded,
  onError,
  autoPlay = true,
  muted = false,
  controls = true,
  className = '',
}) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const audioRef = useRef<HTMLAudioElement>(null);
  const [isReady, setIsReady] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const startTimeoutRef = useRef<number | null>(null);

  // Get media source - prefer cached, fallback to direct URL
  const getMediaSource = useCallback((): string => {
    if (!startMedia) return fallbackUrl || '';

    // Try to get from cache first
    const cachedUrl = mediaCache.getObjectUrl(startMedia.media_id);
    if (cachedUrl) {
      console.log('[SyncMediaPlayer] Using cached media:', startMedia.media_id);
      return cachedUrl;
    }

    // Fallback to direct URL
    console.log('[SyncMediaPlayer] Using direct URL:', startMedia.url);
    return startMedia.url || fallbackUrl || '';
  }, [startMedia, fallbackUrl]);

  // Schedule playback at the specified time
  const schedulePlayback = useCallback(() => {
    if (!startMedia || !autoPlay) return;

    const element = startMedia.media_type === 'video' ? videoRef.current : audioRef.current;
    if (!element) return;

    const now = Date.now();
    const startAt = startMedia.start_at;
    const delay = startAt - now;

    console.log(`[SyncMediaPlayer] Scheduling playback: start_at=${startAt}, now=${now}, delay=${delay}ms`);

    if (delay > 0) {
      // Schedule for future
      startTimeoutRef.current = window.setTimeout(() => {
        console.log('[SyncMediaPlayer] Starting playback now');
        element.play().catch((e) => {
          console.error('[SyncMediaPlayer] Play failed:', e);
          setError('Failed to start playback');
          onError?.(e);
        });
      }, delay);
    } else if (delay > -startMedia.duration_ms) {
      // Already should have started, seek to correct position
      const seekTo = Math.abs(delay) / 1000;
      console.log(`[SyncMediaPlayer] Late start, seeking to ${seekTo}s`);
      element.currentTime = seekTo;
      element.play().catch((e) => {
        console.error('[SyncMediaPlayer] Play failed:', e);
        setError('Failed to start playback');
        onError?.(e);
      });
    } else {
      // Media already ended
      console.log('[SyncMediaPlayer] Media already ended');
      onEnded?.();
    }
  }, [startMedia, autoPlay, onEnded, onError]);

  // Setup when startMedia changes
  useEffect(() => {
    if (!startMedia) return;

    setError(null);
    setIsReady(false);

    // Clear any pending timeout
    if (startTimeoutRef.current) {
      clearTimeout(startTimeoutRef.current);
    }

    return () => {
      if (startTimeoutRef.current) {
        clearTimeout(startTimeoutRef.current);
      }
    };
  }, [startMedia]);

  // Handle media loaded
  const handleCanPlay = useCallback(() => {
    console.log('[SyncMediaPlayer] Media can play');
    setIsReady(true);
    schedulePlayback();
  }, [schedulePlayback]);

  // Handle media ended
  const handleEnded = useCallback(() => {
    console.log('[SyncMediaPlayer] Media ended');
    onEnded?.();
  }, [onEnded]);

  // Handle media error
  const handleError = useCallback((e: React.SyntheticEvent<HTMLMediaElement>) => {
    const mediaError = (e.target as HTMLMediaElement).error;
    const errorMsg = mediaError?.message || 'Unknown media error';
    console.error('[SyncMediaPlayer] Media error:', errorMsg);
    setError(errorMsg);
    onError?.(new Error(errorMsg));
  }, [onError]);

  const mediaSource = getMediaSource();

  if (!startMedia && !fallbackUrl) {
    return null;
  }

  const mediaType = startMedia?.media_type || 'video';

  if (error) {
    return (
      <div className={`sync-media-player sync-media-player--error ${className}`}>
        <div className="sync-media-player__error">
          <span className="sync-media-player__error-icon">⚠️</span>
          <span className="sync-media-player__error-text">{error}</span>
        </div>
      </div>
    );
  }

  if (mediaType === 'image') {
    return (
      <div className={`sync-media-player sync-media-player--image ${className}`}>
        <img
          src={mediaSource}
          alt="Question media"
          className="sync-media-player__image"
          onError={() => setError('Failed to load image')}
        />
      </div>
    );
  }

  if (mediaType === 'audio') {
    return (
      <div className={`sync-media-player sync-media-player--audio ${className}`}>
        <div className="sync-media-player__audio-visual">
          <div className="sync-media-player__audio-icon">🎵</div>
          <div className="sync-media-player__audio-waves">
            <span></span><span></span><span></span><span></span><span></span>
          </div>
        </div>
        <audio
          ref={audioRef}
          src={mediaSource}
          muted={muted}
          controls={controls}
          onCanPlay={handleCanPlay}
          onEnded={handleEnded}
          onError={handleError}
          className="sync-media-player__audio"
        />
      </div>
    );
  }

  // Default: video
  return (
    <div className={`sync-media-player sync-media-player--video ${className}`}>
      {!isReady && (
        <div className="sync-media-player__loading">
          <div className="sync-media-player__spinner"></div>
          <span>Загрузка видео...</span>
        </div>
      )}
      <video
        ref={videoRef}
        src={mediaSource}
        muted={muted}
        controls={controls}
        playsInline
        onCanPlay={handleCanPlay}
        onEnded={handleEnded}
        onError={handleError}
        className="sync-media-player__video"
      />
    </div>
  );
};

export default SyncMediaPlayer;

