package event

type Type string

const (
	TypeGameCreated       Type = "GAME_CREATED"
	TypeGameStarted       Type = "GAME_STARTED"
	TypeGameFinished      Type = "GAME_FINISHED"
	TypeGameCancelled     Type = "GAME_CANCELLED"
	TypePlayerJoined      Type = "PLAYER_JOINED"
	TypePlayerLeft        Type = "PLAYER_LEFT"
	TypePlayerReady       Type = "PLAYER_READY"
	TypeRoundStarted      Type = "ROUND_STARTED"
	TypeRoundFinished     Type = "ROUND_FINISHED"
	TypeQuestionSelected  Type = "QUESTION_SELECTED"
	TypeQuestionShown     Type = "QUESTION_SHOWN"
	TypeButtonPressed     Type = "BUTTON_PRESSED"
	TypeAnswerSubmitted   Type = "ANSWER_SUBMITTED"
	TypeAnswerCorrect     Type = "ANSWER_CORRECT"
	TypeAnswerIncorrect   Type = "ANSWER_INCORRECT"
	TypeScoreChanged      Type = "SCORE_CHANGED"
	TypeTimerStarted      Type = "TIMER_STARTED"
	TypeTimerExpired      Type = "TIMER_EXPIRED"
)

func (t Type) String() string {
	return string(t)
}

