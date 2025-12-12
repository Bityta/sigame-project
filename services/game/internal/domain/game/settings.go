package game

type Settings struct {
	TimeForAnswer int `json:"time_for_answer" binding:"required"`
	TimeForChoice int `json:"time_for_choice" binding:"required"`
}

func DefaultSettings() Settings {
	return Settings{
		TimeForAnswer: 30,
		TimeForChoice: 20,
	}
}

func (s Settings) Validate() error {
	if s.TimeForAnswer <= 0 || s.TimeForAnswer > 300 {
		return ErrInvalidSettings
	}
	if s.TimeForChoice <= 0 || s.TimeForChoice > 300 {
		return ErrInvalidSettings
	}
	return nil
}

