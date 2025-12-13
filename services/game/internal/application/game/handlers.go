package game

import (
	"time"

	"github.com/google/uuid"
	domainGame "sigame/game/internal/domain/game"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/domain/player"
	"sigame/game/internal/infrastructure/logger"
	wsMessage "sigame/game/internal/transport/ws/message"
)

func (m *Manager) findHost() uuid.UUID {
	for userID, p := range m.game.Players {
		if p.Role == player.RoleHost {
			return userID
		}
	}
	return uuid.Nil
}

func (m *Manager) selectActivePlayer() uuid.UUID {
	var selectedPlayer uuid.UUID
	minScore := MaxIntValue

	for userID, p := range m.game.Players {
		if p.Role == player.RoleHost {
			continue
		}
		if p.IsActive && p.Score < minScore {
			minScore = p.Score
			selectedPlayer = userID
		}
	}

	return selectedPlayer
}

func (m *Manager) handleSelectQuestion(action *PlayerAction) {
	if m.game.Status != domainGame.StatusQuestionSelect {
		logger.Warnf(m.ctx, "[SELECT_QUESTION] Invalid game status: %s, expected: %s", m.game.Status, domainGame.StatusQuestionSelect)
		return
	}

	p, ok := m.game.Players[action.UserID]
	if !ok {
		logger.Warnf(m.ctx, "[SELECT_QUESTION] Player not found: %s", action.UserID)
		return
	}
	if p.Role != player.RoleHost {
		logger.Warnf(m.ctx, "[SELECT_QUESTION] Player is not host: %s, role: %s", action.UserID, p.Role)
		return
	}

	payload := action.Message.GetPayload()
	themeIDRaw, ok := payload["theme_id"]
	if !ok {
		logger.Warnf(m.ctx, "[SELECT_QUESTION] Missing theme_id in payload: %v", payload)
		return
	}
	themeID, ok := themeIDRaw.(string)
	if !ok {
		logger.Warnf(m.ctx, "[SELECT_QUESTION] Invalid theme_id type: %T, value: %v", themeIDRaw, themeIDRaw)
		return
	}

	questionIDRaw, ok := payload["question_id"]
	if !ok {
		logger.Warnf(m.ctx, "[SELECT_QUESTION] Missing question_id in payload: %v", payload)
		return
	}
	questionID, ok := questionIDRaw.(string)
	if !ok {
		logger.Warnf(m.ctx, "[SELECT_QUESTION] Invalid question_id type: %T, value: %v", questionIDRaw, questionIDRaw)
		return
	}
	
	logger.Infof(m.ctx, "[SELECT_QUESTION] Processing: theme_id=%s, question_id=%s, user_id=%s", themeID, questionID, action.UserID)

	round := m.pack.GetRound(m.game.CurrentRound)
	if round == nil {
		return
	}

	var selectedTheme *pack.Theme
	var selectedQuestion *pack.Question

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

	if selectedQuestion == nil || selectedTheme == nil {
		return
	}

	m.selectQuestion(selectedTheme, selectedQuestion)
}

func (m *Manager) selectQuestion(theme *pack.Theme, question *pack.Question) {
	question.MarkAsUsed()
	m.game.SetCurrentQuestion(question, theme.Name)

	m.stakeInfo = nil
	m.secretTarget = nil

	questionType := question.GetType()
	logger.Infof(m.ctx, "[selectQuestion] Selected question: id=%s, type=%s, price=%d", question.ID, questionType, question.Price)
	switch questionType {
	case pack.TypeSecret:
		m.startSecretQuestion(question)
	case pack.TypeStake:
		m.startStakeQuestion(question)
	case pack.TypeForAll:
		m.startForAllQuestion(question)
	default:
		m.startNormalQuestion(question)
	}
}

func (m *Manager) startNormalQuestion(question *pack.Question) {
	m.game.UpdateStatus(domainGame.StatusQuestionShow)
	m.BroadcastState()

	if question.HasMedia() {
		m.sendStartMedia(question)
	}

	readTime := QuestionReadDuration
	if question.MediaDurationMs > 0 {
		readTime += time.Duration(question.MediaDurationMs) * time.Millisecond
	}
	m.timer.Start(readTime)
}

func (m *Manager) startSecretQuestion(question *pack.Question) {
	logger.Infof(m.ctx, "[startSecretQuestion] Starting secret question, price: %d", question.Price)
	m.game.UpdateStatus(domainGame.StatusSecretTransfer)
	m.BroadcastState()
	m.timer.Start(SecretTransferDuration)
	logger.Infof(m.ctx, "[startSecretQuestion] Status changed to: %s, timer started for %v", m.game.Status, SecretTransferDuration)
}

