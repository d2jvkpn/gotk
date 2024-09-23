package vocechat

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// kind=["invalid_parameter", "request", "response", "ok"]
func (self *Client) BotSendMsg(typ, msg string, targetId uint) (kind string, err error) {
	var (
		p        string
		reader   *strings.Reader
		request  *http.Request
		response *http.Response
	)

	if msg = strings.TrimSpace(msg); msg == "" {
		return "invalid_parameter", errors.New("msg is empty")
	}

	if targetId <= 0 {
		return "invalid_parameter", errors.New("invalid targetId")
	}

	if self.ApiKey == "" {
		return "invalid_parameter", errors.New("api_key is empty")
	}

	switch typ {
	case "user", "group":
	default:
		return "invalid_parameter", errors.New("invalid type")
	}

	reader = strings.NewReader(msg)

	if typ == "user" {
		p = fmt.Sprintf("%s/api/bot/send_to_user/%d", self.Server, targetId)
	} else {
		p = fmt.Sprintf("%s/api/bot/send_to_group/%d", self.Server, targetId)
	}

	if request, err = http.NewRequest("POST", p, reader); err != nil {
		return "invalid_parameter", err
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Origin", self.Server)
	request.Header.Set("User-Agent", self.UserAgent)

	request.Header.Set("X-API-Key", self.ApiKey)

	if response, err = self.Do(request); err != nil {
		return "request", err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "respone", fmt.Errorf("response: %d", response.StatusCode)
	}

	return "ok", nil
}
