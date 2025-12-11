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

// selectQuestion selects a question and routes based on question type
func (m *Manager) selectQuestion(theme *domain.Theme, question *domain.Question) {
	log.Printf("Question selected: %s - %s (type: %s)", theme.Name, question.ID, question.GetType())

	// Mark question as used
	question.MarkAsUsed()

	// Set current question
	m.game.CurrentQuestion = question
	m.game.CurrentTheme = &theme.Name

	// Clear special type state
	m.stakeInfo = nil
	m.secretTarget = nil

	// Log event
	event := domain.NewGameEvent(m.game.ID, domain.EventQuestionSelected).
		WithRound(m.game.CurrentRound).
		WithQuestion(question.ID).
		WithData("theme", theme.Name).
		WithData("price", question.Price).
		WithData("type", string(question.GetType()))
	m.eventLogger.LogEvent(m.ctx, event)

	// Route based on question type
	questionType := question.GetType()
	switch questionType {
	case domain.QuestionTypeSecret:
		m.startSecretQuestion(question)
	case domain.QuestionTypeStake:
		m.startStakeQuestion(question)
	case domain.QuestionTypeForAll:
		m.startForAllQuestion(question)
	default:
		m.startNormalQuestion(question)
	}
}

// startNormalQuestion handles normal question flow
func (m *Manager) startNormalQuestion(question *domain.Question) {
	// Show question
	m.updateGameStatus(domain.GameStatusQuestionShow)
	m.BroadcastState()

	// If question has media, send START_MEDIA for synchronized playback
	if question.MediaType != "" && question.MediaType != "text" && question.MediaURL != "" {
		m.sendStartMedia(question)
	}

	// Wait a bit for question to be read, then handleTimeout will transition to button_press
	// Add media duration to read time if present
	readTime := 3 * time.Second
	if question.MediaDurationMs > 0 {
		readTime += time.Duration(question.MediaDurationMs) * time.Millisecond
	}
	m.timer.Start(readTime)
}

// startSecretQuestion handles secret (Кот в мешке) question
func (m *Manager) startSecretQuestion(question *domain.Question) {
	log.Printf("Starting secret question - host must transfer to another player")

	// Show that it's a secret question (but don't show the content yet)
	m.updateGameStatus(domain.GameStatusSecretTransfer)
	m.BroadcastState()

	// Give host 30 seconds to choose a player
	m.timer.Start(30 * time.Second)
}

// startStakeQuestion handles stake (Ва-банк) question
func (m *Manager) startStakeQuestion(question *domain.Question) {
	log.Printf("Starting stake question - active player must place bet")

	// The player who selected this question makes the stake
	// Find who selected (the one with lowest score who's not host)
	activePlayerID := m.selectActivePlayer()
	m.game.ActivePlayer = &activePlayerID

	activePlayer := m.game.Players[activePlayerID]

	// Calculate min/max bet
	minBet := question.Price
	maxBet := activePlayer.Score
	if maxBet < minBet {
		maxBet = minBet // Can always bet at least the question price
	}

	m.stakeInfo = &domain.StakeInfo{
		MinBet:     minBet,
		MaxBet:     maxBet,
		CurrentBet: 0,
		IsAllIn:    false,
	}

	m.updateGameStatus(domain.GameStatusStakeBetting)
	m.BroadcastState()

	// Give player 20 seconds to place bet
	m.timer.Start(20 * time.Second)
}