func (m *Manager) startStakeQuestion(question *pack.Question) {
	logger.Infof(m.ctx, "[startStakeQuestion] Starting stake question, price: %d", question.Price)
	activePlayerID := m.selectActivePlayer()
	m.game.SetActivePlayer(activePlayerID)

	activePlayer := m.game.Players[activePlayerID]

	minBet := question.Price
	maxBet := activePlayer.Score
	if maxBet < minBet {
		maxBet = minBet
	}

	m.stakeInfo = &domainGame.StakeInfo{
		MinBet:     minBet,
		MaxBet:     maxBet,
		CurrentBet: 0,
		IsAllIn:    false,
	}

	logger.Infof(m.ctx, "[startStakeQuestion] Active player: %s, minBet: %d, maxBet: %d", activePlayerID, minBet, maxBet)
	m.game.UpdateStatus(domainGame.StatusStakeBetting)
	m.BroadcastState()
	m.timer.Start(StakeBettingDuration)
	logger.Infof(m.ctx, "[startStakeQuestion] Status changed to: %s, timer started for %v", m.game.Status, StakeBettingDuration)
}

func (m *Manager) startForAllQuestion(question *pack.Question) {
	logger.Infof(m.ctx, "[startForAllQuestion] Starting forAll question, price: %d", question.Price)
	m.forAllCollector.Start(question.Answer, question.Price)

	m.game.UpdateStatus(domainGame.StatusQuestionShow)
	m.BroadcastState()

	if question.HasMedia() {
		m.sendStartMedia(question)
	}

	readTime := QuestionReadDuration
	if question.MediaDurationMs > 0 {
		readTime += time.Duration(question.MediaDurationMs) * time.Millisecond
	}

	logger.Infof(m.ctx, "[startForAllQuestion] Status changed to: %s, readTime: %v", m.game.Status, readTime)
	go func() {
		time.Sleep(readTime)
		m.mu.Lock()
		defer m.mu.Unlock()

		if m.game.Status == domainGame.StatusQuestionShow && m.game.CurrentQuestion != nil &&
			m.game.CurrentQuestion.GetType() == pack.TypeForAll {
			logger.Infof(m.ctx, "[startForAllQuestion] Transitioning to forAllAnswering after readTime")
			m.game.UpdateStatus(domainGame.StatusForAllAnswering)
			m.BroadcastState()
			m.timer.Start(time.Duration(m.game.Settings.TimeForAnswer) * time.Second)
			logger.Infof(m.ctx, "[startForAllQuestion] Status changed to: %s, timer started", m.game.Status)
		}
	}()
}

func (m *Manager) sendStartMedia(question *pack.Question) {
	if !question.HasMedia() {
		return
	}

	round := m.pack.GetRound(m.game.CurrentRound)
	if round == nil {
		logger.Errorf(nil, "%v", ErrRoundNotFound)
		return
	}

	var themeIndex int = -1
	if m.game.CurrentTheme != nil {
		for i, theme := range round.Themes {
			if theme.Name == *m.game.CurrentTheme {
				themeIndex = i
				break
			}
		}
	}

	if themeIndex == -1 {
		logger.Errorf(nil, "%v", ErrThemeNotFound)
		return
	}

	mediaItem := m.mediaTracker.FindMediaByQuestion(themeIndex, question.Price)
	if mediaItem == nil {
		logger.Errorf(nil, "%v", ErrMediaItemNotFound)
		return
	}

	now := time.Now().UnixMilli()
	durationMs := int64(question.MediaDurationMs)
	if durationMs == 0 {
		durationMs = DefaultMediaDurationMs
	}

	msg := wsMessage.NewStartMediaMessage(
		mediaItem.ID,
		mediaItem.Type,
		mediaItem.URL,
		now,
		durationMs,
	)

	data, err := msg.ToJSON()
	if err != nil {
		logger.Errorf(nil, "%v", ErrSerializeStartMediaMessage(err))
		return
	}

	m.hub.Broadcast(m.game.ID, data)
}

