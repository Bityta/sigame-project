package grpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestNewPackClient(t *testing.T) {
	client, err := NewPackClient("localhost:8001")
	if err != nil {
		t.Fatalf("NewPackClient() error = %v", err)
	}

	if client == nil {
		t.Error("NewPackClient() returned nil client")
	}

	if client.baseURL != "http:
		t.Errorf("NewPackClient() baseURL = %s, want http:
	}

	if client.httpClient == nil {
		t.Error("NewPackClient() httpClient is nil")
	}
}

func TestPackClient_Close(t *testing.T) {
	client, _ := NewPackClient("localhost:8001")
	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

func TestPackClient_GetPackContent_Success(t *testing.T) {
	packID := uuid.New()
	expectedPack := PackContentResponse{
		ID:     packID.String(),
		Name:   "Test Pack",
		Author: "Test Author",
		Rounds: []RoundJSON{
			{
				ID:          "r1",
				RoundNumber: 1,
				Name:        "Round 1",
				Themes: []ThemeJSON{
					{
						ID:   "t1",
						Name: "Theme 1",
						Questions: []QuestionJSON{
							{
								ID:              "q1",
								Price:           100,
								Text:            "Question 1",
								Answer:          "Answer 1",
								MediaType:       "image",
								MediaURL:        "http:
								MediaDurationMs: 5000,
							},
						},
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedPack)
	}))
	defer server.Close()

	client, _ := NewPackClient(server.URL[7:])
	client.baseURL = server.URL

	ctx := context.Background()
	pack, err := client.GetPackContent(ctx, packID)

	if err != nil {
		t.Fatalf("GetPackContent() error = %v", err)
	}

	if pack.ID != packID {
		t.Errorf("GetPackContent() pack.ID = %v, want %v", pack.ID, packID)
	}

	if pack.Name != expectedPack.Name {
		t.Errorf("GetPackContent() pack.Name = %s, want %s", pack.Name, expectedPack.Name)
	}

	if len(pack.Rounds) != len(expectedPack.Rounds) {
		t.Errorf("GetPackContent() pack.Rounds length = %d, want %d", len(pack.Rounds), len(expectedPack.Rounds))
	}
}

func TestPackClient_GetPackContent_NotFound(t *testing.T) {
	packID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, _ := NewPackClient(server.URL[7:])
	client.baseURL = server.URL

	ctx := context.Background()
	_, err := client.GetPackContent(ctx, packID)

	if err == nil {
		t.Error("GetPackContent() should return error on 404")
	}
}

func TestPackClient_GetPackContent_InvalidJSON(t *testing.T) {
	packID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client, _ := NewPackClient(server.URL[7:])
	client.baseURL = server.URL

	ctx := context.Background()
	_, err := client.GetPackContent(ctx, packID)

	if err == nil {
		t.Error("GetPackContent() should return error on invalid JSON")
	}
}

func TestPackClient_ValidatePackExists_Exists(t *testing.T) {
	packID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewPackClient(server.URL[7:])
	client.baseURL = server.URL

	ctx := context.Background()
	exists, err := client.ValidatePackExists(ctx, packID)

	if err != nil {
		t.Fatalf("ValidatePackExists() error = %v", err)
	}

	if !exists {
		t.Error("ValidatePackExists() = false, want true")
	}
}

func TestPackClient_ValidatePackExists_NotFound(t *testing.T) {
	packID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, _ := NewPackClient(server.URL[7:])
	client.baseURL = server.URL

	ctx := context.Background()
	exists, err := client.ValidatePackExists(ctx, packID)

	if err != nil {
		t.Fatalf("ValidatePackExists() error = %v", err)
	}

	if exists {
		t.Error("ValidatePackExists() = true, want false")
	}
}

func TestPackClient_ValidatePackExists_ServerError(t *testing.T) {
	packID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := NewPackClient(server.URL[7:])
	client.baseURL = server.URL

	ctx := context.Background()
	_, err := client.ValidatePackExists(ctx, packID)

	if err == nil {
		t.Error("ValidatePackExists() should return error on server error")
	}
}

func TestPackClient_ImplementsInterface(t *testing.T) {
	var _ PackServiceClient = (*PackClient)(nil)
}

