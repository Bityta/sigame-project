package domain

import (
	"strings"

	"github.com/google/uuid"
)

// Question represents a question in the game
type Question struct {
	ID        string `json:"id"`
	Price     int    `json:"price"`
	Text      string `json:"text"`
	Answer    string `json:"answer"`
	MediaType string `json:"media_type"`
	Used      bool   `json:"used"`
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
	}

	if includeText {
		state.Text = q.Text
		state.MediaType = q.MediaType
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

