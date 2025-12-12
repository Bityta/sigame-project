package pack

import (
	"testing"

	"github.com/google/uuid"
)

func TestConvertQuestion(t *testing.T) {
	q := QuestionJSON{
		ID:              "q1",
		Price:          100,
		Text:            "Test question",
		Answer:          "Answer",
		MediaType:       "image",
		MediaURL:        "http:
		MediaDurationMs: 5000,
	}

	result := convertQuestion(q)

	if result.ID != q.ID {
		t.Errorf("convertQuestion() ID = %s, want %s", result.ID, q.ID)
	}
	if result.Price != q.Price {
		t.Errorf("convertQuestion() Price = %d, want %d", result.Price, q.Price)
	}
	if result.Text != q.Text {
		t.Errorf("convertQuestion() Text = %s, want %s", result.Text, q.Text)
	}
	if result.Answer != q.Answer {
		t.Errorf("convertQuestion() Answer = %s, want %s", result.Answer, q.Answer)
	}
	if result.Used != false {
		t.Error("convertQuestion() Used should be false")
	}
}

func TestConvertTheme(t *testing.T) {
	tm := ThemeJSON{
		ID:   "t1",
		Name: "Theme 1",
		Questions: []QuestionJSON{
			{ID: "q1", Price: 100, Text: "Q1", Answer: "A1"},
			{ID: "q2", Price: 200, Text: "Q2", Answer: "A2"},
		},
	}

	result := convertTheme(tm)

	if result.ID != tm.ID {
		t.Errorf("convertTheme() ID = %s, want %s", result.ID, tm.ID)
	}
	if result.Name != tm.Name {
		t.Errorf("convertTheme() Name = %s, want %s", result.Name, tm.Name)
	}
	if len(result.Questions) != len(tm.Questions) {
		t.Errorf("convertTheme() Questions length = %d, want %d", len(result.Questions), len(tm.Questions))
	}
}

func TestConvertRound(t *testing.T) {
	r := RoundJSON{
		ID:          "r1",
		RoundNumber: 1,
		Name:        "Round 1",
		Themes: []ThemeJSON{
			{ID: "t1", Name: "Theme 1", Questions: []QuestionJSON{}},
			{ID: "t2", Name: "Theme 2", Questions: []QuestionJSON{}},
		},
	}

	result := convertRound(r)

	if result.ID != r.ID {
		t.Errorf("convertRound() ID = %s, want %s", result.ID, r.ID)
	}
	if result.RoundNumber != r.RoundNumber {
		t.Errorf("convertRound() RoundNumber = %d, want %d", result.RoundNumber, r.RoundNumber)
	}
	if len(result.Themes) != len(r.Themes) {
		t.Errorf("convertRound() Themes length = %d, want %d", len(result.Themes), len(r.Themes))
	}
}

func TestConvertPackResponse(t *testing.T) {
	packID := uuid.New()
	packResp := &PackContentResponse{
		ID:     packID.String(),
		Name:   "Test Pack",
		Author: "Test Author",
		Rounds: []RoundJSON{
			{
				ID:          "r1",
				RoundNumber: 1,
				Name:        "Round 1",
				Themes:      []ThemeJSON{},
			},
		},
	}

	result, err := convertPackResponse(packResp, packID)
	if err != nil {
		t.Fatalf("convertPackResponse() error = %v", err)
	}

	if result.ID != packID {
		t.Errorf("convertPackResponse() ID = %v, want %v", result.ID, packID)
	}
	if result.Name != packResp.Name {
		t.Errorf("convertPackResponse() Name = %s, want %s", result.Name, packResp.Name)
	}
	if len(result.Rounds) != len(packResp.Rounds) {
		t.Errorf("convertPackResponse() Rounds length = %d, want %d", len(result.Rounds), len(packResp.Rounds))
	}
}

func TestConvertPackResponse_WithError(t *testing.T) {
	packID := uuid.New()
	packResp := &PackContentResponse{
		Error: "Pack not found",
	}

	_, err := convertPackResponse(packResp, packID)
	if err == nil {
		t.Error("convertPackResponse() should return error when Error field is set")
	}
}

