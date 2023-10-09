package vocechat

import (
	"fmt"
	"testing"

	"github.com/d2jvkpn/gotk"
)

func TestUserSend(t *testing.T) {
	var (
		fp     string
		err    error
		client *Client
	)

	if fp, err = gotk.RootFile("configs/local.yaml"); err != nil {
		t.Fatal(err)
	}

	if client, err = NewClient(fp, "vocechat_user"); err != nil {
		t.Fatal(err)
	}

	msg := fmt.Sprintf("Hello, %s! My name is %s", client.SendToName, client.Username)
	if err = client.UserSend(msg); err != nil {
		t.Fatal(err)
	}
}

func TestBotSend(t *testing.T) {
	var (
		fp     string
		err    error
		client *Client
	)

	if fp, err = gotk.RootFile("configs/local.yaml"); err != nil {
		t.Fatal(err)
	}

	if client, err = NewClient(fp, "vocechat_bot"); err != nil {
		t.Fatal(err)
	}

	msg := fmt.Sprintf("Hello, %s! My name is %s", client.SendToName, client.Username)
	if err = client.BotSend(msg); err != nil {
		t.Fatal(err)
	}
}
