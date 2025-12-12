package answer

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewForAllCollector(t *testing.T) {
	collector := NewForAllCollector()

	if collector == nil {
		t.Fatal("NewForAllCollector() returned nil")
	}
	if collector.answers == nil {
		t.Error("NewForAllCollector() answers is nil")
	}
	if len(collector.answers) != 0 {
		t.Errorf("NewForAllCollector() answers length = %d, want 0", len(collector.answers))
	}
}

func TestForAllCollector_Start(t *testing.T) {
	collector := NewForAllCollector()
	correctAnswer := "Correct Answer"
	questionPrice := 100

	collector.Start(correctAnswer, questionPrice)

	if collector.correctAnswer != correctAnswer {
		t.Errorf("Start() correctAnswer = %s, want %s", collector.correctAnswer, correctAnswer)
	}
	if collector.questionPrice != questionPrice {
		t.Errorf("Start() questionPrice = %d, want %d", collector.questionPrice, questionPrice)
	}
	if collector.startedAt.IsZero() {
		t.Error("Start() startedAt is zero")
	}
	if collector.closed {
		t.Error("Start() closed = true, want false")
	}
	if len(collector.answers) != 0 {
		t.Errorf("Start() answers length = %d, want 0", len(collector.answers))
	}
}

func TestForAllCollector_SubmitAnswer(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	userID := uuid.New()
	username := "testuser"
	answer := "My Answer"

	success := collector.SubmitAnswer(userID, username, answer)

	if !success {
		t.Error("SubmitAnswer() returned false, want true")
	}
	if !collector.HasAnswered(userID) {
		t.Error("SubmitAnswer() HasAnswered() = false, want true")
	}
	if collector.GetAnswerCount() != 1 {
		t.Errorf("SubmitAnswer() GetAnswerCount() = %d, want 1", collector.GetAnswerCount())
	}

	answers := collector.GetAllAnswers()
	if len(answers) != 1 {
		t.Errorf("SubmitAnswer() GetAllAnswers() length = %d, want 1", len(answers))
	}
	if answers[0].UserID != userID {
		t.Errorf("SubmitAnswer() UserID = %v, want %v", answers[0].UserID, userID)
	}
	if answers[0].Answer != answer {
		t.Errorf("SubmitAnswer() Answer = %s, want %s", answers[0].Answer, answer)
	}
}

func TestForAllCollector_SubmitAnswer_Duplicate(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	userID := uuid.New()

	collector.SubmitAnswer(userID, "user1", "Answer 1")
	success := collector.SubmitAnswer(userID, "user1", "Answer 2")

	if success {
		t.Error("SubmitAnswer() duplicate returned true, want false")
	}
	if collector.GetAnswerCount() != 1 {
		t.Errorf("SubmitAnswer() duplicate GetAnswerCount() = %d, want 1", collector.GetAnswerCount())
	}

	answers := collector.GetAllAnswers()
	if answers[0].Answer != "Answer 1" {
		t.Errorf("SubmitAnswer() duplicate Answer = %s, want Answer 1", answers[0].Answer)
	}
}

func TestForAllCollector_SubmitAnswer_Closed(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)
	collector.Close()

	userID := uuid.New()
	success := collector.SubmitAnswer(userID, "user1", "Answer")

	if success {
		t.Error("SubmitAnswer() after Close() returned true, want false")
	}
	if collector.GetAnswerCount() != 0 {
		t.Errorf("SubmitAnswer() after Close() GetAnswerCount() = %d, want 0", collector.GetAnswerCount())
	}
}

func TestForAllCollector_HasAnswered(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	userID := uuid.New()

	if collector.HasAnswered(userID) {
		t.Error("HasAnswered() empty = true, want false")
	}

	collector.SubmitAnswer(userID, "user1", "Answer")

	if !collector.HasAnswered(userID) {
		t.Error("HasAnswered() with answer = false, want true")
	}
}

func TestForAllCollector_GetAnswerCount(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	if collector.GetAnswerCount() != 0 {
		t.Errorf("GetAnswerCount() empty = %d, want 0", collector.GetAnswerCount())
	}

	collector.SubmitAnswer(uuid.New(), "user1", "Answer 1")
	collector.SubmitAnswer(uuid.New(), "user2", "Answer 2")

	if collector.GetAnswerCount() != 2 {
		t.Errorf("GetAnswerCount() = %d, want 2", collector.GetAnswerCount())
	}
}

func TestForAllCollector_Close(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	collector.Close()

	if !collector.IsClosed() {
		t.Error("Close() IsClosed() = false, want true")
	}
}

func TestForAllCollector_IsClosed(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	if collector.IsClosed() {
		t.Error("IsClosed() = true, want false")
	}

	collector.Close()

	if !collector.IsClosed() {
		t.Error("IsClosed() after Close() = false, want true")
	}
}

