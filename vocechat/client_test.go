package vocechat

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	var (
		// fp  string
		err    error
		client *Client
	)

	// fp = "configs/local.yaml"
	// bot, err = NewBot(fp)
	client, err = ClientFromViper(_TestConfig)
	require.Nil(t, err)

	_, err = client.BotSendMsg("user", "Hello! My name is d2jvkpn", 1)
	require.Nil(t, err)
}
