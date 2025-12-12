package media

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
)

func TestNewMediaTracker(t *testing.T) {
	roundNumber := 1
	tracker := NewMediaTracker(roundNumber)

	if tracker == nil {
		t.Fatal("NewMediaTracker returned nil")
	}

	if tracker.roundNumber != roundNumber {
		t.Errorf("Expected roundNumber %d, got %d", roundNumber, tracker.roundNumber)
	}

	if tracker.manifest == nil {
		t.Error("manifest should not be nil")
	}

	if len(tracker.manifest) != 0 {
		t.Errorf("Expected empty manifest, got %d items", len(tracker.manifest))
	}

	if tracker.clients == nil {
		t.Error("clients map should not be nil")
	}

	if len(tracker.clients) != 0 {
		t.Errorf("Expected empty clients map, got %d clients", len(tracker.clients))
	}

	if tracker.totalSize != 0 {
		t.Errorf("Expected totalSize 0, got %d", tracker.totalSize)
	}
}

func TestMediaTracker_BuildManifest(t *testing.T) {
	tracker := NewMediaTracker(1)

	round := &domain.Round{
		ID:          "round1",
		RoundNumber: 1,
		Name:        "Round 1",
		Themes: []*domain.Theme{
			{
				ID:   "theme1",
				Name: "Theme 1",
				Questions: []*domain.Question{
					{
						ID:        "q1",
						Price:      100,
						MediaType:  MediaTypeImage,
						MediaURL:   "http:
					},
					{
						ID:        "q2",
						Price:      200,
						MediaType:  MediaTypeAudio,
						MediaURL:   "http:
					},
				},
			},
			{
				ID:   "theme2",
				Name: "Theme 2",
				Questions: []*domain.Question{
					{
						ID:        "q3",
						Price:      300,
						MediaType:  MediaTypeVideo,
						MediaURL:   "http:
					},
				},
			},
		},
	}

	tracker.BuildManifest(round)

	manifest, totalSize := tracker.GetManifest()

	if len(manifest) != 3 {
		t.Errorf("Expected 3 media items, got %d", len(manifest))
	}

	expectedSize := int64(MediaSizeImage + MediaSizeAudio + MediaSizeVideo)
	if totalSize != expectedSize {
		t.Errorf("Expected totalSize %d, got %d", expectedSize, totalSize)
	}

	if !tracker.HasMedia() {
		t.Error("HasMedia should return true")
	}
}

func TestMediaTracker_BuildManifest_SkipsTextOnly(t *testing.T) {
	tracker := NewMediaTracker(1)

	round := &domain.Round{
		ID:          "round1",
		RoundNumber: 1,
		Name:        "Round 1",
		Themes: []*domain.Theme{
			{
				ID:   "theme1",
				Name: "Theme 1",
				Questions: []*domain.Question{
					{
						ID:        "q1",
						Price:      100,
						MediaType:  MediaTypeText,
						MediaURL:   "",
					},
					{
						ID:        "q2",
						Price:      200,
						MediaType:  "",
						MediaURL:   "",
					},
					{
						ID:        "q3",
						Price:      300,
						MediaType:  MediaTypeImage,
						MediaURL:   "",
					},
					{
						ID:        "q4",
						Price:      400,
						MediaType:  MediaTypeImage,
						MediaURL:   "http:
					},
				},
			},
		},
	}

	tracker.BuildManifest(round)

	manifest, _ := tracker.GetManifest()

	if len(manifest) != 1 {
		t.Errorf("Expected 1 media item, got %d", len(manifest))
	}

	if manifest[0].QuestionRef.QuestionPrice != 400 {
		t.Errorf("Expected question price 400, got %d", manifest[0].QuestionRef.QuestionPrice)
	}
}

