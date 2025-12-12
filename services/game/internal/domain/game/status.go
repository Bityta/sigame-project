package game

type Status string

const (
	StatusWaiting        Status = "waiting"
	StatusRoundsOverview Status = "rounds_overview"
	StatusRoundStart     Status = "round_start"
	StatusQuestionSelect Status = "question_select"
	StatusQuestionShow   Status = "question_show"
	StatusButtonPress    Status = "button_press"
	StatusAnswering      Status = "answering"
	StatusAnswerJudging  Status = "answer_judging"
	StatusRoundEnd       Status = "round_end"
	StatusGameEnd        Status = "game_end"
	StatusFinished       Status = "finished"
	StatusCancelled      Status = "cancelled"
	StatusSecretTransfer Status = "secret_transfer"
	StatusStakeBetting   Status = "stake_betting"
	StatusForAllAnswering Status = "for_all_answering"
	StatusForAllResults   Status = "for_all_results"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsActive() bool {
	return s != StatusFinished && s != StatusCancelled
}

func (s Status) IsPlaying() bool {
	return s != StatusWaiting && s != StatusFinished && s != StatusCancelled
}