// startForAllQuestion handles forAll question
func (m *Manager) startForAllQuestion(question *domain.Question) {
	log.Printf("Starting forAll question - all players will answer")

	// Initialize collector
	m.forAllCollector.Start(question.Answer, question.Price)

	// Show question
	m.updateGameStatus(domain.GameStatusQuestionShow)
	m.BroadcastState()

	// If question has media, send START_MEDIA for synchronized playback
	if question.MediaType != "" && question.MediaType != "text" && question.MediaURL != "" {
		m.sendStartMedia(question)
	}

	// After showing question, transition to answering phase
	readTime := 3 * time.Second
	if question.MediaDurationMs > 0 {
		readTime += time.Duration(question.MediaDurationMs) * time.Millisecond
	}

	// We'll transition in a goroutine after the read time
	go func() {
		time.Sleep(readTime)
		m.mu.Lock()
		defer m.mu.Unlock()

		if m.game.Status == domain.GameStatusQuestionShow && m.game.CurrentQuestion != nil &&
			m.game.CurrentQuestion.GetType() == domain.QuestionTypeForAll {
			m.updateGameStatus(domain.GameStatusForAllAnswering)
			m.BroadcastState()
			m.timer.Start(time.Duration(m.game.Settings.TimeForAnswer) * time.Second)
		}
	}()
}

// sendStartMedia sends START_MEDIA command for synchronized playback
func (m *Manager) sendStartMedia(question *domain.Question) {
	// Start playback 300ms from now to allow for network delay
	startAt := time.Now().Add(300 * time.Millisecond).UnixMilli()

	mediaID := question.ID + "_media"
	msg := websocket.NewStartMediaMessage(
		mediaID,
		question.MediaType,
		question.MediaURL,
		startAt,
		int64(question.MediaDurationMs),
	)

	if data, err := msg.ToJSON(); err == nil {
		m.hub.Broadcast(m.game.ID, data)
		log.Printf("Sent START_MEDIA for %s, start_at: %d", mediaID, startAt)
	}
}

// handlePressButton handles button press with RTT compensation
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

	// Get client's RTT for compensation
	rtt := m.hub.GetClientRTT(m.game.ID, userID)

	// Record press with RTT compensation
	if m.buttonPress.Press(userID, player.Username, rtt) {
		log.Printf("Button pressed by %s (RTT: %v)", player.Username, rtt)

		// Check if this is the first press - start collection window
		if m.buttonPress.GetPressCount() == 1 {
			// Start a short collection window (150ms) to gather all near-simultaneous presses
			// After this window, we determine the winner
			go m.finishButtonPressCollection()
		}
	}
}

