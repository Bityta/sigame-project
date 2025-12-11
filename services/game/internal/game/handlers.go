package game

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
	"github.com/sigame/game/internal/transport/websocket"
)

// handleSelectQuestion handles question selection
func (m *Manager) handleSelectQuestion(action *PlayerAction) {
	if m.game.Status != domain.GameStatusQuestionSelect {
		return
	}

	// Only host can select questions
	player, ok := m.game.Players[action.UserID]
	if !ok || player.Role != domain.PlayerRoleHost {
		return
	}

	// Parse payload
	payload, ok := action.Message.Payload["theme_id"].(string)
	if !ok {
		return
	}
	themeID := payload

	questionID, ok := action.Message.Payload["question_id"].(string)
	if !ok {
		return
	}

	// Find the question
	round := m.pack.Rounds[m.game.CurrentRound-1]
	var selectedTheme *domain.Theme
	var selectedQuestion *domain.Question

	for _, theme := range round.Themes {
		if theme.ID == themeID || theme.Name == themeID {
			selectedTheme = theme
			for _, q := range theme.Questions {
				if q.ID == questionID && q.IsAvailable() {
					selectedQuestion = q
					break
				}
			}
			break
		}
	}

	if selectedQuestion == nil {
		log.Printf("Question not found or not available")
		return
	}

	m.selectQuestion(selectedTheme, selectedQuestion)
}

// selectQuestion selects a question and shows it
func (m *Manager) selectQuestion(theme *domain.Theme, question *domain.Question) {
	log.Printf("Question selected: %s - %s", theme.Name, question.ID)

	// Mark question as used
	question.MarkAsUsed()

	// Set current question
	m.game.CurrentQuestion = question
	m.game.CurrentTheme = &theme.Name

	// Log event
	event := domain.NewGameEvent(m.game.ID, domain.EventQuestionSelected).
		WithRound(m.game.CurrentRound).
		WithQuestion(question.ID).
		WithData("theme", theme.Name).
		WithData("price", question.Price)
	m.eventLogger.LogEvent(m.ctx, event)

	// Show question
	m.updateGameStatus(domain.GameStatusQuestionShow)
	m.BroadcastState()

	// Wait a bit for question to be read, then handleTimeout will transition to button_press
	m.timer.Start(3 * time.Second) // 3 seconds to read
}

// handlePressButton handles button press
func (m *Manager) handlePressButton(userID uuid.UUID) {
	if m.game.Status != domain.GameStatusButtonPress {
		return
	}

	player, ok := m.game.Players[userID]
	if !ok || !player.IsActive {
		return
	}

	// Host cannot press button - they are the game master
	if player.Role == domain.PlayerRoleHost {
		return
	}

	// Try to press button (atomic)
	if m.buttonPress.Press(userID) {
		log.Printf("Button pressed by %s", player.Username)

		// This player gets to answer
		m.game.ActivePlayer = &userID

		// Stop timer
		m.timer.Stop()

		// Log event
		event := domain.NewGameEvent(m.game.ID, domain.EventButtonPressed).
			WithUser(userID).
			WithQuestion(m.game.CurrentQuestion.ID)
		m.eventLogger.LogEvent(m.ctx, event)

		// Broadcast button press
		latency := m.buttonPress.GetLatency()
		msg := websocket.NewButtonPressedMessage(userID, player.Username, latency)
		if data, err := msg.ToJSON(); err == nil {
			m.hub.Broadcast(m.game.ID, data)
		}

		// Transition to answer judging phase (player answers verbally, host judges)
		m.updateGameStatus(domain.GameStatusAnswerJudging)
		m.BroadcastState()

		// Start timer for host to judge the answer
		m.timer.Start(30 * time.Second)
	}
}

// handleSubmitAnswer handles answer submission
func (m *Manager) handleSubmitAnswer(action *PlayerAction) {
	if m.game.Status != domain.GameStatusAnswering {
		return
	}

	// Verify it's the active player
	if m.game.ActivePlayer == nil || *m.game.ActivePlayer != action.UserID {
		return
	}

	// Get answer
	answerStr, ok := action.Message.Payload["answer"].(string)
	if !ok {
		return
	}

	player := m.game.Players[action.UserID]

	// Stop timer
	m.timer.Stop()

	// Validate answer
	correct := m.game.CurrentQuestion.ValidateAnswer(answerStr)

	// Log event
	event := domain.NewGameEvent(m.game.ID, domain.EventAnswerSubmitted).
		WithUser(action.UserID).
		WithQuestion(m.game.CurrentQuestion.ID).
		WithData("answer", answerStr).
		WithData("correct", correct)
	m.eventLogger.LogEvent(m.ctx, event)

	// Update score
	questionPrice := m.game.CurrentQuestion.Price
	scoreDelta := 0

	if correct {
		player.AddScore(questionPrice)
		scoreDelta = questionPrice
	} else {
		player.SubtractScore(questionPrice)
		scoreDelta = -questionPrice
	}

	log.Printf("Player %s answered %v. Score: %d (delta: %d)", player.Username, correct, player.Score, scoreDelta)

	// Broadcast answer result
	msg := websocket.NewAnswerResultMessage(
		action.UserID,
		player.Username,
		correct,
		m.game.CurrentQuestion.Answer,
		player.Score,
		scoreDelta,
	)
	if data, err := msg.ToJSON(); err == nil {
		m.hub.Broadcast(m.game.ID, data)
	}

	// Continue to next question
	m.continueGame()
}

