package line

import (
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Line struct {
	ChannelSecret string
	ChannelToken  string
	Client        *linebot.Client
}

func LineConnection() (*Line, error) {
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	lineInfo := Line {
		ChannelSecret: os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		ChannelToken: os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
		Client: bot,
	}

	if err != nil {
		return nil, err
	}

	return &lineInfo, nil
}