func TestMediaTracker_GetManifest(t *testing.T) {
	tracker := NewMediaTracker(1)

	manifest, totalSize := tracker.GetManifest()

	if manifest == nil {
		t.Error("manifest should not be nil")
	}

	if len(manifest) != 0 {
		t.Errorf("Expected empty manifest, got %d items", len(manifest))
	}

	if totalSize != 0 {
		t.Errorf("Expected totalSize 0, got %d", totalSize)
	}

	round := &domain.Round{
		ID:          "round1",
		RoundNumber: 1,
		Name:        "Round 1",
		Themes: []*domain.Theme{
			{
				ID:   "theme1",
				Name: "Theme 1",
				Questions: []*domain.Question{
					{
						ID:        "q1",
						Price:      100,
						MediaType:  MediaTypeImage,
						MediaURL:   "http:
					},
				},
			},
		},
	}

	tracker.BuildManifest(round)

	manifest, totalSize = tracker.GetManifest()

	if len(manifest) != 1 {
		t.Errorf("Expected 1 media item, got %d", len(manifest))
	}

	expectedSize := int64(MediaSizeImage)
	if totalSize != expectedSize {
		t.Errorf("Expected totalSize %d, got %d", expectedSize, totalSize)
	}

	if manifest[0].Type != MediaTypeImage {
		t.Errorf("Expected media type %s, got %s", MediaTypeImage, manifest[0].Type)
	}

	if manifest[0].URL != "http:
		t.Errorf("Expected URL http:
	}
}

func TestMediaTracker_HasMedia(t *testing.T) {
	tracker := NewMediaTracker(1)

	if tracker.HasMedia() {
		t.Error("HasMedia should return false for empty tracker")
	}

	round := &domain.Round{
		ID:          "round1",
		RoundNumber: 1,
		Name:        "Round 1",
		Themes: []*domain.Theme{
			{
				ID:   "theme1",
				Name: "Theme 1",
				Questions: []*domain.Question{
					{
						ID:        "q1",
						Price:      100,
						MediaType:  MediaTypeImage,
						MediaURL:   "http:
					},
				},
			},
		},
	}

	tracker.BuildManifest(round)

	if !tracker.HasMedia() {
		t.Error("HasMedia should return true after building manifest with media")
	}
}

func TestMediaTracker_RegisterClient(t *testing.T) {
	tracker := NewMediaTracker(1)

	userID := uuid.New()
	tracker.RegisterClient(userID)

	if tracker.clients[userID] == nil {
		t.Fatal("Client should be registered")
	}

	if tracker.clients[userID].Total != 0 {
		t.Errorf("Expected Total 0 for empty manifest, got %d", tracker.clients[userID].Total)
	}

	if !tracker.clients[userID].Complete {
		t.Error("Client should be marked as complete when no media")
	}

	round := &domain.Round{
		ID:          "round1",
		RoundNumber: 1,
		Name:        "Round 1",
		Themes: []*domain.Theme{
			{
				ID:   "theme1",
				Name: "Theme 1",
				Questions: []*domain.Question{
					{
						ID:        "q1",
						Price:      100,
						MediaType:  MediaTypeImage,
						MediaURL:   "http:
					},
				},
			},
		},
	}

	tracker2 := NewMediaTracker(1)
	tracker2.BuildManifest(round)

	userID2 := uuid.New()
	tracker2.RegisterClient(userID2)

	manifest, _ := tracker2.GetManifest()
	expectedTotal := len(manifest)

	if tracker2.clients[userID2].Total != expectedTotal {
		t.Errorf("Expected Total %d, got %d", expectedTotal, tracker2.clients[userID2].Total)
	}

	if tracker2.clients[userID2].Complete {
		t.Error("Client should not be marked as complete when manifest has media")
	}
}

func TestMediaTracker_UpdateProgress(t *testing.T) {
	tracker := NewMediaTracker(1)

	userID := uuid.New()
	loaded := 5
	total := 10
	bytesLoaded := int64(500000)
	percent := 50

	tracker.UpdateProgress(userID, loaded, total, bytesLoaded, percent)

	status := tracker.clients[userID]
	if status == nil {
		t.Fatal("Client status should exist after UpdateProgress")
	}

	if status.Loaded != loaded {
		t.Errorf("Expected Loaded %d, got %d", loaded, status.Loaded)
	}

	if status.Total != total {
		t.Errorf("Expected Total %d, got %d", total, status.Total)
	}

	if status.BytesLoaded != bytesLoaded {
		t.Errorf("Expected BytesLoaded %d, got %d", bytesLoaded, status.BytesLoaded)
	}

	if status.Percent != percent {
		t.Errorf("Expected Percent %d, got %d", percent, status.Percent)
	}

	if status.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}

	tracker.UpdateProgress(userID, 8, total, bytesLoaded*2, 80)

	if tracker.clients[userID].Loaded != 8 {
		t.Errorf("Expected Loaded 8, got %d", tracker.clients[userID].Loaded)
	}

	if tracker.clients[userID].Percent != 80 {
		t.Errorf("Expected Percent 80, got %d", tracker.clients[userID].Percent)
	}
}

