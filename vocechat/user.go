package vocechat

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func (client *Client) UserSend(msg string) (err error) {
	var (
		bts      []byte
		p, token string
		reader   *strings.Reader
		request  *http.Request
		response *http.Response
	)

	if client.Email == "" || client.Password == "" {
		return errors.New("email or password is empty")
	}
	if msg = strings.TrimSpace(msg); msg == "" {
		return errors.New("msg is empty")
	}

	//
	reader = strings.NewReader(fmt.Sprintf(
		`{"credential":{"email":%q,"password":%q,"type":"password"}}`,
		client.Email, client.Password,
	))

	p = fmt.Sprintf("%s/api/token/login", client.Server)
	if request, err = http.NewRequest("POST", p, reader); err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", client.Server)

	if response, err = client.cli.Do(request); err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("login response: %d, %s", response.StatusCode, response.Status)
	}

	if bts, err = io.ReadAll(response.Body); err != nil {
		response.Body.Close()
		return fmt.Errorf("read body: %w", err)
	}
	response.Body.Close()

	if token = gjson.GetBytes(bts, "token").String(); token == "" {
		return fmt.Errorf("empty token")
	}

	//
	reader = strings.NewReader(msg)

	if client.SendToKind == "group" {
		p = fmt.Sprintf("%s/api/group/%d/send", client.Server, client.SendToId)
	} else { // user
		p = fmt.Sprintf("%s/api/user/%d/send", client.Server, client.SendToId)
	}

	if request, err = http.NewRequest("POST", p, reader); err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("X-API-Key", token)

	if response, err = client.cli.Do(request); err != nil {
		return err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("send response: %d, %s", response.StatusCode, response.Status)
	}

	return nil
}