func (m *Manager) handlePressButton(userID uuid.UUID) {
	logger.Infof(m.ctx, "[PRESS_BUTTON] Received from user: %s, game status: %s", userID, m.game.Status)
	if m.game.Status != domainGame.StatusButtonPress {
		logger.Warnf(m.ctx, "[PRESS_BUTTON] Invalid game status: %s, expected: %s", m.game.Status, domainGame.StatusButtonPress)
		return
	}

	p, ok := m.game.Players[userID]
	if !ok {
		logger.Warnf(m.ctx, "[PRESS_BUTTON] Player not found: %s", userID)
		return
	}
	if !p.IsActive {
		logger.Warnf(m.ctx, "[PRESS_BUTTON] Player is not active: %s", userID)
		return
	}
	if p.Role == player.RoleHost {
		logger.Warnf(m.ctx, "[PRESS_BUTTON] Host cannot press button: %s, role: %s", userID, p.Role)
		return
	}
	
	logger.Infof(m.ctx, "[PRESS_BUTTON] Processing button press: user=%s, username=%s", userID, p.Username)

	rtt := m.hub.GetClientRTT(m.game.ID, userID)

	if m.buttonPress.Press(userID, p.Username, rtt) {
		if m.buttonPress.GetPressCount() == 1 {
			go m.finishButtonPressCollection()
		}
	}
}

func (m *Manager) finishButtonPressCollection() {
	logger.Infof(m.ctx, "[finishButtonPressCollection] Starting, waiting %v", ButtonPressCollectionWindow)
	time.Sleep(ButtonPressCollectionWindow)

	m.mu.Lock()
	defer m.mu.Unlock()

	logger.Infof(m.ctx, "[finishButtonPressCollection] After sleep, game status: %s", m.game.Status)
	if m.game.Status != domainGame.StatusButtonPress {
		logger.Warnf(m.ctx, "[finishButtonPressCollection] Game status changed, aborting: %s", m.game.Status)
		return
	}

	m.buttonPress.Close()
	winner := m.buttonPress.GetWinner()
	if winner == nil {
		logger.Warnf(m.ctx, "[finishButtonPressCollection] No winner found")
		return
	}

	logger.Infof(m.ctx, "[finishButtonPressCollection] Winner: %s (%s), setting active player and transitioning to answer_judging immediately", winner.UserID, winner.Username)
	m.timer.Stop()
	m.game.SetActivePlayer(winner.UserID)

	// Переходим сразу в answer_judging, чтобы кнопки судейства появились сразу
	// Таймер все равно работает для отсчета времени ответа
	m.game.UpdateStatus(domainGame.StatusAnswerJudging)
	logger.Infof(m.ctx, "[finishButtonPressCollection] Status changed to: %s, activePlayer: %v, broadcasting state", m.game.Status, m.game.ActivePlayer)
	m.BroadcastState()
	m.timer.Start(time.Duration(m.game.Settings.TimeForAnswer) * time.Second)
	logger.Infof(m.ctx, "[finishButtonPressCollection] Timer started for %d seconds (for answer timeout)", m.game.Settings.TimeForAnswer)
}

func (m *Manager) handleSubmitAnswer(action *PlayerAction) {
	if m.game.Status != domainGame.StatusAnswering {
		return
	}

	if m.game.ActivePlayer == nil || *m.game.ActivePlayer != action.UserID {
		return
	}

	answerStr, ok := action.Message.GetPayload()["answer"].(string)
	if !ok {
		return
	}

	p := m.game.Players[action.UserID]
	m.timer.Stop()

	correct := m.game.CurrentQuestion.ValidateAnswer(answerStr)

	questionPrice := m.game.CurrentQuestion.Price
	if correct {
		p.AddScore(questionPrice)
	} else {
		p.SubtractScore(questionPrice)
	}

	m.transitionToAnswerJudging()
}

func (m *Manager) handleJudgeAnswer(action *PlayerAction) {
	hostPlayer := m.game.Players[action.UserID]
	if hostPlayer.Role != player.RoleHost {
		return
	}

	if m.game.Status != domainGame.StatusAnswerJudging {
		return
	}

	correct, ok := action.Message.GetPayload()["correct"].(bool)
	if !ok {
		return
	}

	answeringUserIDStr, ok := action.Message.GetPayload()["user_id"].(string)
	if !ok {
		return
	}

	answeringUserID, err := uuid.Parse(answeringUserIDStr)
	if err != nil {
		return
	}

	p := m.game.Players[answeringUserID]
	questionPrice := m.game.CurrentQuestion.Price

	if correct {
		p.AddScore(questionPrice)
	} else {
		p.SubtractScore(questionPrice)
	}

	m.continueGame()
}

func (m *Manager) handleMediaLoadProgress(action *PlayerAction) {
	loaded, _ := action.Message.GetPayload()["loaded"].(float64)
	total, _ := action.Message.GetPayload()["total"].(float64)
	bytesLoaded, _ := action.Message.GetPayload()["bytes_loaded"].(float64)
	percent, _ := action.Message.GetPayload()["percent"].(float64)

	m.mediaTracker.UpdateProgress(
		action.UserID,
		int(loaded),
		int(total),
		int64(bytesLoaded),
		int(percent),
	)
}

