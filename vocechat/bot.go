package vocechat

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (client *Client) BotSend(msg string) (err error) {
	var (
		p        string
		reader   *strings.Reader
		request  *http.Request
		response *http.Response
	)

	if client.ApiKey == "" {
		return errors.New("api_key is empty")
	}
	if msg = strings.TrimSpace(msg); msg == "" {
		return errors.New("msg is empty")
	}

	reader = strings.NewReader(msg)

	if client.SendToKind == "group" {
		p = fmt.Sprintf("%s/api/bot/send_to_group/%d", client.Server, client.SendToId)
	} else { // user
		p = fmt.Sprintf("%s/api/bot/send_to_user/%d", client.Server, client.SendToId)
	}

	if request, err = http.NewRequest("POST", p, reader); err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("X-API-Key", client.ApiKey)

	if response, err = client.cli.Do(request); err != nil {
		return err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("send response: %d, %s", response.StatusCode, response.Status)
	}

	return nil
}
