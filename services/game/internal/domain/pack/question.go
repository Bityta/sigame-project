package pack

import "strings"

type Type string

const (
	TypeNormal Type = "normal"
	TypeSecret Type = "secret"
	TypeStake  Type = "stake"
	TypeForAll Type = "forAll"
)

func (t Type) String() string {
	return string(t)
}

func (t Type) IsSpecial() bool {
	return t == TypeSecret || t == TypeStake || t == TypeForAll
}

type Question struct {
	ID              string
	Price           int
	Text            string
	Answer          string
	Type            Type
	MediaType       string
	MediaURL        string
	MediaDurationMs int
	Used            bool
}

type QuestionState struct {
	ID              string `json:"id" binding:"required"`
	Price           int    `json:"price" binding:"required"`
	Available       bool   `json:"available" binding:"required"`
	Type            string `json:"type" binding:"required"`
	Text            string `json:"text,omitempty"`
	MediaType       string `json:"mediaType,omitempty"`
	MediaURL        string `json:"mediaUrl,omitempty"`
	MediaDurationMs int    `json:"mediaDurationMs,omitempty"`
	Answer          string `json:"answer,omitempty"`
}

func (q *Question) GetType() Type {
	if q.Type == "" {
		return TypeNormal
	}
	return q.Type
}

func (q *Question) IsSpecialType() bool {
	return q.GetType().IsSpecial()
}

func (q *Question) ValidateAnswer(userAnswer string) bool {
	normalized := strings.TrimSpace(strings.ToLower(userAnswer))
	correctAnswer := strings.TrimSpace(strings.ToLower(q.Answer))
	return normalized == correctAnswer
}

func (q *Question) MarkAsUsed() {
	q.Used = true
}

func (q *Question) IsAvailable() bool {
	return !q.Used
}

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

func (q *Question) ToStateWithAnswer(includeText bool) QuestionState {
	state := q.ToState(includeText)
	state.Answer = q.Answer
	return state
}

func (q *Question) HasMedia() bool {
	return q.MediaType != "" && q.MediaType != "text" && q.MediaURL != ""
}

