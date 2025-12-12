package game

import (
	"context"
	"time"

	"sigame/game/internal/domain/event"
	domainGame "sigame/game/internal/domain/game"
)

func (m *Manager) startGame() {
	now := time.Now()
	m.game.StartedAt = &now
	m.game.CurrentRound = InitialRoundNumber

	m.logEvent(event.TypeGameStarted)

	m.showRoundsOverview()
}

func (m *Manager) showRoundsOverview() {
	m.game.UpdateStatus(domainGame.StatusRoundsOverview)
	m.BroadcastState()

	m.timer.Start(RoundsOverviewDuration)
}

func (m *Manager) startRound(roundNumber int) {
	if roundNumber > m.pack.TotalRounds() {
		m.endGame()
		return
	}

	m.game.CurrentRound = roundNumber
	m.game.UpdateStatus(domainGame.StatusRoundStart)

	evt := event.New(m.game.ID, event.TypeRoundStarted).WithRound(roundNumber)
	m.eventLogger.LogEvent(context.Background(), evt)

	round := m.pack.GetRound(roundNumber)
	m.mediaTracker.Reset(roundNumber)
	m.mediaTracker.BuildManifest(round)

	for userID := range m.game.Players {
		m.mediaTracker.RegisterClient(userID)
	}

	if m.mediaTracker.HasMedia() {
		manifest, totalSize := m.mediaTracker.GetManifest()
		m.sendRoundMediaManifest(roundNumber, manifest, totalSize)
	}

	m.BroadcastState()
	m.timer.Start(RoundIntroDuration)
}

func (m *Manager) endRound() {
	m.game.UpdateStatus(domainGame.StatusRoundEnd)

	evt := event.New(m.game.ID, event.TypeRoundFinished).WithRound(m.game.CurrentRound)
	m.eventLogger.LogEvent(context.Background(), evt)

	m.BroadcastState()

	currentRound := m.game.CurrentRound
	totalRounds := m.pack.TotalRounds()

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

func (m *Manager) endGame() {
	m.game.UpdateStatus(domainGame.StatusGameEnd)
	now := time.Now()
	m.game.FinishedAt = &now

	m.game.Winners = m.calculateWinners()
	m.game.FinalScores = m.calculateFinalScores()

	m.logEvent(event.TypeGameFinished)

	m.BroadcastState()

	if m.timerTicker != nil {
		m.timerTicker.Stop()
	}
}

func (m *Manager) transitionToQuestionSelect() {
	host, err := m.game.GetHost()
	if err != nil {
		return
	}
	m.game.SetActivePlayer(host.UserID)

	m.game.UpdateStatus(domainGame.StatusQuestionSelect)

	m.timer.Start(time.Duration(m.game.Settings.TimeForChoice) * time.Second)
	m.BroadcastState()
}

func (m *Manager) transitionToButtonPress() {
	m.game.UpdateStatus(domainGame.StatusButtonPress)
	m.buttonPress.Reset()

	m.timer.Start(time.Duration(m.game.Settings.TimeForAnswer) * time.Second)
	m.BroadcastState()
}

func (m *Manager) transitionToAnswerJudging() {
	m.game.UpdateStatus(domainGame.StatusAnswerJudging)
	m.BroadcastState()
	m.timer.Start(30 * time.Second)
}

func (m *Manager) transitionFromQuestionShow() {
	if m.game.CurrentQuestion == nil {
		m.transitionToButtonPress()
		return
	}

	questionType := m.game.CurrentQuestion.GetType()
	if questionType == "secret" || questionType == "stake" {
		m.transitionToAnswerJudging()
	} else {
		m.transitionToButtonPress()
	}
}

func (m *Manager) continueGame() {
	m.game.ClearCurrentQuestion()
	m.stakeInfo = nil
	m.secretTarget = nil
	m.forAllCollector.Reset()

	round := m.pack.GetRound(m.game.CurrentRound)
	if round.IsComplete() {
		m.endRound()
		return
	}

	host, _ := m.game.GetHost()
	m.game.SetActivePlayer(host.UserID)

	m.game.UpdateStatus(domainGame.StatusQuestionSelect)
	m.BroadcastState()

	m.timer.Start(time.Duration(m.game.Settings.TimeForChoice) * time.Second)
}

