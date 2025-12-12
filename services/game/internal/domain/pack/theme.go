package pack

type Theme struct {
	ID        string
	Name      string
	Questions []*Question
}

type ThemeState struct {
	Name      string          `json:"name" binding:"required"`
	Questions []QuestionState `json:"questions" binding:"required"`
}

func (t *Theme) GetAvailableQuestions() []*Question {
	available := make([]*Question, 0)
	for _, q := range t.Questions {
		if q.IsAvailable() {
			available = append(available, q)
		}
	}
	return available
}

func (t *Theme) HasAvailableQuestions() bool {
	return len(t.GetAvailableQuestions()) > 0
}

func (t *Theme) ToState(includeQuestionText bool) ThemeState {
	questions := make([]QuestionState, len(t.Questions))
	for i, q := range t.Questions {
		questions[i] = q.ToState(includeQuestionText)
	}

	return ThemeState{
		Name:      t.Name,
		Questions: questions,
	}
}

func (t *Theme) FindQuestion(questionID string) *Question {
	for _, q := range t.Questions {
		if q.ID == questionID {
			return q
		}
	}
	return nil
}