func TestForAllCollector_GetResults(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Correct", 100)

	user1 := uuid.New()
	user2 := uuid.New()
	user3 := uuid.New()

	collector.SubmitAnswer(user1, "user1", "Correct")
	collector.SubmitAnswer(user2, "user2", "Wrong")
	collector.SubmitAnswer(user3, "user3", "correct")

	validateAnswer := func(userAnswer, correctAnswer string) bool {
		return strings.TrimSpace(strings.ToLower(userAnswer)) == strings.TrimSpace(strings.ToLower(correctAnswer))
	}

	results := collector.GetResults(validateAnswer)

	if len(results) != 3 {
		t.Errorf("GetResults() length = %d, want 3", len(results))
	}

	if !results[user1].IsCorrect {
		t.Error("GetResults() user1 IsCorrect = false, want true")
	}
	if results[user1].ScoreDelta != 100 {
		t.Errorf("GetResults() user1 ScoreDelta = %d, want 100", results[user1].ScoreDelta)
	}

	if results[user2].IsCorrect {
		t.Error("GetResults() user2 IsCorrect = true, want false")
	}
	if results[user2].ScoreDelta != -100 {
		t.Errorf("GetResults() user2 ScoreDelta = %d, want -100", results[user2].ScoreDelta)
	}

	if !results[user3].IsCorrect {
		t.Error("GetResults() user3 IsCorrect = false, want true")
	}
	if results[user3].ScoreDelta != 100 {
		t.Errorf("GetResults() user3 ScoreDelta = %d, want 100", results[user3].ScoreDelta)
	}
}

func TestForAllCollector_GetAllAnswers(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	if collector.GetAllAnswers() != nil && len(collector.GetAllAnswers()) != 0 {
		t.Error("GetAllAnswers() empty = non-empty, want empty")
	}

	user1 := uuid.New()
	user2 := uuid.New()

	collector.SubmitAnswer(user1, "user1", "Answer 1")
	collector.SubmitAnswer(user2, "user2", "Answer 2")

	answers := collector.GetAllAnswers()

	if len(answers) != 2 {
		t.Errorf("GetAllAnswers() length = %d, want 2", len(answers))
	}

	userIDs := make(map[uuid.UUID]bool)
	for _, answer := range answers {
		userIDs[answer.UserID] = true
	}

	if !userIDs[user1] || !userIDs[user2] {
		t.Error("GetAllAnswers() missing expected user IDs")
	}
}

func TestForAllCollector_GetCorrectAnswer(t *testing.T) {
	collector := NewForAllCollector()
	correctAnswer := "Correct Answer"

	collector.Start(correctAnswer, 100)

	if collector.GetCorrectAnswer() != correctAnswer {
		t.Errorf("GetCorrectAnswer() = %s, want %s", collector.GetCorrectAnswer(), correctAnswer)
	}
}

func TestForAllCollector_GetQuestionPrice(t *testing.T) {
	collector := NewForAllCollector()
	questionPrice := 200

	collector.Start("Answer", questionPrice)

	if collector.GetQuestionPrice() != questionPrice {
		t.Errorf("GetQuestionPrice() = %d, want %d", collector.GetQuestionPrice(), questionPrice)
	}
}

func TestForAllCollector_Reset(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	userID := uuid.New()
	collector.SubmitAnswer(userID, "user1", "Answer")
	collector.Close()

	collector.Reset()

	if collector.closed {
		t.Error("Reset() closed = true, want false")
	}
	if len(collector.answers) != 0 {
		t.Errorf("Reset() answers length = %d, want 0", len(collector.answers))
	}
	if collector.correctAnswer != "" {
		t.Errorf("Reset() correctAnswer = %s, want empty", collector.correctAnswer)
	}
	if collector.questionPrice != 0 {
		t.Errorf("Reset() questionPrice = %d, want 0", collector.questionPrice)
	}
}

func TestForAllCollector_ConcurrentAccess(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	done := make(chan bool)
	users := make([]uuid.UUID, 10)
	for i := range users {
		users[i] = uuid.New()
	}

	for i := 0; i < 10; i++ {
		go func(idx int) {
			collector.SubmitAnswer(users[idx], "user", "Answer")
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if collector.GetAnswerCount() != 10 {
		t.Errorf("Concurrent access GetAnswerCount() = %d, want 10", collector.GetAnswerCount())
	}
}

func TestForAllAnswer_SubmittedAt(t *testing.T) {
	collector := NewForAllCollector()
	collector.Start("Answer", 100)

	userID := uuid.New()
	beforeSubmit := time.Now()

	collector.SubmitAnswer(userID, "user1", "Answer")

	afterSubmit := time.Now()
	answers := collector.GetAllAnswers()

	if len(answers) == 0 {
		t.Fatal("No answers recorded")
	}

	submittedAt := answers[0].SubmittedAt
	if submittedAt.Before(beforeSubmit) || submittedAt.After(afterSubmit.Add(10*time.Millisecond)) {
		t.Error("SubmittedAt is not within expected range")
	}
}