func (m *Manager) handleMediaLoadComplete(action *PlayerAction) {
	loadedCount, _ := action.Message.GetPayload()["loaded_count"].(float64)
	m.mediaTracker.MarkComplete(action.UserID, int(loadedCount))
}

func (m *Manager) handleTransferSecret(action *PlayerAction) {
	logger.Infof(m.ctx, "[TRANSFER_SECRET] Received from user: %s, game status: %s", action.UserID, m.game.Status)
	if m.game.Status != domainGame.StatusSecretTransfer {
		logger.Warnf(m.ctx, "[TRANSFER_SECRET] Invalid game status: %s, expected: %s", m.game.Status, domainGame.StatusSecretTransfer)
		return
	}

	hostPlayer := m.game.Players[action.UserID]
	if hostPlayer.Role != player.RoleHost {
		logger.Warnf(m.ctx, "[TRANSFER_SECRET] User is not host: %s, role: %s", action.UserID, hostPlayer.Role)
		return
	}

	targetUserIDStr, ok := action.Message.GetPayload()["target_user_id"].(string)
	if !ok {
		return
	}

	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return
	}

	targetPlayer, exists := m.game.Players[targetUserID]
	if !exists || targetPlayer.Role == player.RoleHost {
		return
	}

	m.transferSecretToPlayer(action.UserID, targetUserID)
}

func (m *Manager) transferSecretToPlayer(fromUserID, toUserID uuid.UUID) {
	m.timer.Stop()
	m.game.SetActivePlayer(toUserID)
	m.secretTarget = &toUserID

	m.game.UpdateStatus(domainGame.StatusQuestionShow)
	m.BroadcastState()

	if m.game.CurrentQuestion.HasMedia() {
		m.sendStartMedia(m.game.CurrentQuestion)
	}

	readTime := QuestionReadDuration
	if m.game.CurrentQuestion.MediaDurationMs > 0 {
		readTime += time.Duration(m.game.CurrentQuestion.MediaDurationMs) * time.Millisecond
	}
	m.timer.Start(readTime)
}

func (m *Manager) handlePlaceStake(action *PlayerAction) {
	logger.Infof(m.ctx, "[PLACE_STAKE] Received from user: %s, game status: %s", action.UserID, m.game.Status)
	if m.game.Status != domainGame.StatusStakeBetting {
		logger.Warnf(m.ctx, "[PLACE_STAKE] Invalid game status: %s, expected: %s", m.game.Status, domainGame.StatusStakeBetting)
		return
	}

	if m.game.ActivePlayer == nil || *m.game.ActivePlayer != action.UserID {
		logger.Warnf(m.ctx, "[PLACE_STAKE] User is not active player: %s, active: %v", action.UserID, m.game.ActivePlayer)
		return
	}

	amount, ok := action.Message.GetPayload()["amount"].(float64)
	if !ok {
		return
	}

	allIn, _ := action.Message.GetPayload()["all_in"].(bool)

	m.placeStake(action.UserID, int(amount), allIn)
}

func (m *Manager) placeStake(userID uuid.UUID, amount int, allIn bool) {
	m.timer.Stop()

	p := m.game.Players[userID]

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
		amount = p.Score
		if amount < m.stakeInfo.MinBet {
			amount = m.stakeInfo.MinBet
		}
	}

	m.stakeInfo.CurrentBet = amount
	m.stakeInfo.IsAllIn = allIn
	m.game.CurrentQuestion.Price = amount

	m.game.UpdateStatus(domainGame.StatusQuestionShow)
	m.BroadcastState()

	if m.game.CurrentQuestion.HasMedia() {
		m.sendStartMedia(m.game.CurrentQuestion)
	}

	readTime := QuestionReadDuration
	if m.game.CurrentQuestion.MediaDurationMs > 0 {
		readTime += time.Duration(m.game.CurrentQuestion.MediaDurationMs) * time.Millisecond
	}
	m.timer.Start(readTime)
}

func (m *Manager) handleSubmitForAllAnswer(action *PlayerAction) {
	if m.game.Status != domainGame.StatusForAllAnswering {
		return
	}

	p, ok := m.game.Players[action.UserID]
	if !ok || !p.IsActive || p.Role == player.RoleHost {
		return
	}

	answerStr, ok := action.Message.GetPayload()["answer"].(string)
	if !ok {
		return
	}

	if m.forAllCollector.SubmitAnswer(action.UserID, p.Username, answerStr) {
		expectedAnswers := 0
		for _, p := range m.game.Players {
			if p.Role != player.RoleHost && p.IsActive {
				expectedAnswers++
			}
		}

		if m.forAllCollector.GetAnswerCount() >= expectedAnswers {
			m.timer.Stop()
			m.finishForAllQuestion()
		}
	}
}

