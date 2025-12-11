package game

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// ForAllAnswer represents a single player's answer in a forAll question
type ForAllAnswer struct {
	UserID      uuid.UUID
	Username    string
	Answer      string
	SubmittedAt time.Time
}

// ForAllCollector collects answers from all players for a forAll question
type ForAllCollector struct {
	answers       map[uuid.UUID]*ForAllAnswer
	correctAnswer string
	questionPrice int
	startedAt     time.Time
	closed        bool
	mu            sync.Mutex
}

// NewForAllCollector creates a new collector for forAll answers
func NewForAllCollector() *ForAllCollector {
	return &ForAllCollector{
		answers: make(map[uuid.UUID]*ForAllAnswer),
	}
}

// Start initializes the collector for a new forAll question
func (c *ForAllCollector) Start(correctAnswer string, questionPrice int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.answers = make(map[uuid.UUID]*ForAllAnswer)
	c.correctAnswer = correctAnswer
	c.questionPrice = questionPrice
	c.startedAt = time.Now()
	c.closed = false
}

// SubmitAnswer records a player's answer
func (c *ForAllCollector) SubmitAnswer(userID uuid.UUID, username, answer string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return false
	}

	// Check if already answered
	if _, exists := c.answers[userID]; exists {
		return false
	}

	c.answers[userID] = &ForAllAnswer{
		UserID:      userID,
		Username:    username,
		Answer:      answer,
		SubmittedAt: time.Now(),
	}

	return true
}

// HasAnswered checks if a player has already submitted an answer
func (c *ForAllCollector) HasAnswered(userID uuid.UUID) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.answers[userID]
	return exists
}

// GetAnswerCount returns the number of answers collected
func (c *ForAllCollector) GetAnswerCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.answers)
}

// Close closes the collector to prevent more submissions
func (c *ForAllCollector) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
}

// IsClosed returns whether the collector is closed
func (c *ForAllCollector) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// GetResults calculates and returns the results for all players
// Returns map of userID -> (isCorrect, scoreDelta)
func (c *ForAllCollector) GetResults(validateAnswer func(userAnswer, correctAnswer string) bool) map[uuid.UUID]ForAllResult {
	c.mu.Lock()
	defer c.mu.Unlock()

	results := make(map[uuid.UUID]ForAllResult)

	for userID, answer := range c.answers {
		isCorrect := validateAnswer(answer.Answer, c.correctAnswer)
		scoreDelta := 0
		if isCorrect {
			scoreDelta = c.questionPrice
		} else {
			scoreDelta = -c.questionPrice
		}

		results[userID] = ForAllResult{
			UserID:     userID,
			Username:   answer.Username,
			Answer:     answer.Answer,
			IsCorrect:  isCorrect,
			ScoreDelta: scoreDelta,
		}
	}

	return results
}

// GetAllAnswers returns all submitted answers
func (c *ForAllCollector) GetAllAnswers() []*ForAllAnswer {
	c.mu.Lock()
	defer c.mu.Unlock()

	answers := make([]*ForAllAnswer, 0, len(c.answers))
	for _, answer := range c.answers {
		answers = append(answers, answer)
	}
	return answers
}

// GetCorrectAnswer returns the correct answer
func (c *ForAllCollector) GetCorrectAnswer() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.correctAnswer
}

// GetQuestionPrice returns the question price
func (c *ForAllCollector) GetQuestionPrice() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.questionPrice
}

// Reset clears the collector
func (c *ForAllCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.answers = make(map[uuid.UUID]*ForAllAnswer)
	c.correctAnswer = ""
	c.questionPrice = 0
	c.closed = false
}

// ForAllResult represents the result for a single player
type ForAllResult struct {
	UserID     uuid.UUID
	Username   string
	Answer     string
	IsCorrect  bool
	ScoreDelta int
}

