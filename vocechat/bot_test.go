package vocechat

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBot(t *testing.T) {
	var (
		fp  string
		err error
		bot *Bot
	)

	fp = "configs/vocechat.yaml"
	bot, err = NewBot(fp)
	require.Nil(t, err)

	_, err = bot.SendMsg("user", "Hello! My name is d2jvkpn")
	require.Nil(t, err)
}