// finishButtonPressCollection waits for collection window and determines winner
func (m *Manager) finishButtonPressCollection() {
	// Wait for collection window (150ms to account for network jitter)
	time.Sleep(150 * time.Millisecond)

	m.mu.Lock()
	defer m.mu.Unlock()

	// Make sure we're still in button press phase
	if m.game.Status != domain.GameStatusButtonPress {
		return
	}

	// Close collection and determine winner
	m.buttonPress.Close()

	winner := m.buttonPress.GetWinner()
	if winner == nil {
		// No presses recorded (shouldn't happen if we got here)
		return
	}

	// Stop main timer
	m.timer.Stop()

	// Set winner as active player
	m.game.ActivePlayer = &winner.UserID

	// Log event
	event := domain.NewGameEvent(m.game.ID, domain.EventButtonPressed).
		WithUser(winner.UserID).
		WithQuestion(m.game.CurrentQuestion.ID).
		WithData("rtt_ms", winner.RTT.Milliseconds()).
		WithData("reaction_time_ms", m.buttonPress.GetReactionTime(winner))
	m.eventLogger.LogEvent(m.ctx, event)

	// Build all_presses for response
	allPresses := m.buttonPress.GetAllPresses()
	pressInfos := make([]websocket.PressInfo, len(allPresses))
	for i, entry := range allPresses {
		pressInfos[i] = websocket.PressInfo{
			UserID:   entry.UserID,
			Username: entry.Username,
			TimeMS:   m.buttonPress.GetReactionTime(&entry),
		}
	}

	// Broadcast button press result with all presses
	reactionTime := m.buttonPress.GetReactionTime(winner)
	msg := websocket.NewButtonPressedMessage(winner.UserID, winner.Username, reactionTime, pressInfos)
	if data, err := msg.ToJSON(); err == nil {
		m.hub.Broadcast(m.game.ID, data)
	}

	log.Printf("Button winner: %s (reaction: %dms, total presses: %d)",
		winner.Username, reactionTime, len(allPresses))

	// Transition to answer judging phase (player answers verbally, host judges)
	m.updateGameStatus(domain.GameStatusAnswerJudging)
	m.BroadcastState()

	// Start timer for host to judge the answer
	m.timer.Start(30 * time.Second)
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

// skipQuestion skips the current question (no one pressed button or timeout during collection)
func (m *Manager) skipQuestion() {
	// Close button press collection if still open
	m.buttonPress.Close()

	// Check if there were any presses during timeout
	if m.buttonPress.HasPresses() {
		// Someone pressed but we timed out - determine winner anyway
		winner := m.buttonPress.GetWinner()
		if winner != nil {
			log.Printf("Question timeout but had presses, winner: %s", winner.Username)
			m.game.ActivePlayer = &winner.UserID

			// Build all_presses for response
			allPresses := m.buttonPress.GetAllPresses()
			pressInfos := make([]websocket.PressInfo, len(allPresses))
			for i, entry := range allPresses {
				pressInfos[i] = websocket.PressInfo{
					UserID:   entry.UserID,
					Username: entry.Username,
					TimeMS:   m.buttonPress.GetReactionTime(&entry),
				}
			}

			// Broadcast button press result
			reactionTime := m.buttonPress.GetReactionTime(winner)
			msg := websocket.NewButtonPressedMessage(winner.UserID, winner.Username, reactionTime, pressInfos)
			if data, err := msg.ToJSON(); err == nil {
				m.hub.Broadcast(m.game.ID, data)
			}

			// Transition to answer judging
			m.updateGameStatus(domain.GameStatusAnswerJudging)
			m.BroadcastState()
			m.timer.Start(30 * time.Second)
			return
		}
	}

	log.Printf("Question skipped (no button press)")
	m.continueGame()
}

// continueGame continues to next question or ends round
func (m *Manager) continueGame() {
	// Clear current question and special type state
	m.game.CurrentQuestion = nil
	m.game.CurrentTheme = nil
	m.stakeInfo = nil
	m.secretTarget = nil
	m.forAllCollector.Reset()

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

// handleMediaLoadProgress handles client's media loading progress
func (m *Manager) handleMediaLoadProgress(action *PlayerAction) {
	// Extract progress data from payload
	loaded, _ := action.Message.Payload["loaded"].(float64)
	total, _ := action.Message.Payload["total"].(float64)
	bytesLoaded, _ := action.Message.Payload["bytes_loaded"].(float64)
	percent, _ := action.Message.Payload["percent"].(float64)

	m.mediaTracker.UpdateProgress(
		action.UserID,
		int(loaded),
		int(total),
		int64(bytesLoaded),
		int(percent),
	)

	log.Printf("Media load progress from %s: %d%% (%d/%d files)",
		action.UserID, int(percent), int(loaded), int(total))
}

// handleMediaLoadComplete handles client finishing media loading
func (m *Manager) handleMediaLoadComplete(action *PlayerAction) {
	// Extract data from payload
	loadedCount, _ := action.Message.Payload["loaded_count"].(float64)

	m.mediaTracker.MarkComplete(action.UserID, int(loadedCount))

	log.Printf("Media load complete from %s: %d files loaded", action.UserID, int(loadedCount))

	// Check if all clients are ready
	if m.mediaTracker.AllClientsReady() {
		log.Printf("All clients have loaded media for round %d", m.game.CurrentRound)
	}
}

// ==================== Special Question Type Handlers ====================

// handleTransferSecret handles host transferring secret question to another player
func (m *Manager) handleTransferSecret(action *PlayerAction) {
	if m.game.Status != domain.GameStatusSecretTransfer {
		return
	}

	// Only host can transfer
	hostPlayer := m.game.Players[action.UserID]
	if hostPlayer.Role != domain.PlayerRoleHost {
		return
	}

	// Get target player ID
	targetUserIDStr, ok := action.Message.Payload["target_user_id"].(string)
	if !ok {
		return
	}

	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return
	}

	// Verify target is a valid player (not host)
	targetPlayer, exists := m.game.Players[targetUserID]
	if !exists || targetPlayer.Role == domain.PlayerRoleHost {
		return
	}

	m.transferSecretToPlayer(action.UserID, targetUserID)
}

// transferSecretToPlayer executes the secret transfer
func (m *Manager) transferSecretToPlayer(fromUserID, toUserID uuid.UUID) {
	m.timer.Stop()

	fromPlayer := m.game.Players[fromUserID]
	toPlayer := m.game.Players[toUserID]

	log.Printf("Secret question transferred from %s to %s", fromPlayer.Username, toPlayer.Username)

	// Set target as active player
	m.game.ActivePlayer = &toUserID
	m.secretTarget = &toUserID

	// Broadcast transfer notification
	msg := websocket.NewSecretTransferredMessage(fromUserID, fromPlayer.Username, toUserID, toPlayer.Username)
	if data, err := msg.ToJSON(); err == nil {
		m.hub.Broadcast(m.game.ID, data)
	}

	// Now show the question content
	m.updateGameStatus(domain.GameStatusQuestionShow)
	m.BroadcastState()

	// If question has media, send START_MEDIA
	if m.game.CurrentQuestion.MediaType != "" && m.game.CurrentQuestion.MediaType != "text" {
		m.sendStartMedia(m.game.CurrentQuestion)
	}

	// Wait for question to be read, then go to answer judging (target must answer)
	readTime := 3 * time.Second
	if m.game.CurrentQuestion.MediaDurationMs > 0 {
		readTime += time.Duration(m.game.CurrentQuestion.MediaDurationMs) * time.Millisecond
	}
	m.timer.Start(readTime)
}

// handleSecretTransferTimeout auto-selects first available player for secret transfer
func (m *Manager) handleSecretTransferTimeout() {
	log.Printf("Secret transfer timeout - auto-selecting player")

	hostID := m.findHost()

	// Find first non-host player
	for userID, player := range m.game.Players {
		if player.Role != domain.PlayerRoleHost {
			m.transferSecretToPlayer(hostID, userID)
			return
		}
	}

	// No players found (shouldn't happen) - skip question
	m.continueGame()
}

// handlePlaceStake handles player placing a stake
func (m *Manager) handlePlaceStake(action *PlayerAction) {
	if m.game.Status != domain.GameStatusStakeBetting {
		return
	}

	// Only active player can place stake
	if m.game.ActivePlayer == nil || *m.game.ActivePlayer != action.UserID {
		return
	}

	// Get stake amount
	amount, ok := action.Message.Payload["amount"].(float64)
	if !ok {
		return
	}

	allIn, _ := action.Message.Payload["all_in"].(bool)

	m.placeStake(action.UserID, int(amount), allIn)
}

// placeStake executes the stake placement
func (m *Manager) placeStake(userID uuid.UUID, amount int, allIn bool) {
	m.timer.Stop()

	player := m.game.Players[userID]

	// Validate amount
	if m.stakeInfo == nil {
		return
	}

	if amount < m.stakeInfo.MinBet {
		amount = m.stakeInfo.MinBet
	}
	if amount > m.stakeInfo.MaxBet {
		amount = m.stakeInfo.MaxBet
	}
	if allIn {
		amount = player.Score
		if amount < m.stakeInfo.MinBet {
			amount = m.stakeInfo.MinBet
		}
	}

	m.stakeInfo.CurrentBet = amount
	m.stakeInfo.IsAllIn = allIn

	log.Printf("Player %s placed stake: %d (all-in: %v)", player.Username, amount, allIn)

	// Broadcast stake notification
	msg := websocket.NewStakePlacedMessage(userID, player.Username, amount, allIn)
	if data, err := msg.ToJSON(); err == nil {
		m.hub.Broadcast(m.game.ID, data)
	}

	// Update question price to the stake amount for scoring
	m.game.CurrentQuestion.Price = amount

	// Now show the question
	m.updateGameStatus(domain.GameStatusQuestionShow)
	m.BroadcastState()

	// If question has media, send START_MEDIA
	if m.game.CurrentQuestion.MediaType != "" && m.game.CurrentQuestion.MediaType != "text" {
		m.sendStartMedia(m.game.CurrentQuestion)
	}

	// Wait for question to be read, then go to answer judging
	readTime := 3 * time.Second
	if m.game.CurrentQuestion.MediaDurationMs > 0 {
		readTime += time.Duration(m.game.CurrentQuestion.MediaDurationMs) * time.Millisecond
	}
	m.timer.Start(readTime)
}

// handleStakeBettingTimeout uses minimum bet when player doesn't respond
func (m *Manager) handleStakeBettingTimeout() {
	if m.game.ActivePlayer == nil || m.stakeInfo == nil {
		m.continueGame()
		return
	}

	log.Printf("Stake betting timeout - using minimum bet")
	m.placeStake(*m.game.ActivePlayer, m.stakeInfo.MinBet, false)
}

// handleSubmitForAllAnswer handles player submitting answer for forAll question
func (m *Manager) handleSubmitForAllAnswer(action *PlayerAction) {
	if m.game.Status != domain.GameStatusForAllAnswering {
		return
	}

	player, ok := m.game.Players[action.UserID]
	if !ok || !player.IsActive {
		return
	}

	// Host doesn't participate
	if player.Role == domain.PlayerRoleHost {
		return
	}

	// Get answer
	answerStr, ok := action.Message.Payload["answer"].(string)
	if !ok {
		return
	}

	// Submit to collector
	if m.forAllCollector.SubmitAnswer(action.UserID, player.Username, answerStr) {
		log.Printf("ForAll answer from %s: %s", player.Username, answerStr)
	}

	// Check if all players have answered
	expectedAnswers := 0
	for _, p := range m.game.Players {
		if p.Role != domain.PlayerRoleHost && p.IsActive {
			expectedAnswers++
		}
	}

	if m.forAllCollector.GetAnswerCount() >= expectedAnswers {
		// All answers received, finish early
		m.timer.Stop()
		m.finishForAllQuestion()
	}
}

// finishForAllQuestion processes all answers and shows results
func (m *Manager) finishForAllQuestion() {
	log.Printf("Finishing forAll question")

	m.forAllCollector.Close()

	// Get results
	results := m.forAllCollector.GetResults(func(userAnswer, correctAnswer string) bool {
		return m.game.CurrentQuestion.ValidateAnswer(userAnswer)
	})

	// Apply scores and build result list
	resultList := make([]domain.ForAllAnswerResult, 0, len(results))
	for userID, result := range results {
		player := m.game.Players[userID]
		if result.IsCorrect {
			player.AddScore(result.ScoreDelta)
		} else {
			player.SubtractScore(-result.ScoreDelta) // ScoreDelta is negative for wrong answers
		}

		resultList = append(resultList, domain.ForAllAnswerResult{
			UserID:     result.UserID,
			Username:   result.Username,
			Answer:     result.Answer,
			IsCorrect:  result.IsCorrect,
			ScoreDelta: result.ScoreDelta,
		})

		log.Printf("ForAll result for %s: correct=%v, delta=%d", result.Username, result.IsCorrect, result.ScoreDelta)
	}

	// Broadcast results
	msg := websocket.NewForAllResultsMessage(m.forAllCollector.GetCorrectAnswer(), resultList)
	if data, err := msg.ToJSON(); err == nil {
		m.hub.Broadcast(m.game.ID, data)
	}

	// Show results screen
	m.updateGameStatus(domain.GameStatusForAllResults)
	m.BroadcastState()

	// Wait 5 seconds to show results, then continue
	m.timer.Start(5 * time.Second)
}

