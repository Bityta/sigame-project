package game

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
	"github.com/sigame/game/internal/transport/websocket"
)

// MediaLoadStatus represents a client's media loading state
type MediaLoadStatus struct {
	UserID      uuid.UUID
	Loaded      int
	Total       int
	BytesLoaded int64
	Percent     int
	Complete    bool
	UpdatedAt   time.Time
}

// MediaTracker tracks media loading progress for all clients in a game
type MediaTracker struct {
	roundNumber int
	manifest    []websocket.MediaItem
	totalSize   int64
	clients     map[uuid.UUID]*MediaLoadStatus
	mu          sync.RWMutex
}

// NewMediaTracker creates a new media tracker for a round
func NewMediaTracker(roundNumber int) *MediaTracker {
	return &MediaTracker{
		roundNumber: roundNumber,
		manifest:    make([]websocket.MediaItem, 0),
		clients:     make(map[uuid.UUID]*MediaLoadStatus),
	}
}

// BuildManifest builds the media manifest from pack round data
func (mt *MediaTracker) BuildManifest(round *domain.Round) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.manifest = make([]websocket.MediaItem, 0)
	mt.totalSize = 0

	for themeIdx, theme := range round.Themes {
		for _, question := range theme.Questions {
			// Skip text-only questions
			if question.MediaType == "" || question.MediaType == "text" || question.MediaURL == "" {
				continue
			}

			mediaID := buildMediaID(mt.roundNumber, themeIdx, question.Price)
			
			item := websocket.MediaItem{
				ID:   mediaID,
				Type: question.MediaType,
				URL:  question.MediaURL,
				Size: estimateMediaSize(question.MediaType), // Estimate since we don't have actual size
				QuestionRef: websocket.QuestionRef{
					ThemeIndex:    themeIdx,
					QuestionPrice: question.Price,
				},
			}

			mt.manifest = append(mt.manifest, item)
			mt.totalSize += item.Size
		}
	}
}

// GetManifest returns the media manifest
func (mt *MediaTracker) GetManifest() ([]websocket.MediaItem, int64) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	return mt.manifest, mt.totalSize
}

// HasMedia returns true if the round has any media to load
func (mt *MediaTracker) HasMedia() bool {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	return len(mt.manifest) > 0
}

// RegisterClient registers a client to track their loading progress
func (mt *MediaTracker) RegisterClient(userID uuid.UUID) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.clients[userID] = &MediaLoadStatus{
		UserID:    userID,
		Total:     len(mt.manifest),
		Complete:  len(mt.manifest) == 0, // Already complete if no media
		UpdatedAt: time.Now(),
	}
}

// UpdateProgress updates a client's loading progress
func (mt *MediaTracker) UpdateProgress(userID uuid.UUID, loaded, total int, bytesLoaded int64, percent int) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	status, exists := mt.clients[userID]
	if !exists {
		status = &MediaLoadStatus{UserID: userID}
		mt.clients[userID] = status
	}

	status.Loaded = loaded
	status.Total = total
	status.BytesLoaded = bytesLoaded
	status.Percent = percent
	status.UpdatedAt = time.Now()
}

// MarkComplete marks a client as having finished loading
func (mt *MediaTracker) MarkComplete(userID uuid.UUID, loadedCount int) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	status, exists := mt.clients[userID]
	if !exists {
		status = &MediaLoadStatus{UserID: userID}
		mt.clients[userID] = status
	}

	status.Loaded = loadedCount
	status.Complete = true
	status.Percent = 100
	status.UpdatedAt = time.Now()
}

// AllClientsReady returns true if all registered clients have finished loading
func (mt *MediaTracker) AllClientsReady() bool {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	if len(mt.clients) == 0 {
		return true // No clients registered yet
	}

	for _, status := range mt.clients {
		if !status.Complete {
			return false
		}
	}

	return true
}

// GetPendingClients returns list of clients still loading
func (mt *MediaTracker) GetPendingClients() []uuid.UUID {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	pending := make([]uuid.UUID, 0)
	for userID, status := range mt.clients {
		if !status.Complete {
			pending = append(pending, userID)
		}
	}
	return pending
}

// GetOverallProgress returns aggregate loading progress
func (mt *MediaTracker) GetOverallProgress() (totalPercent int, clientsReady int, totalClients int) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	totalClients = len(mt.clients)
	if totalClients == 0 {
		return 100, 0, 0
	}

	var totalPct int
	for _, status := range mt.clients {
		totalPct += status.Percent
		if status.Complete {
			clientsReady++
		}
	}

	totalPercent = totalPct / totalClients
	return
}

// FindMediaByQuestion finds media item for a specific question
func (mt *MediaTracker) FindMediaByQuestion(themeIndex, price int) *websocket.MediaItem {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	for i := range mt.manifest {
		item := &mt.manifest[i]
		if item.QuestionRef.ThemeIndex == themeIndex && item.QuestionRef.QuestionPrice == price {
			return item
		}
	}
	return nil
}

// Reset clears all tracking data for a new round
func (mt *MediaTracker) Reset(roundNumber int) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.roundNumber = roundNumber
	mt.manifest = make([]websocket.MediaItem, 0)
	mt.totalSize = 0
	mt.clients = make(map[uuid.UUID]*MediaLoadStatus)
}

// Helper functions

func buildMediaID(round, themeIndex, price int) string {
	return uuid.NewString() // Simple unique ID for now
}

func estimateMediaSize(mediaType string) int64 {
	// Rough estimates since we don't have actual file sizes
	switch mediaType {
	case "image":
		return 500_000 // 500KB
	case "audio":
		return 3_000_000 // 3MB
	case "video":
		return 10_000_000 // 10MB
	default:
		return 100_000
	}
}

