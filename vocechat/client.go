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

type Client struct {
	Server        string        `mapstructure:"server"`
	Timeout       time.Duration `mapstructure:"timeout"`
	TlsSkipVerify bool          `mapstructure:"tls_skip_verify"`
	UserAgent     string        `mapstructure:"user_agent"`

	// bot
	ApiKey string `mapstructure:"api_key"`

	// user
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`

	*http.Client
}

func NewClient(fp string, keys ...string) (client *Client, err error) {
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

	return ClientFromViper(vp, keys...)
}

func ClientFromViper(vp *viper.Viper, keys ...string) (client *Client, err error) {
	client = new(Client)

	if len(keys) > 0 {
		err = vp.UnmarshalKey(keys[0], client)
	} else {
		err = vp.Unmarshal(client)
	}

	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if client.Server = strings.TrimRight(client.Server, "/"); client.Server == "" {
		return nil, errors.New("empty server")
	}
	if client.Timeout < 0 {
		return nil, fmt.Errorf("invalid timeout: %v", client.Timeout)
	}

	client.Client = new(http.Client)
	client.Client.Timeout = client.Timeout

	if client.TlsSkipVerify {
		client.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true},
		}
	}

	return client, nil
}

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
