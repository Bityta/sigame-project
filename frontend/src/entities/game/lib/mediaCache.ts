/**
 * Media Cache - Preloading and caching media files for synchronized playback
 */

import type { MediaItem, RoundMediaManifestPayload } from '@/shared/types';

export interface MediaLoadProgress {
  loaded: number;
  total: number;
  bytesLoaded: number;
  percent: number;
}

type ProgressCallback = (progress: MediaLoadProgress) => void;
type CompleteCallback = (round: number, loadedCount: number) => void;

/**
 * MediaCache handles preloading and caching of media files
 */
export class MediaCache {
  private cache: Map<string, Blob> = new Map();
  private objectUrls: Map<string, string> = new Map();
  private currentRound = 0;
  private manifest: MediaItem[] = [];
  private loadedCount = 0;
  private totalBytes = 0;
  private loadedBytes = 0;
  private isLoading = false;
  private abortController: AbortController | null = null;

  private onProgress: ProgressCallback | null = null;
  private onComplete: CompleteCallback | null = null;

  /**
   * Set progress callback
   */
  setProgressCallback(callback: ProgressCallback | null): void {
    this.onProgress = callback;
  }

  /**
   * Set completion callback
   */
  setCompleteCallback(callback: CompleteCallback | null): void {
    this.onComplete = callback;
  }

  /**
   * Start preloading media from manifest
   */
  async preloadRound(payload: RoundMediaManifestPayload): Promise<void> {
    // Cancel any ongoing preload
    this.cancelPreload();

    this.currentRound = payload.round;
    this.manifest = payload.media;
    this.totalBytes = payload.total_size;
    this.loadedCount = 0;
    this.loadedBytes = 0;
    this.isLoading = true;
    this.abortController = new AbortController();

    console.log(`[MediaCache] Starting preload for round ${payload.round}: ${payload.total_count} files, ${payload.total_size} bytes`);

    if (payload.media.length === 0) {
      this.isLoading = false;
      this.onComplete?.(payload.round, 0);
      return;
    }

    // Load all media files in parallel (with concurrency limit)
    const concurrencyLimit = 3;
    const queue = [...payload.media];
    const inProgress: Promise<void>[] = [];

    while (queue.length > 0 || inProgress.length > 0) {
      // Fill up to concurrency limit
      while (queue.length > 0 && inProgress.length < concurrencyLimit) {
        const item = queue.shift()!;
        const promise = this.loadMediaItem(item).then(() => {
          const index = inProgress.indexOf(promise);
          if (index > -1) inProgress.splice(index, 1);
        });
        inProgress.push(promise);
      }

      // Wait for at least one to complete
      if (inProgress.length > 0) {
        await Promise.race(inProgress);
      }
    }

    this.isLoading = false;
    console.log(`[MediaCache] Preload complete: ${this.loadedCount} files loaded`);
    this.onComplete?.(this.currentRound, this.loadedCount);
  }

  /**
   * Load a single media item
   */
  private async loadMediaItem(item: MediaItem): Promise<void> {
    try {
      const response = await fetch(item.url, {
        signal: this.abortController?.signal,
      });

      if (!response.ok) {
        console.error(`[MediaCache] Failed to load ${item.id}: ${response.status}`);
        return;
      }

      const blob = await response.blob();
      this.cache.set(item.id, blob);

      // Create object URL for quick access
      const objectUrl = URL.createObjectURL(blob);
      this.objectUrls.set(item.id, objectUrl);

      this.loadedCount++;
      this.loadedBytes += item.size;

      this.reportProgress();

      console.log(`[MediaCache] Loaded ${item.id} (${item.type}): ${blob.size} bytes`);
    } catch (error) {
      if ((error as Error).name === 'AbortError') {
        console.log(`[MediaCache] Load aborted for ${item.id}`);
      } else {
        console.error(`[MediaCache] Error loading ${item.id}:`, error);
      }
    }
  }

  /**
   * Report current progress
   */
  private reportProgress(): void {
    const progress: MediaLoadProgress = {
      loaded: this.loadedCount,
      total: this.manifest.length,
      bytesLoaded: this.loadedBytes,
      percent: this.manifest.length > 0 
        ? Math.round((this.loadedCount / this.manifest.length) * 100)
        : 100,
    };

    this.onProgress?.(progress);
  }

  /**
   * Cancel ongoing preload
   */
  cancelPreload(): void {
    if (this.abortController) {
      this.abortController.abort();
      this.abortController = null;
    }
    this.isLoading = false;
  }

  /**
   * Get cached media blob by ID
   */
  getBlob(mediaId: string): Blob | undefined {
    return this.cache.get(mediaId);
  }

  /**
   * Get cached media object URL by ID (for use in src attributes)
   */
  getObjectUrl(mediaId: string): string | undefined {
    return this.objectUrls.get(mediaId);
  }

  /**
   * Get media by question reference
   */
  getMediaByQuestion(themeIndex: number, price: number): { blob?: Blob; url?: string } | undefined {
    const item = this.manifest.find(
      m => m.question_ref.theme === themeIndex && m.question_ref.price === price
    );

    if (!item) return undefined;

    return {
      blob: this.cache.get(item.id),
      url: this.objectUrls.get(item.id),
    };
  }

  /**
   * Check if media is cached
   */
  isCached(mediaId: string): boolean {
    return this.cache.has(mediaId);
  }

  /**
   * Get current loading progress
   */
  getProgress(): MediaLoadProgress {
    return {
      loaded: this.loadedCount,
      total: this.manifest.length,
      bytesLoaded: this.loadedBytes,
      percent: this.manifest.length > 0 
        ? Math.round((this.loadedCount / this.manifest.length) * 100)
        : 100,
    };
  }

  /**
   * Check if all media is loaded
   */
  isComplete(): boolean {
    return !this.isLoading && this.loadedCount === this.manifest.length;
  }

  /**
   * Clear cache and revoke object URLs
   */
  clear(): void {
    this.cancelPreload();

    // Revoke all object URLs to free memory
    for (const url of this.objectUrls.values()) {
      URL.revokeObjectURL(url);
    }

    this.cache.clear();
    this.objectUrls.clear();
    this.manifest = [];
    this.loadedCount = 0;
    this.loadedBytes = 0;
    this.totalBytes = 0;
  }
}

// Singleton instance
export const mediaCache = new MediaCache();