func TestMediaTracker_MarkComplete(t *testing.T) {
	tracker := NewMediaTracker(1)

	userID := uuid.New()
	loadedCount := 10

	tracker.MarkComplete(userID, loadedCount)

	status := tracker.clients[userID]
	if status == nil {
		t.Fatal("Client status should exist after MarkComplete")
	}

	if status.Loaded != loadedCount {
		t.Errorf("Expected Loaded %d, got %d", loadedCount, status.Loaded)
	}

	if !status.Complete {
		t.Error("Complete should be true")
	}

	if status.Percent != PercentComplete {
		t.Errorf("Expected Percent %d, got %d", PercentComplete, status.Percent)
	}

	if status.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

func TestMediaTracker_AllClientsReady(t *testing.T) {
	tracker := NewMediaTracker(1)

	if !tracker.AllClientsReady() {
		t.Error("AllClientsReady should return true when no clients registered")
	}

	userID1 := uuid.New()
	userID2 := uuid.New()

	tracker.RegisterClient(userID1)
	tracker.RegisterClient(userID2)

	if !tracker.AllClientsReady() {
		t.Error("AllClientsReady should return true when all clients complete")
	}

	tracker.UpdateProgress(userID1, 5, 10, 500000, 50)

	if tracker.AllClientsReady() {
		t.Error("AllClientsReady should return false when some clients not complete")
	}

	tracker.MarkComplete(userID1, 10)

	if !tracker.AllClientsReady() {
		t.Error("AllClientsReady should return true when all clients complete")
	}
}

func TestMediaTracker_GetPendingClients(t *testing.T) {
	tracker := NewMediaTracker(1)

	pending := tracker.GetPendingClients()
	if len(pending) != 0 {
		t.Errorf("Expected 0 pending clients, got %d", len(pending))
	}

	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()

	tracker.RegisterClient(userID1)
	tracker.RegisterClient(userID2)
	tracker.RegisterClient(userID3)

	pending = tracker.GetPendingClients()
	if len(pending) != 0 {
		t.Errorf("Expected 0 pending clients, got %d", len(pending))
	}

	tracker.UpdateProgress(userID1, 5, 10, 500000, 50)
	tracker.MarkComplete(userID2, 10)

	pending = tracker.GetPendingClients()
	if len(pending) != 2 {
		t.Errorf("Expected 2 pending clients, got %d", len(pending))
	}

	found := false
	for _, id := range pending {
		if id == userID1 || id == userID3 {
			found = true
		}
	}
	if !found {
		t.Error("Expected userID1 or userID3 in pending clients")
	}
}

func TestMediaTracker_GetOverallProgress(t *testing.T) {
	tracker := NewMediaTracker(1)

	totalPercent, clientsReady, totalClients := tracker.GetOverallProgress()

	if totalPercent != PercentComplete {
		t.Errorf("Expected totalPercent %d, got %d", PercentComplete, totalPercent)
	}

	if clientsReady != 0 {
		t.Errorf("Expected clientsReady 0, got %d", clientsReady)
	}

	if totalClients != 0 {
		t.Errorf("Expected totalClients 0, got %d", totalClients)
	}

	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()

	tracker.RegisterClient(userID1)
	tracker.RegisterClient(userID2)
	tracker.RegisterClient(userID3)

	tracker.UpdateProgress(userID1, 5, 10, 500000, 50)
	tracker.UpdateProgress(userID2, 8, 10, 800000, 80)
	tracker.MarkComplete(userID3, 10)

	totalPercent, clientsReady, totalClients = tracker.GetOverallProgress()

	expectedPercent := (50 + 80 + 100) / 3
	if totalPercent != expectedPercent {
		t.Errorf("Expected totalPercent %d, got %d", expectedPercent, totalPercent)
	}

	if clientsReady != 1 {
		t.Errorf("Expected clientsReady 1, got %d", clientsReady)
	}

	if totalClients != 3 {
		t.Errorf("Expected totalClients 3, got %d", totalClients)
	}
}

func TestMediaTracker_FindMediaByQuestion(t *testing.T) {
	tracker := NewMediaTracker(1)

	round := &domain.Round{
		ID:          "round1",
		RoundNumber: 1,
		Name:        "Round 1",
		Themes: []*domain.Theme{
			{
				ID:   "theme1",
				Name: "Theme 1",
				Questions: []*domain.Question{
					{
						ID:        "q1",
						Price:      100,
						MediaType:  MediaTypeImage,
						MediaURL:   "http:
					},
					{
						ID:        "q2",
						Price:      200,
						MediaType:  MediaTypeAudio,
						MediaURL:   "http:
					},
				},
			},
			{
				ID:   "theme2",
				Name: "Theme 2",
				Questions: []*domain.Question{
					{
						ID:        "q3",
						Price:      300,
						MediaType:  MediaTypeVideo,
						MediaURL:   "http:
					},
				},
			},
		},
	}

	tracker.BuildManifest(round)

	item := tracker.FindMediaByQuestion(0, 100)
	if item == nil {
		t.Fatal("Expected to find media item for theme 0, price 100")
	}

	if item.Type != MediaTypeImage {
		t.Errorf("Expected media type %s, got %s", MediaTypeImage, item.Type)
	}

	if item.URL != "http:
		t.Errorf("Expected URL http:
	}

	item = tracker.FindMediaByQuestion(1, 300)
	if item == nil {
		t.Fatal("Expected to find media item for theme 1, price 300")
	}

	if item.Type != MediaTypeVideo {
		t.Errorf("Expected media type %s, got %s", MediaTypeVideo, item.Type)
	}

	item = tracker.FindMediaByQuestion(0, 999)
	if item != nil {
		t.Error("Expected nil for non-existent question")
	}

	item = tracker.FindMediaByQuestion(999, 100)
	if item != nil {
		t.Error("Expected nil for non-existent theme")
	}
}

func TestMediaTracker_Reset(t *testing.T) {
	tracker := NewMediaTracker(1)

	round := &domain.Round{
		ID:          "round1",
		RoundNumber: 1,
		Name:        "Round 1",
		Themes: []*domain.Theme{
			{
				ID:   "theme1",
				Name: "Theme 1",
				Questions: []*domain.Question{
					{
						ID:        "q1",
						Price:      100,
						MediaType:  MediaTypeImage,
						MediaURL:   "http:
					},
				},
			},
		},
	}

	tracker.BuildManifest(round)

	userID := uuid.New()
	tracker.RegisterClient(userID)
	tracker.UpdateProgress(userID, 5, 10, 500000, 50)

	newRoundNumber := 2
	tracker.Reset(newRoundNumber)

	if tracker.roundNumber != newRoundNumber {
		t.Errorf("Expected roundNumber %d, got %d", newRoundNumber, tracker.roundNumber)
	}

	manifest, totalSize := tracker.GetManifest()
	if len(manifest) != 0 {
		t.Errorf("Expected empty manifest, got %d items", len(manifest))
	}

	if totalSize != 0 {
		t.Errorf("Expected totalSize 0, got %d", totalSize)
	}

	if len(tracker.clients) != 0 {
		t.Errorf("Expected empty clients map, got %d clients", len(tracker.clients))
	}

	if !tracker.HasMedia() {
		t.Error("HasMedia should return false after reset")
	}
}

func TestMediaLoadStatus(t *testing.T) {
	userID := uuid.New()
	status := MediaLoadStatus{
		UserID:      userID,
		Loaded:      5,
		Total:       10,
		BytesLoaded: 500000,
		Percent:     50,
		Complete:    false,
		UpdatedAt:   time.Now(),
	}

	if status.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, status.UserID)
	}

	if status.Loaded != 5 {
		t.Errorf("Expected Loaded 5, got %d", status.Loaded)
	}

	if status.Total != 10 {
		t.Errorf("Expected Total 10, got %d", status.Total)
	}

	if status.BytesLoaded != 500000 {
		t.Errorf("Expected BytesLoaded 500000, got %d", status.BytesLoaded)
	}

	if status.Percent != 50 {
		t.Errorf("Expected Percent 50, got %d", status.Percent)
	}

	if status.Complete {
		t.Error("Expected Complete false")
	}

	if status.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

