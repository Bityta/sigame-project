package game

import "time"

const (
	ManagerActionChannelBuffer = 100
	RoundsOverviewDuration     = 5 * time.Second
	RoundIntroDuration         = 3 * time.Second
	AnswerJudgingDuration      = 30 * time.Second
	TimerUpdateInterval        = 1 * time.Second
	InitialRoundNumber         = 0
	MaxIntValue                = int(^uint(0) >> 1)

	QuestionReadDuration         = 3 * time.Second
	SecretTransferDuration       = 30 * time.Second
	StakeBettingDuration         = 20 * time.Second
	ButtonPressCollectionWindow  = 150 * time.Millisecond
	MediaStartDelay              = 300 * time.Millisecond
	ForAllResultsDisplayDuration = 5 * time.Second
	RoundEndDelay                = 5 * time.Second
	TopWinnersCount              = 3
	MediaIDSuffix                = "_media"
	DefaultMediaDurationMs       = 5000
)

