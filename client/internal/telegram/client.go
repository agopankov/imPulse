package telegram

import (
	"time"

	tele "gopkg.in/telebot.v3"
)

type Client struct {
	botToken string
	bot      *tele.Bot
}

func NewClient(botToken string) (*Client, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		botToken: botToken,
		bot:      bot,
	}, nil
}

func (c *Client) Start() {
	c.bot.Start()
}

func (c *Client) Bot() *tele.Bot {
	return c.bot
}

func (c *Client) HandleText(handler func(m *tele.Message)) {
	c.bot.Handle(tele.OnText, func(c tele.Context) error {
		handler(c.Message())
		return nil
	})
}

func (c *Client) SendMessage(recipient *tele.User, text string) (*tele.Message, error) {
	return c.bot.Send(recipient, text)
}

func (c *Client) HandleCommand(command string, handler func(m *tele.Message)) {
	c.bot.Handle(command, func(c tele.Context) error {
		handler(c.Message())
		return nil
	})
}

func (c *Client) HandleOnMessage(fn func(m *tele.Message)) {
	c.bot.Handle(tele.OnText, func(c tele.Context) error {
		fn(c.Message())
		return nil
	})
}
