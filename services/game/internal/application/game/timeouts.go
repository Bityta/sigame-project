package game

import (
	"time"

	domainGame "sigame/game/internal/domain/game"
)

func (m *Manager) handleTimeout() {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.game.Status {
	case domainGame.StatusRoundsOverview:
		m.startRound(1)

	case domainGame.StatusRoundStart:
		m.transitionToQuestionSelect()

	case domainGame.StatusQuestionSelect:
		m.autoSelectQuestion()

	case domainGame.StatusQuestionShow:
		m.transitionFromQuestionShow()

	case domainGame.StatusButtonPress:
		m.skipQuestion()

	case domainGame.StatusAnswering:
		m.handleAnswerTimeout()

	case domainGame.StatusAnswerJudging:
		m.handleAnswerTimeout()

	case domainGame.StatusSecretTransfer:
		m.handleSecretTransferTimeout()

	case domainGame.StatusStakeBetting:
		m.handleStakeBettingTimeout()

	case domainGame.StatusForAllAnswering:
		m.finishForAllQuestion()

	case domainGame.StatusForAllResults:
		m.continueGame()
	}
}

func (m *Manager) handleTimerTick() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch m.game.Status {
	case domainGame.StatusQuestionSelect,
		domainGame.StatusButtonPress,
		domainGame.StatusAnswering,
		domainGame.StatusAnswerJudging,
		domainGame.StatusSecretTransfer,
		domainGame.StatusStakeBetting,
		domainGame.StatusForAllAnswering:
		m.BroadcastStateUnlocked()
	}
}

func (m *Manager) autoSelectQuestion() {
	round := m.pack.GetRound(m.game.CurrentRound)
	if round == nil {
		m.endRound()
		return
	}

	for _, theme := range round.Themes {
		for _, question := range theme.Questions {
			if question.IsAvailable() {
				m.selectQuestion(theme, question)
				return
			}
		}
	}

	m.endRound()
}

func (m *Manager) skipQuestion() {
	m.buttonPress.Close()

	if m.buttonPress.HasPresses() {
		winner := m.buttonPress.GetWinner()
		if winner != nil {
			m.game.SetActivePlayer(winner.UserID)
			m.game.UpdateStatus(domainGame.StatusAnswerJudging)
			m.BroadcastState()
			m.timer.Start(30 * time.Second)
			return
		}
	}

	m.continueGame()
}

func (m *Manager) handleAnswerTimeout() {
	if m.game.ActivePlayer == nil {
		return
	}

	p := m.game.Players[*m.game.ActivePlayer]
	questionPrice := m.game.CurrentQuestion.Price
	p.SubtractScore(questionPrice)

	m.continueGame()
}

func (m *Manager) handleSecretTransferTimeout() {
	hostID := m.findHost()

	for userID, p := range m.game.Players {
		if p.Role != 0 {
			m.transferSecretToPlayer(hostID, userID)
			return
		}
	}

	m.continueGame()
}

func (m *Manager) handleStakeBettingTimeout() {
	if m.game.ActivePlayer == nil || m.stakeInfo == nil {
		m.continueGame()
		return
	}

	m.placeStake(*m.game.ActivePlayer, m.stakeInfo.MinBet, false)
}

func (m *Manager) finishForAllQuestion() {
	m.forAllCollector.Close()

	results := m.forAllCollector.GetResults(func(userAnswer, correctAnswer string) bool {
		return m.game.CurrentQuestion.ValidateAnswer(userAnswer)
	})

	for userID, result := range results {
		p := m.game.Players[userID]
		if result.IsCorrect {
			p.AddScore(result.ScoreDelta)
		} else {
			p.SubtractScore(-result.ScoreDelta)
		}
	}

	m.game.UpdateStatus(domainGame.StatusForAllResults)
	m.BroadcastState()
	m.timer.Start(5 * time.Second)
}

