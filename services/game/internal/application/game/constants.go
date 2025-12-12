package game

import "time"

const (
	RTTCompensationFactor = 2

	MediaTypeText  = "text"
	MediaTypeImage = "image"
	MediaTypeAudio = "audio"
	MediaTypeVideo = "video"

	MediaSizeImage   = 500_000
	MediaSizeAudio   = 3_000_000
	MediaSizeVideo   = 10_000_000
	MediaSizeDefault = 100_000

	PercentComplete = 100

	TimerChannelBufferSize = 1
	TimerInactiveRemaining = 0
	RankStartIndex = 1

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
)

