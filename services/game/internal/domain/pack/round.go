package pack

type Round struct {
	ID          string
	RoundNumber int
	Name        string
	Themes      []*Theme
}

func (r *Round) GetAvailableThemes() []*Theme {
	available := make([]*Theme, 0)
	for _, theme := range r.Themes {
		if theme.HasAvailableQuestions() {
			available = append(available, theme)
		}
	}
	return available
}

func (r *Round) IsComplete() bool {
	for _, theme := range r.Themes {
		for _, question := range theme.Questions {
			if !question.Used {
				return false
			}
		}
	}
	return true
}

func (r *Round) FindQuestion(themeID, questionID string) *Question {
	for _, theme := range r.Themes {
		if theme.ID == themeID || theme.Name == themeID {
			question := theme.FindQuestion(questionID)
			if question != nil {
				return question
			}
		}
	}
	return nil
}

func (r *Round) FindTheme(themeIDOrName string) *Theme {
	for _, theme := range r.Themes {
		if theme.ID == themeIDOrName || theme.Name == themeIDOrName {
			return theme
		}
	}
	return nil
}

func (r *Round) TotalQuestions() int {
	total := 0
	for _, theme := range r.Themes {
		total += len(theme.Questions)
	}
	return total
}