// handleJudgeAnswer handles manual answer judging by host
func (m *Manager) handleJudgeAnswer(action *PlayerAction) {
	// Only host can judge
	hostPlayer := m.game.Players[action.UserID]
	if hostPlayer.Role != domain.PlayerRoleHost {
		return
	}

	if m.game.Status != domain.GameStatusAnswerJudging {
		return
	}

	// Get judging result
	correct, ok := action.Message.Payload["correct"].(bool)
	if !ok {
		return
	}

	answeringUserIDStr, ok := action.Message.Payload["user_id"].(string)
	if !ok {
		return
	}

	answeringUserID, err := uuid.Parse(answeringUserIDStr)
	if err != nil {
		return
	}

	player := m.game.Players[answeringUserID]

	// Update score based on host's judgment
	questionPrice := m.game.CurrentQuestion.Price
	scoreDelta := 0

	if correct {
		player.AddScore(questionPrice)
		scoreDelta = questionPrice
	} else {
		player.SubtractScore(questionPrice)
		scoreDelta = -questionPrice
	}

	// Broadcast result
	msg := websocket.NewAnswerResultMessage(
		answeringUserID,
		player.Username,
		correct,
		m.game.CurrentQuestion.Answer,
		player.Score,
		scoreDelta,
	)
	if data, err := msg.ToJSON(); err == nil {
		m.hub.Broadcast(m.game.ID, data)
	}

	// Continue game
	m.continueGame()
}

// handleAnswerTimeout handles timeout during answering phase
func (m *Manager) handleAnswerTimeout() {
	if m.game.ActivePlayer == nil {
		return
	}

	player := m.game.Players[*m.game.ActivePlayer]

	// Penalize for timeout
	questionPrice := m.game.CurrentQuestion.Price
	player.SubtractScore(questionPrice)

	log.Printf("Answer timeout for player %s", player.Username)

	// Broadcast timeout
	msg := websocket.NewAnswerResultMessage(
		*m.game.ActivePlayer,
		player.Username,
		false,
		m.game.CurrentQuestion.Answer,
		player.Score,
		-m.game.CurrentQuestion.Price,
	)
	if data, err := msg.ToJSON(); err == nil {
		m.hub.Broadcast(m.game.ID, data)
	}

	m.continueGame()
}

// skipQuestion skips the current question (no one pressed button)
func (m *Manager) skipQuestion() {
	log.Printf("Question skipped (no button press)")
	m.continueGame()
}

// continueGame continues to next question or ends round
func (m *Manager) continueGame() {
	// Clear current question
	m.game.CurrentQuestion = nil
	m.game.CurrentTheme = nil

	// Check if round is complete
	round := m.pack.Rounds[m.game.CurrentRound-1]
	if round.IsComplete() {
		m.endRound()
		return
	}

	// Host selects next question (they are the game master)
	hostID := m.findHost()
	m.game.ActivePlayer = &hostID

	// Back to question selection
	m.updateGameStatus(domain.GameStatusQuestionSelect)
	m.BroadcastState()

	// Start timer for next selection
	m.timer.Start(time.Duration(m.game.Settings.TimeForChoice) * time.Second)
}

// endRound ends the current round
func (m *Manager) endRound() {
	log.Printf("Round %d ended", m.game.CurrentRound)

	m.updateGameStatus(domain.GameStatusRoundEnd)

	// Log event
	event := domain.NewGameEvent(m.game.ID, domain.EventRoundFinished).
		WithRound(m.game.CurrentRound)
	m.eventLogger.LogEvent(m.ctx, event)

	// Broadcast round end
	m.BroadcastState()

	// Schedule next round/game end in a goroutine to avoid deadlock
	currentRound := m.game.CurrentRound
	totalRounds := len(m.pack.Rounds)
	
	go func() {
		time.Sleep(5 * time.Second)
		
		m.mu.Lock()
		defer m.mu.Unlock()

		if currentRound < totalRounds {
			m.startRound(currentRound + 1)
		} else {
			m.endGame()
		}
	}()
}

// endGame ends the game
func (m *Manager) endGame() {
	log.Printf("Game %s ended", m.game.ID)

	m.updateGameStatus(domain.GameStatusGameEnd)
	now := time.Now()
	m.game.FinishedAt = &now

	// Calculate winners
	m.game.Winners = m.calculateWinners()
	m.game.FinalScores = m.calculateFinalScores()

	// Log event
	event := domain.NewGameEvent(m.game.ID, domain.EventGameFinished)
	m.eventLogger.LogEvent(m.ctx, event)

	// Broadcast game end with winners
	m.BroadcastState()

	// Stop timer ticker as game is done
	if m.timerTicker != nil {
		m.timerTicker.Stop()
	}
}

// calculateWinners calculates the game winners (top 3 players)
func (m *Manager) calculateWinners() []domain.PlayerScore {
	scores := m.calculateFinalScores()
	
	// Top 3 are winners
	winners := make([]domain.PlayerScore, 0)
	for i, score := range scores {
		if i >= 3 {
			break
		}
		winners = append(winners, score)
	}
	
	return winners
}

// calculateFinalScores calculates and ranks all players by score
func (m *Manager) calculateFinalScores() []domain.PlayerScore {
	scores := make([]domain.PlayerScore, 0)
	
	// Collect scores (exclude host)
	for userID, player := range m.game.Players {
		if player.Role == domain.PlayerRoleHost {
			continue
		}
		scores = append(scores, domain.PlayerScore{
			UserID:   userID,
			Username: player.Username,
			Score:    player.Score,
		})
	}
	
	// Sort by score descending
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].Score > scores[i].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}
	
	// Assign ranks
	for i := range scores {
		scores[i].Rank = i + 1
	}
	
	return scores
}

