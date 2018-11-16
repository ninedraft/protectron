package bot

type Bot struct {
	token string
}

func New(token string) Bot {
	return Bot{
		token: token,
	}
}

func (bot Bot) Run() error {
	return nil
}
