package domain

import (
	"strings"

	"github.com/google/uuid"
)

// QuestionType represents the type of question
type QuestionType string

// Question type constants
const (
	QuestionTypeNormal QuestionType = "normal" // Обычный вопрос
	QuestionTypeSecret QuestionType = "secret" // Кот в мешке - передаётся другому игроку
	QuestionTypeStake  QuestionType = "stake"  // Ва-банк - игрок делает ставку
	QuestionTypeForAll QuestionType = "forAll" // Вопрос для всех одновременно
)

// Question represents a question in the game
type Question struct {
	ID              string       `json:"id"`
	Price           int          `json:"price"`
	Text            string       `json:"text"`
	Answer          string       `json:"answer"`
	Type            QuestionType `json:"type"`
	MediaType       string       `json:"media_type"`
	MediaURL        string       `json:"media_url"`
	MediaDurationMs int          `json:"media_duration_ms"`
	Used            bool         `json:"used"`
}

// GetType returns the question type, defaulting to normal if not set
func (q *Question) GetType() QuestionType {
	if q.Type == "" {
		return QuestionTypeNormal
	}
	return q.Type
}

// IsSpecialType returns true if the question has a special type (not normal)
func (q *Question) IsSpecialType() bool {
	t := q.GetType()
	return t == QuestionTypeSecret || t == QuestionTypeStake || t == QuestionTypeForAll
}

// Theme represents a theme with questions
type Theme struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Questions []*Question `json:"questions"`
}

// Round represents a round in the game
type Round struct {
	ID          string   `json:"id"`
	RoundNumber int      `json:"round_number"`
	Name        string   `json:"name"`
	Themes      []*Theme `json:"themes"`
}

// Pack represents the complete question pack
type Pack struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Rounds      []*Round  `json:"rounds"`
}

// ValidateAnswer checks if the given answer matches the correct answer
func (q *Question) ValidateAnswer(userAnswer string) bool {
	// Normalize both answers: trim spaces and convert to lowercase
	normalized := strings.TrimSpace(strings.ToLower(userAnswer))
	correctAnswer := strings.TrimSpace(strings.ToLower(q.Answer))

	// Simple exact match for now
	// TODO: Add fuzzy matching or Levenshtein distance for typos
	return normalized == correctAnswer
}

// MarkAsUsed marks the question as used
func (q *Question) MarkAsUsed() {
	q.Used = true
}

// IsAvailable checks if the question can be selected
func (q *Question) IsAvailable() bool {
	return !q.Used
}

// ToState converts Question to QuestionState
func (q *Question) ToState(includeText bool) QuestionState {
	state := QuestionState{
		ID:        q.ID,
		Price:     q.Price,
		Available: !q.Used,
		Type:      string(q.GetType()),
	}

	if includeText {
		state.Text = q.Text
		state.MediaType = q.MediaType
		state.MediaURL = q.MediaURL
		state.MediaDurationMs = q.MediaDurationMs
	}

	return state
}

// ToStateWithAnswer converts Question to QuestionState including the answer (for host)
func (q *Question) ToStateWithAnswer(includeText bool) QuestionState {
	state := q.ToState(includeText)
	state.Answer = q.Answer
	return state
}

// GetAvailableQuestions returns all available questions in a theme
func (t *Theme) GetAvailableQuestions() []*Question {
	available := make([]*Question, 0)
	for _, q := range t.Questions {
		if q.IsAvailable() {
			available = append(available, q)
		}
	}
	return available
}

// ToState converts Theme to ThemeState
func (t *Theme) ToState(includeQuestionText bool) ThemeState {
	questions := make([]QuestionState, len(t.Questions))
	for i, q := range t.Questions {
		questions[i] = q.ToState(includeQuestionText)
	}

	return ThemeState{
		Name:      t.Name,
		Questions: questions,
	}
}

// GetAvailableThemes returns themes that still have available questions
func (r *Round) GetAvailableThemes() []*Theme {
	available := make([]*Theme, 0)
	for _, theme := range r.Themes {
		if len(theme.GetAvailableQuestions()) > 0 {
			available = append(available, theme)
		}
	}
	return available
}

// IsComplete checks if all questions in the round have been used
func (r *Round) IsComplete() bool {
	for _, theme := range r.Themes {
		for _, question := range theme.Questions {
			if !question.Used {
				return false
			}
		}
	}
	return true
}

// FindQuestion finds a question by theme and question ID
func (r *Round) FindQuestion(themeID, questionID string) *Question {
	for _, theme := range r.Themes {
		if theme.ID == themeID || theme.Name == themeID {
			for _, q := range theme.Questions {
				if q.ID == questionID {
					return q
				}
			}
		}
	}
	return nil
}

// TotalQuestions returns the total number of questions in the pack
func (p *Pack) TotalQuestions() int {
	total := 0
	for _, round := range p.Rounds {
		for _, theme := range round.Themes {
			total += len(theme.Questions)
		}
	}
	return total
}

