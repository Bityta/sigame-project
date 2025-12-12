package media

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
	"github.com/sigame/game/internal/transport/websocket"
)

type MediaLoadStatus struct {
	UserID      uuid.UUID
	Loaded      int
	Total       int
	BytesLoaded int64
	Percent     int
	Complete    bool
	UpdatedAt   time.Time
}

type MediaTracker struct {
	roundNumber int
	manifest    []websocket.MediaItem
	totalSize   int64
	clients     map[uuid.UUID]*MediaLoadStatus
	mu          sync.RWMutex
}

func NewMediaTracker(roundNumber int) *MediaTracker {
	return &MediaTracker{
		roundNumber: roundNumber,
		manifest:    make([]websocket.MediaItem, 0),
		clients:     make(map[uuid.UUID]*MediaLoadStatus),
	}
}

func (mt *MediaTracker) BuildManifest(round *domain.Round) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.manifest = make([]websocket.MediaItem, 0)
	mt.totalSize = 0

	for themeIdx, theme := range round.Themes {
		for _, question := range theme.Questions {
			if question.MediaType == "" || question.MediaType == MediaTypeText || question.MediaURL == "" {
				continue
			}

			mediaID := buildMediaID(mt.roundNumber, themeIdx, question.Price)

			item := websocket.MediaItem{
				ID:   mediaID,
				Type: question.MediaType,
				URL:  question.MediaURL,
				Size: estimateMediaSize(question.MediaType),
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

func (mt *MediaTracker) GetManifest() ([]websocket.MediaItem, int64) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	return mt.manifest, mt.totalSize
}

func (mt *MediaTracker) HasMedia() bool {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	return len(mt.manifest) > 0
}

func (mt *MediaTracker) RegisterClient(userID uuid.UUID) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.clients[userID] = &MediaLoadStatus{
		UserID:    userID,
		Total:     len(mt.manifest),
		Complete:  len(mt.manifest) == 0,
		UpdatedAt: time.Now(),
	}
}

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
	status.Percent = PercentComplete
	status.UpdatedAt = time.Now()
}

func (mt *MediaTracker) AllClientsReady() bool {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	if len(mt.clients) == 0 {
		return true
	}

	for _, status := range mt.clients {
		if !status.Complete {
			return false
		}
	}

	return true
}

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

func (mt *MediaTracker) GetOverallProgress() (totalPercent int, clientsReady int, totalClients int) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	totalClients = len(mt.clients)
	if totalClients == 0 {
		return PercentComplete, 0, 0
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

func (mt *MediaTracker) Reset(roundNumber int) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.roundNumber = roundNumber
	mt.manifest = make([]websocket.MediaItem, 0)
	mt.totalSize = 0
	mt.clients = make(map[uuid.UUID]*MediaLoadStatus)
}

