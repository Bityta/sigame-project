package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
)

// PackServiceClient defines the interface for pack service client
type PackServiceClient interface {
	GetPackContent(ctx context.Context, packID uuid.UUID) (*domain.Pack, error)
	ValidatePackExists(ctx context.Context, packID uuid.UUID) (bool, error)
}

// PackClient is a HTTP client for Pack Service (temporary, will migrate to gRPC later)
type PackClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewPackClient creates a new Pack Service HTTP client
func NewPackClient(address string) (*PackClient, error) {
	return &PackClient{
		baseURL: fmt.Sprintf("http://%s", address),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// Close closes the HTTP client (no-op for HTTP, kept for interface compatibility)
func (c *PackClient) Close() error {
	return nil
}

// PackContentResponse represents the JSON response from Pack Service
type PackContentResponse struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Author string       `json:"author"`
	Rounds []RoundJSON  `json:"rounds"`
	Error  string       `json:"error,omitempty"`
}

type RoundJSON struct {
	ID          string      `json:"id"`
	RoundNumber int         `json:"round_number"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Themes      []ThemeJSON `json:"themes"`
}

type ThemeJSON struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Questions []QuestionJSON `json:"questions"`
}

type QuestionJSON struct {
	ID              string   `json:"id"`
	Price           int      `json:"price"`
	Type            string   `json:"type"`
	Text            string   `json:"text"`
	Answer          string   `json:"answer"`
	AltAnswers      []string `json:"alt_answers"`
	MediaType       string   `json:"media_type"`
	MediaURL        string   `json:"media_url"`
	MediaDurationMs int      `json:"media_duration_ms"`
}

// GetPackContent retrieves full pack content from Pack Service via HTTP
func (c *PackClient) GetPackContent(ctx context.Context, packID uuid.UUID) (*domain.Pack, error) {
	url := fmt.Sprintf("%s/api/packs/%s/content", c.baseURL, packID.String())
	
	log.Printf("[PackClient] Fetching pack content from: %s", url)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pack content: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pack service returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var packResp PackContentResponse
	if err := json.Unmarshal(body, &packResp); err != nil {
		return nil, fmt.Errorf("failed to parse pack content: %w", err)
	}
	
	log.Printf("[PackClient] Got pack: %s with %d rounds", packResp.Name, len(packResp.Rounds))
	
	// Convert to domain.Pack
	pack := &domain.Pack{
		ID:     packID,
		Name:   packResp.Name,
		Author: packResp.Author,
		Rounds: make([]*domain.Round, len(packResp.Rounds)),
	}
	
	for i, r := range packResp.Rounds {
		round := &domain.Round{
			ID:          r.ID,
			RoundNumber: r.RoundNumber,
			Name:        r.Name,
			Themes:      make([]*domain.Theme, len(r.Themes)),
		}
		
		for j, t := range r.Themes {
			theme := &domain.Theme{
				ID:        t.ID,
				Name:      t.Name,
				Questions: make([]*domain.Question, len(t.Questions)),
			}
			
			for k, q := range t.Questions {
				theme.Questions[k] = &domain.Question{
					ID:        q.ID,
					Price:     q.Price,
					Text:      q.Text,
					Answer:    q.Answer,
					MediaType: q.MediaType,
					Used:      false,
				}
			}
			
			round.Themes[j] = theme
		}
		
		pack.Rounds[i] = round
	}
	
	log.Printf("[PackClient] Converted pack with %d rounds, first round has %d themes", 
		len(pack.Rounds), 
		func() int { if len(pack.Rounds) > 0 { return len(pack.Rounds[0].Themes) }; return 0 }())
	
	return pack, nil
}

// ValidatePackExists checks if a pack exists via HTTP
func (c *PackClient) ValidatePackExists(ctx context.Context, packID uuid.UUID) (bool, error) {
	url := fmt.Sprintf("%s/api/packs/%s", c.baseURL, packID.String())
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to validate pack: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("pack service returned status %d", resp.StatusCode)
	}
	
	return true, nil
}
