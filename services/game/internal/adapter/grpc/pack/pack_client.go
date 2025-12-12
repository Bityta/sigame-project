package pack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/infrastructure/logger"
)

type PackServiceClient interface {
	GetPackContent(ctx context.Context, packID uuid.UUID) (*pack.Pack, error)
	ValidatePackExists(ctx context.Context, packID uuid.UUID) (bool, error)
}

type PackClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPackClient(address string) (*PackClient, error) {
	return &PackClient{
		baseURL:    fmt.Sprintf("http://%s", address),
		httpClient: http.DefaultClient,
	}, nil
}

func (c *PackClient) Close() error {
	return nil
}

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

func (c *PackClient) GetPackContent(ctx context.Context, packID uuid.UUID) (*pack.Pack, error) {
	url := buildURL(c.baseURL, PathPackContent, packID.String())

	logger.Debugf(ctx, "Fetching pack content from: %s", url)

	body, statusCode, err := doGetRequest(ctx, c.httpClient, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pack content: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("pack service returned status %d: %s", statusCode, string(body))
	}

	var packResp PackContentResponse
	if err := json.Unmarshal(body, &packResp); err != nil {
		return nil, fmt.Errorf("failed to parse pack content: %w", err)
	}

	logger.Debugf(ctx, "Got pack: %s with %d rounds", packResp.Name, len(packResp.Rounds))

	pack, err := convertPackResponse(&packResp, packID)
	if err != nil {
		return nil, err
	}

	logger.Debugf(ctx, "Converted pack with %d rounds", len(pack.Rounds))

	return pack, nil
}

func (c *PackClient) ValidatePackExists(ctx context.Context, packID uuid.UUID) (bool, error) {
	url := buildURL(c.baseURL, PathPack, packID.String())

	_, statusCode, err := doGetRequest(ctx, c.httpClient, url)
	if err != nil {
		return false, fmt.Errorf("failed to validate pack: %w", err)
	}

	if statusCode == http.StatusNotFound {
		return false, nil
	}

	if statusCode != http.StatusOK {
		return false, fmt.Errorf("pack service returned status %d", statusCode)
	}

	return true, nil
}
