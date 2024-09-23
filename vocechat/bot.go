package vocechat

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Bot struct {
	Server        string        `mapstructure:"server"`
	Timeout       time.Duration `mapstructure:"timeout"`
	TlsSkipVerify bool          `mapstructure:"tls_skip_verify"`
	UserAgent     string        `mapstructure:"user_agent"`

	Bot struct {
		ApiKey string `mapstructure:"api_key"`

		SendToUser  uint32 `mapstructure:"send_to_user"`
		SendToGroup uint32 `mapstructure:"send_to_group"`
	} `mapstructure:"bot"`

	*http.Client
}

func NewBot(fp string, keys ...string) (bot *Bot, err error) {
	var vp *viper.Viper

	vp = viper.New()
	// vp.SetConfigName(name)
	vp.SetConfigType("yaml")

	vp.SetConfigFile(fp)
	if err = vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config %q: %w", fp, err)
	}
	vp.SetDefault("timeout", "5s")
	vp.SetDefault("user_agent", "Mozilla/5.0 Gecko/20100101 Firefox/130.0")

	bot = new(Bot)

	if len(keys) > 0 {
		err = vp.UnmarshalKey(keys[0], bot)
	} else {
		err = vp.Unmarshal(bot)
	}

	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if bot.Server = strings.TrimRight(bot.Server, "/"); bot.Server == "" {
		return nil, errors.New("empty server")
	}
	if bot.Bot.ApiKey == "" {
		return nil, errors.New("api_key is empty")
	}
	if bot.Timeout < 0 {
		return nil, fmt.Errorf("invalid timeout: %v", bot.Timeout)
	}

	bot.Client = new(http.Client)
	bot.Client.Timeout = bot.Timeout

	if bot.TlsSkipVerify {
		bot.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true},
		}
	}

	return bot, nil
}

func (self *Bot) SendMsg(typ, msg string) (kind string, err error) {
	var (
		p        string
		reader   *strings.Reader
		request  *http.Request
		response *http.Response
	)

	if msg = strings.TrimSpace(msg); msg == "" {
		return "invalid_parameter", errors.New("msg is empty")
	}

	switch typ {
	case "user":
		if self.Bot.SendToUser <= 0 {
			return "invalid_parameter", errors.New("bot.send_to_user is empty")
		}
	case "group":
		if self.Bot.SendToGroup <= 0 {
			return "invalid_parameter", errors.New("bot.send_to_group is empty")
		}
	default:
		return "invalid_parameter", errors.New("invalid type")
	}

	reader = strings.NewReader(msg)

	if typ == "user" {
		p = fmt.Sprintf("%s/api/bot/send_to_user/%d", self.Server, self.Bot.SendToUser)
	} else {
		p = fmt.Sprintf("%s/api/bot/send_to_group/%d", self.Server, self.Bot.SendToGroup)
	}

	if request, err = http.NewRequest("POST", p, reader); err != nil {
		return "create_request", err
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Origin", self.Server)
	request.Header.Set("User-Agent", self.UserAgent)

	request.Header.Set("X-API-Key", self.Bot.ApiKey)

	if response, err = self.Do(request); err != nil {
		return "request_failed", err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "respone_failed", fmt.Errorf("send response: %d", response.StatusCode)
	}

	return "ok", nil
}
