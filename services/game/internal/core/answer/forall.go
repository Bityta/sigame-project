package answer

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type ForAllAnswer struct {
	UserID      uuid.UUID
	Username    string
	Answer      string
	SubmittedAt time.Time
}

type ForAllCollector struct {
	answers       map[uuid.UUID]*ForAllAnswer
	correctAnswer string
	questionPrice int
	startedAt     time.Time
	closed        bool
	mu            sync.Mutex
}

func NewForAllCollector() *ForAllCollector {
	return &ForAllCollector{
		answers: make(map[uuid.UUID]*ForAllAnswer),
	}
}

func (c *ForAllCollector) Start(correctAnswer string, questionPrice int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.answers = make(map[uuid.UUID]*ForAllAnswer)
	c.correctAnswer = correctAnswer
	c.questionPrice = questionPrice
	c.startedAt = time.Now()
	c.closed = false
}

func (c *ForAllCollector) SubmitAnswer(userID uuid.UUID, username, answer string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return false
	}

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

func (c *ForAllCollector) HasAnswered(userID uuid.UUID) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.answers[userID]
	return exists
}

func (c *ForAllCollector) GetAnswerCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.answers)
}

func (c *ForAllCollector) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
}

func (c *ForAllCollector) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

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

func (c *ForAllCollector) GetAllAnswers() []*ForAllAnswer {
	c.mu.Lock()
	defer c.mu.Unlock()

	answers := make([]*ForAllAnswer, 0, len(c.answers))
	for _, answer := range c.answers {
		answers = append(answers, answer)
	}
	return answers
}

func (c *ForAllCollector) GetCorrectAnswer() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.correctAnswer
}

func (c *ForAllCollector) GetQuestionPrice() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.questionPrice
}

func (c *ForAllCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.answers = make(map[uuid.UUID]*ForAllAnswer)
	c.correctAnswer = ""
	c.questionPrice = 0
	c.closed = false
}

type ForAllResult struct {
	UserID     uuid.UUID
	Username   string
	Answer     string
	IsCorrect  bool
	ScoreDelta int
}

