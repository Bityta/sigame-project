package game

import (
	domainGame "github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
)

func (m *Manager) BroadcastStateUnlocked() {
	state := m.buildGameState()
	m.broadcastState(state)
}

func (m *Manager) BroadcastState() {
	state := m.buildGameState()
	m.broadcastState(state)
}

func (m *Manager) buildGameState() *domainGame.State {
	state := &domainGame.State{
		GameID:        m.game.ID,
		Status:        m.game.Status,
		CurrentRound:  m.game.CurrentRound,
		Players:       make([]player.State, 0, len(m.game.Players)),
		ActivePlayer:  m.game.ActivePlayer,
		TimeRemaining: m.timer.Remaining(),
	}

	for _, p := range m.game.Players {
		state.Players = append(state.Players, p.ToState())
	}

	if m.game.Status == domainGame.StatusRoundsOverview {
		state.AllRounds = make([]domainGame.RoundOverview, 0, len(m.pack.Rounds))
		for i, round := range m.pack.Rounds {
			themeNames := make([]string, 0, len(round.Themes))
			for _, theme := range round.Themes {
				themeNames = append(themeNames, theme.Name)
			}
			state.AllRounds = append(state.AllRounds, domainGame.RoundOverview{
				RoundNumber: i + RankStartIndex,
				Name:        round.Name,
				ThemeNames:  themeNames,
			})
		}
	}

	if m.game.CurrentRound > 0 && m.game.CurrentRound <= len(m.pack.Rounds) {
		round := m.pack.Rounds[m.game.CurrentRound-1]
		state.RoundName = round.Name

		state.Themes = make([]pack.ThemeState, 0, len(round.Themes))
		for _, theme := range round.Themes {
			includeText := m.game.Status == domainGame.StatusQuestionShow
			state.Themes = append(state.Themes, theme.ToState(includeText))
		}
	}

	if m.game.CurrentQuestion != nil {
		questionState := m.game.CurrentQuestion.ToStateWithAnswer(true)
		state.CurrentQuestion = &questionState
	}

	if m.game.Status == domainGame.StatusGameEnd {
		state.Winners = m.game.Winners
		state.FinalScores = m.game.FinalScores
	}

	if m.stakeInfo != nil && m.game.Status == domainGame.StatusStakeBetting {
		state.StakeInfo = m.stakeInfo
	}

	if m.secretTarget != nil {
		state.SecretTarget = m.secretTarget
	}

	return state
}

func (m *Manager) broadcastState(state *domainGame.State) {
	data := m.serializeState(state)
	if data != nil {
		m.hub.Broadcast(m.game.ID, data)
	}
}

func (m *Manager) serializeState(state *domainGame.State) []byte {
	return nil
}

func (m *Manager) sendStateToClient(client interface{}, state *domainGame.State) {
}

func (m *Manager) sendRoundMediaManifest(roundNumber int, manifest interface{}, totalSize int64) {
}

